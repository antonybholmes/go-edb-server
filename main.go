package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"runtime"

	"github.com/gorilla/sessions"
	"github.com/rs/zerolog"

	"github.com/antonybholmes/go-auth"
	"github.com/antonybholmes/go-auth/tokengen"
	"github.com/antonybholmes/go-auth/userdbcache"
	"github.com/antonybholmes/go-dna/dnadbcache"
	"github.com/antonybholmes/go-edb-server/consts"
	adminroutes "github.com/antonybholmes/go-edb-server/routes/admin"
	"github.com/antonybholmes/go-edb-server/routes/authentication"
	"github.com/antonybholmes/go-edb-server/routes/authorization"
	dnaroutes "github.com/antonybholmes/go-edb-server/routes/modules/dna"
	generoutes "github.com/antonybholmes/go-edb-server/routes/modules/gene"
	geneconvroutes "github.com/antonybholmes/go-edb-server/routes/modules/geneconv"
	gexroutes "github.com/antonybholmes/go-edb-server/routes/modules/gex"
	motiftogeneroutes "github.com/antonybholmes/go-edb-server/routes/modules/motiftogene"
	mutationroutes "github.com/antonybholmes/go-edb-server/routes/modules/mutation"
	pathwayroutes "github.com/antonybholmes/go-edb-server/routes/modules/pathway"
	"github.com/antonybholmes/go-geneconv/geneconvdbcache"
	"github.com/antonybholmes/go-genes/genedbcache"
	"github.com/antonybholmes/go-gex/gexdbcache"
	"github.com/antonybholmes/go-motiftogene/motiftogenedb"
	"github.com/antonybholmes/go-mutations/mutationdbcache"
	"github.com/antonybholmes/go-pathway/pathwaydbcache"
	"github.com/antonybholmes/go-sys/env"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo-contrib/session"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	_ "github.com/mattn/go-sqlite3"
	"gopkg.in/natefinch/lumberjack.v2"
)

type AboutResp struct {
	Name      string `json:"name"`
	Copyright string `json:"copyright"`
	Version   string `json:"version"`
	Updated   string `json:"updated"`
}

type InfoResp struct {
	IpAddr string `json:"ipAddr"`
	Arch   string `json:"arch"`
}

// var store *sqlitestore.SqliteStore
var store *sessions.CookieStore

func init() {

	env.Ls()
	// store = sys.Must(sqlitestore.NewSqliteStore("data/users.db",
	// 	"sessions",
	// 	"/",
	// 	auth.MAX_AGE_7_DAYS_SECS,
	// 	[]byte(consts.SESSION_SECRET)))

	store = sessions.NewCookieStore([]byte(consts.SESSION_SECRET))
	// store.Options = &sessions.Options{
	// 	Path:     "/",
	// 	MaxAge:   auth.MAX_AGE_7_DAYS_SECS,
	// 	HttpOnly: false,
	// 	Secure:   true,
	// 	SameSite: http.SameSiteNoneMode}

	userdbcache.InitCache("data/users.db")

	//mailserver.Init()

	dnadbcache.InitCache("data/modules/dna")
	genedbcache.InitCache("data/modules/genes")
	//microarraydb.InitDB("data/microarray")

	gexdbcache.InitCache("data/modules/gex")

	mutationdbcache.InitCache("data/modules/mutations")

	geneconvdbcache.InitCache("data/modules/geneconv/geneconv.db")

	motiftogenedb.InitCache("data/modules/motiftogene/motiftogene.db")

	pathwaydbcache.InitCache("data/modules/pathway/pathway.db")
}

func main() {
	//env.Reload()
	//env.Load("consts.env")
	//env.Load("version.env")

	//consts.Init()

	tokengen.Init(consts.JWT_RSA_PRIVATE_KEY)

	//env.Load()

	// list env to see what is loaded
	//env.Ls()

	//initCache()

	// test redis
	//email := gomailer.RedisQueueEmail{To: "antony@antonybholmes.dev"}
	//rdb.PublishEmail(&email)

	//
	// Set logging to file
	//

	//log.SetOutput(logFile)

	//
	// end logging setup
	//

	e := echo.New()

	e.Use(middleware.BodyLimit("2M"))

	//e.Use(middleware.Logger())

	e.Use(session.Middleware(store))

	// write to both stdout and log file
	// f := env.GetStr("LOG_FILE", fmt.Sprintf("logs/%s.log", consts.APP_NAME))

	// logFile, err := os.OpenFile(f, os.O_APPEND|os.O_RDWR|os.O_CREATE, 0644)

	// if err != nil {
	// 	//return nil, err
	// }

	// // to prevent file closing before program exits
	// defer logFile.Close()

	fileLogger := &lumberjack.Logger{
		Filename:   fmt.Sprintf("logs/%s.log", consts.APP_NAME),
		MaxSize:    5, //
		MaxBackups: 10,
		MaxAge:     14,
		Compress:   true,
	}

	logger := zerolog.New(io.MultiWriter(zerolog.ConsoleWriter{Out: os.Stderr}, fileLogger)).With().Timestamp().Logger()

	// we use != development because it means we need to set the env variable in order
	// to see debugging work. The default is to assume production, in which case we use
	// lumberjack
	if os.Getenv("APP_ENV") != "development" {
		logger = zerolog.New(io.MultiWriter(zerolog.ConsoleWriter{Out: os.Stderr}, fileLogger)).With().Timestamp().Logger()
	}

	// We cache options regarding ttl so some session routes need to be in an object
	sr := authentication.NewSessionRoutes()

	e.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogURI:    true,
		LogStatus: true,
		LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
			logger.Info().
				Str("URI", v.URI).
				Int("status", v.Status).
				Msg("request")

			return nil
		},
	}))

	//e.Use(loggerMiddleware)
	e.Use(middleware.Recover())
	//e.Use(middleware.CORS())

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"https://edb.rdf-lab.org", "http://localhost:8000"},
		AllowMethods: []string{http.MethodGet, http.MethodDelete, http.MethodPost},
		// for sharing session cookie for validating logins etc
		AllowCredentials: true,
	}))

	// Configure middleware with the custom claims type which
	// will parse our jwt with scope etc
	config := echojwt.Config{
		NewClaimsFunc: func(c echo.Context) jwt.Claims {
			return &auth.TokenClaims{}
		},
		SigningKey: consts.JWT_RSA_PRIVATE_KEY,
		// Have to tell it to use the public key for verification
		KeyFunc: func(token *jwt.Token) (interface{}, error) {
			return consts.JWT_RSA_PUBLIC_KEY, nil
		},
	}
	jwtMiddleWare := echojwt.WithConfig(config)

	//
	// Routes
	//

	e.GET("/about", func(c echo.Context) error {
		return c.JSON(http.StatusOK,
			AboutResp{Name: consts.NAME,
				Version:   consts.VERSION,
				Updated:   consts.UPDATED,
				Copyright: consts.COPYRIGHT})
	})

	e.GET("/info", func(c echo.Context) error {
		return c.JSON(http.StatusOK, InfoResp{Arch: runtime.GOARCH, IpAddr: c.RealIP()})
	})

	adminGroup := e.Group("/admin")
	adminGroup.Use(jwtMiddleWare,
		JwtIsAccessTokenMiddleware,
		JwtHasAdminPermissionMiddleware)

	adminGroup.GET("/roles", adminroutes.RolesRoute)

	adminUsersGroup := adminGroup.Group("/users")

	adminUsersGroup.POST("", adminroutes.UsersRoute)
	adminUsersGroup.GET("/stats", adminroutes.UserStatsRoute)
	adminUsersGroup.POST("/update", adminroutes.UpdateUserRoute)
	adminUsersGroup.POST("/add", adminroutes.AddUserRoute)
	adminUsersGroup.DELETE("/delete/:publicId", adminroutes.DeleteUserRoute)

	// Allow users to sign up for an account
	e.POST("/signup", authentication.SignupRoute)

	//
	// user groups: start
	//

	authGroup := e.Group("/auth")

	authGroup.POST("/signin", authentication.UsernamePasswordSignInRoute)

	emailGroup := authGroup.Group("/email")

	emailGroup.POST("/verify",
		authentication.EmailAddressVerifiedRoute,
		jwtMiddleWare)

	// with the correct token, performs the update
	emailGroup.POST("/reset", authentication.SendResetEmailEmailRoute, jwtMiddleWare)
	// with the correct token, performs the update
	emailGroup.POST("/update", authentication.UpdateEmailRoute, jwtMiddleWare)

	passwordGroup := authGroup.Group("/passwords")

	// sends a reset link
	passwordGroup.POST("/reset", authentication.SendResetPasswordFromUsernameEmailRoute)

	// with the correct token, updates a password
	passwordGroup.POST("/update", authentication.UpdatePasswordRoute, jwtMiddleWare)

	passwordlessGroup := authGroup.Group("/passwordless")

	passwordlessGroup.POST("/email", func(c echo.Context) error {
		return authentication.PasswordlessSigninEmailRoute(c, nil)
	})

	passwordlessGroup.POST("/signin",
		authentication.PasswordlessSignInRoute,
		jwtMiddleWare)

	tokenGroup := authGroup.Group("/tokens")
	tokenGroup.Use(jwtMiddleWare)
	tokenGroup.POST("/info", authorization.TokenInfoRoute)
	tokenGroup.POST("/access", authorization.NewAccessTokenRoute)

	usersGroup := authGroup.Group("/users")
	usersGroup.Use(jwtMiddleWare,
		JwtIsAccessTokenMiddleware)

	usersGroup.POST("", authorization.UserRoute)

	usersGroup.POST("/update", authorization.UpdateUserRoute)

	//usersGroup.POST("/passwords/update", authentication.UpdatePasswordRoute)

	//
	// Deal with logins where we want a session
	//

	sessionGroup := e.Group("/sessions")

	//sessionAuthGroup := sessionGroup.Group("/auth")

	sessionGroup.POST("/signin", sr.SessionUsernamePasswordSignInRoute)

	sessionGroup.POST("/passwordless/signin",
		sr.SessionPasswordlessValidateSignInRoute,
		jwtMiddleWare)

	sessionGroup.POST("/signout", authentication.SessionSignOutRoute)

	//sessionGroup.POST("/email/reset", authentication.SessionSendResetEmailEmailRoute)

	//sessionGroup.POST("/password/reset", authentication.SessionSendResetPasswordEmailRoute)

	sessionGroup.POST("/tokens/access", authentication.SessionNewAccessTokenRoute, SessionIsValidMiddleware)

	sessionUsersGroup := sessionGroup.Group("/users")
	sessionUsersGroup.Use(SessionIsValidMiddleware)

	sessionUsersGroup.GET("", authentication.SessionUserRoute)

	sessionUsersGroup.POST("/update", authorization.SessionUpdateUserRoute)

	// sessionPasswordGroup := sessionAuthGroup.Group("/passwords")
	// sessionPasswordGroup.Use(SessionIsValidMiddleware)

	// sessionPasswordGroup.POST("/update", func(c echo.Context) error {
	// 	return authentication.SessionUpdatePasswordRoute(c)
	// })

	//
	// sessions: end
	//

	//
	// passwordless groups: end
	//

	//
	// module groups: start
	//

	moduleGroup := e.Group("/modules")
	//moduleGroup.Use(jwtMiddleWare,JwtIsAccessTokenMiddleware)

	dnaGroup := moduleGroup.Group("/dna")

	dnaGroup.POST("/:assembly", dnaroutes.DNARoute)

	dnaGroup.POST("/assemblies", dnaroutes.AssembliesRoute)

	genesGroup := moduleGroup.Group("/genes")

	genesGroup.POST("/assemblies", generoutes.AssembliesRoute)

	genesGroup.POST("/within/:assembly", generoutes.WithinGenesRoute)

	genesGroup.POST("/closest/:assembly", generoutes.ClosestGeneRoute)

	genesGroup.POST("/annotate/:assembly", generoutes.AnnotateRoute)

	// mutationsGroup := moduleGroup.Group("/mutations",
	// 	jwtMiddleWare,
	// 	JwtIsAccessTokenMiddleware,
	// 	NewJwtPermissionsMiddleware("rdf"))

	mutationsGroup := moduleGroup.Group("/mutations")

	mutationsGroup.POST("/datasets/:assembly", mutationroutes.MutationDatasetsRoute)

	mutationsGroup.POST("/:assembly/:name", mutationroutes.MutationsRoute)

	mutationsGroup.POST("/maf/:assembly", mutationroutes.PileupRoute)

	mutationsGroup.POST("/pileup/:assembly", mutationroutes.PileupRoute)

	gexGroup := moduleGroup.Group("/gex")

	gexGroup.GET("/platforms", gexroutes.PlaformsRoute)

	gexGroup.POST("/types", gexroutes.GexValueTypesRoute)

	gexGroup.POST("/datasets", gexroutes.GexDatasetsRoute)

	gexGroup.POST("/exp", gexroutes.GexGeneExpRoute)

	geneConvGroup := moduleGroup.Group("/geneconv")

	geneConvGroup.POST("/convert/:from/:to", geneconvroutes.ConvertRoute)

	// geneConvGroup.POST("/:species", func(c echo.Context) error {
	// 	return geneconvroutes.GeneInfoRoute(c, "")
	// })

	motifToGeneGroup := moduleGroup.Group("/motiftogene")

	motifToGeneGroup.POST("/convert", motiftogeneroutes.ConvertRoute)

	pathwayGroup := moduleGroup.Group("/pathway")

	pathwayGroup.POST("/datasets", pathwayroutes.DatasetsRoute)

	pathwayGroup.POST("/overlap", pathwayroutes.PathwayOverlapRoute)

	//
	// module groups: end
	//

	httpPort := os.Getenv("PORT")
	if httpPort == "" {
		httpPort = "8080"
	}

	e.Logger.Fatal(e.Start(":" + httpPort))
}
