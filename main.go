package main

import (
	"io"
	"net/http"
	"os"
	"runtime"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/antonybholmes/go-auth"
	"github.com/antonybholmes/go-auth/userdb"
	"github.com/antonybholmes/go-dna/dnadbcache"
	"github.com/antonybholmes/go-edb-api/consts"
	"github.com/antonybholmes/go-edb-api/routes/authroutes"
	"github.com/antonybholmes/go-edb-api/routes/modroutes/dnaroutes"
	"github.com/antonybholmes/go-edb-api/routes/modroutes/geneconvroutes"
	"github.com/antonybholmes/go-edb-api/routes/modroutes/generoutes"
	"github.com/antonybholmes/go-edb-api/routes/modroutes/mutationroutes"
	"github.com/antonybholmes/go-gene-conversion/geneconvdb"
	"github.com/antonybholmes/go-genes/genedbcache"
	"github.com/antonybholmes/go-mailer/mailer"
	"github.com/antonybholmes/go-mutations/mutationdbcache"
	"github.com/antonybholmes/go-sys/env"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo-contrib/session"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	_ "github.com/mattn/go-sqlite3"
	"github.com/michaeljs1990/sqlitestore"
)

type AboutResp struct {
	Name      string `json:"name"`
	Copyright string `json:"copyright"`
	Version   string `json:"version"`
}

type InfoResp struct {
	IpAddr string `json:"ipAddr"`
	Arch   string `json:"arch"`
}

var store *sqlitestore.SqliteStore

func initCache() {
	var err error

	store, err = sqlitestore.NewSqliteStore("./data/users.db", "sessions", "/", 3600, []byte(consts.SESSION_SECRET))

	if err != nil {
		log.Fatal().Msgf("error opening %s", "./data/users.db")
	}

	err = userdb.InitDB("data/users.db")

	if err != nil {
		log.Fatal().Msgf("Error loading user db: %s", err)
	}

	mailer.InitMailer()

	dnadbcache.InitCache("data/modules/dna")
	genedbcache.InitCache("data/modules/genes")
	//microarraydb.InitDB("data/microarray")

	mutationdbcache.InitCache("data/modules/mutations")

	geneconvdb.InitCache("data/modules/geneconv/geneconv.db")
}

func main() {
	consts.LoadConsts()

	//env.Load()

	// list env to see what is loaded
	env.Ls()

	initCache()

	//buildMode := env.GetStr("BUILD", "dev")

	//
	// Set logging to file
	//

	//log.SetOutput(logFile)

	//
	// end logging setup
	//

	e := echo.New()

	//e.Use(middleware.Logger())

	e.Use(session.Middleware(store))

	// write to both stdout and log file
	f := env.GetStr("LOG_FILE", "logs/app.log")

	logFile, err := os.OpenFile(f, os.O_APPEND|os.O_RDWR|os.O_CREATE, 0644)

	if err != nil {
		log.Fatal().Msgf("%s", err)
	}

	defer logFile.Close()

	logger := zerolog.New(io.MultiWriter(os.Stdout, logFile)).With().Timestamp().Logger() //os.Stderr)

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
		AllowOrigins:     []string{"https://edb.rdf-lab.org", "http://localhost:8000"},
		AllowMethods:     []string{http.MethodGet, http.MethodPut, http.MethodPost},
		AllowCredentials: true,
	}))

	//e.Logger.SetLevel(log.DEBUG)

	// e.GET("/write", func(c echo.Context) error {
	// 	sess, err := session.Get("session", c)
	// 	if err != nil {
	// 		return err
	// 	}
	// 	sess.Options = authroutes.SESSION_OPT_24H
	// 	sess.Values["name"] = "Steve"
	// 	sess.Save(c.Request(), c.Response())

	// 	return c.NoContent(http.StatusOK)
	// })

	// e.POST("/login", func(c echo.Context) error {
	// 	validator, err := routes.NewValidator(c).ReqBind().Ok()

	// 	if err != nil {
	// 		return err
	// 	}

	// 	if validator.Req.Password == "" {
	// 		return routes.ErrorReq("empty password: use passwordless")
	// 	}

	// 	authUser, err := userdb.FindUserByUsername(validator.Req.Username)

	// 	if err != nil {
	// 		email, err := mail.ParseAddress(validator.Req.Username)

	// 		if err != nil {
	// 			return routes.ErrorReq("email address not valid")
	// 		}

	// 		// also check if username is valid email and try to login
	// 		// with that
	// 		authUser, err = userdb.FindUserByEmail(email)

	// 		if err != nil {
	// 			return routes.UserDoesNotExistReq()
	// 		}
	// 	}

	// 	if !authUser.EmailVerified {
	// 		return routes.ErrorReq("email address not verified")
	// 	}

	// 	if !authUser.CanSignIn {
	// 		return routes.ErrorReq("user not allowed to login")
	// 	}

	// 	if !authUser.CheckPasswords(validator.Req.Password) {
	// 		return routes.InvalidPasswordReq()
	// 	}

	// 	sess, err := session.Get("session", c)
	// 	if err != nil {
	// 		return err
	// 	}
	// 	sess.Options = authroutes.SESSION_OPT_30D
	// 	sess.Values["uuid"] = authUser.Uuid
	// 	sess.Save(c.Request(), c.Response())

	// 	return c.NoContent(http.StatusOK)
	// })

	// e.GET("/read", func(c echo.Context) error {
	// 	sess, err := session.Get("session", c)
	// 	if err != nil {
	// 		return err
	// 	}
	// 	log.Debug().Msgf("%s", sess.ID)

	// 	return c.JSON(http.StatusOK, sess.Values[authroutes.SESSION_UUID])
	// })

	e.GET("/about", func(c echo.Context) error {
		return c.JSON(http.StatusOK, AboutResp{Name: consts.NAME, Version: consts.VERSION, Copyright: consts.COPYRIGHT})
	})

	e.GET("/info", func(c echo.Context) error {
		return c.JSON(http.StatusOK, InfoResp{Arch: runtime.GOARCH, IpAddr: c.RealIP()})
	})

	// Configure middleware with the custom claims type
	config := echojwt.Config{
		NewClaimsFunc: func(c echo.Context) jwt.Claims {
			return new(auth.JwtCustomClaims)
		},
		SigningKey: consts.JWT_PRIVATE_KEY,
		// Have to tell it to use the public key for verification
		KeyFunc: func(token *jwt.Token) (interface{}, error) {
			return consts.JWT_PUBLIC_KEY, nil
		},
	}
	jwtMiddleWare := echojwt.WithConfig(config)

	//
	// Routes
	//

	e.POST("/signup", func(c echo.Context) error {
		return authroutes.SignupRoute(c)
	})

	//
	// user groups: start
	//

	authGroup := e.Group("/auth")

	authGroup.POST("/signin", func(c echo.Context) error {
		return authroutes.UsernamePasswordSignInRoute(c)
	})

	authGroup.POST("/verify", func(c echo.Context) error {
		return authroutes.EmailAddressWasVerifiedRoute(c)
	}, jwtMiddleWare)

	passwordGroup := authGroup.Group("/passwords")

	passwordGroup.POST("/reset", func(c echo.Context) error {
		return authroutes.SendChangeEmailRoute(c)
	})

	passwordGroup.POST("/update", func(c echo.Context) error {
		return authroutes.UpdatePasswordRoute(c)
	}, jwtMiddleWare)

	emailGroup := authGroup.Group("/email")

	emailGroup.POST("/update", func(c echo.Context) error {
		return authroutes.UpdateEmailRoute(c)
	}, jwtMiddleWare)

	passwordlessGroup := authGroup.Group("/passwordless")

	passwordlessGroup.POST("/email", func(c echo.Context) error {
		return authroutes.PasswordlessEmailRoute(c, nil)
	})

	passwordlessGroup.POST("/signin", func(c echo.Context) error {
		return authroutes.PasswordlessSignInRoute(c)
	}, jwtMiddleWare)

	tokenGroup := authGroup.Group("/tokens")
	tokenGroup.Use(jwtMiddleWare)
	tokenGroup.POST("/info", authroutes.TokenInfoRoute)
	tokenGroup.POST("/access", authroutes.NewAccessTokenRoute)

	//
	// Deal with logins where we want a session
	//

	sessionGroup := e.Group("/sessions")

	sessionAuthGroup := sessionGroup.Group("/auth")

	sessionAuthGroup.POST("/signin", func(c echo.Context) error {
		return authroutes.SessionUsernamePasswordSignInRoute(c)
	})

	sessionAuthGroup.POST("/email/change", func(c echo.Context) error {
		return authroutes.SessionSendChangeEmailRoute(c)
	})

	sessionAuthGroup.POST("/password/reset", func(c echo.Context) error {
		return authroutes.SessionSendResetPasswordRoute(c)
	})

	sessionAuthGroup.POST("/passwordless/signin", func(c echo.Context) error {
		return authroutes.SessionPasswordlessSignInRoute(c)
	}, jwtMiddleWare)

	sessionAuthGroup.POST("/tokens/access", authroutes.SessionNewAccessTokenRoute, SessionIsValidMiddleware)

	sessionUsersGroup := sessionGroup.Group("/users")
	sessionUsersGroup.Use(SessionIsValidMiddleware)

	sessionUsersGroup.POST("/info", func(c echo.Context) error {
		return authroutes.SessionUserInfoRoute(c)
	})

	sessionUsersGroup.POST("/update", func(c echo.Context) error {
		return authroutes.SessionUpdateUserInfoRoute(c)
	})

	// sessionPasswordGroup := sessionAuthGroup.Group("/passwords")
	// sessionPasswordGroup.Use(SessionIsValidMiddleware)

	// sessionPasswordGroup.POST("/update", func(c echo.Context) error {
	// 	return authroutes.SessionUpdatePasswordRoute(c)
	// })

	//
	// sessions: end
	//

	//
	// passwordless groups: end
	//

	usersGroup := e.Group("/users")
	usersGroup.Use(jwtMiddleWare)
	usersGroup.Use(JwtIsAccessTokenMiddleware)

	usersGroup.POST("/info", func(c echo.Context) error {
		return authroutes.UserInfoRoute(c)
	})

	usersGroup.POST("/update", func(c echo.Context) error {
		return authroutes.UpdateAccountRoute(c)
	})

	usersGroup.POST("/passwords/update", func(c echo.Context) error {
		return authroutes.UpdatePasswordRoute(c)
	})

	//
	// module groups: start
	//

	moduleGroup := e.Group("/modules")
	//moduleGroup.Use(jwtMiddleWare)
	//moduleGroup.Use(JwtIsAccessTokenMiddleware)

	dnaGroup := moduleGroup.Group("/dna")

	dnaGroup.POST("/:assembly", func(c echo.Context) error {
		return dnaroutes.DNARoute(c)
	})

	dnaGroup.POST("/assemblies", func(c echo.Context) error {
		return dnaroutes.AssembliesRoute(c)
	})

	genesGroup := moduleGroup.Group("/genes")

	genesGroup.POST("/assemblies", func(c echo.Context) error {
		return generoutes.AssembliesRoute(c)
	})

	genesGroup.POST("/within/:assembly", func(c echo.Context) error {
		return generoutes.WithinGenesRoute(c)
	})

	genesGroup.POST("/closest/:assembly", func(c echo.Context) error {
		return generoutes.ClosestGeneRoute(c)
	})

	genesGroup.POST("/annotate/:assembly", func(c echo.Context) error {
		return generoutes.AnnotateRoute(c)
	})

	mutationsGroup := moduleGroup.Group("/mutations",
		jwtMiddleWare,
		JwtIsAccessTokenMiddleware,
		NewJwtPermissionsMiddleware("GetMutations"))

	mutationsGroup.POST("/databases", func(c echo.Context) error {
		return mutationroutes.MutationDatabasesRoute(c)
	})

	mutationsGroup.POST("/:assembly/:name", func(c echo.Context) error {
		return mutationroutes.MutationsRoute(c)
	})

	mutationsGroup.POST("/maf", func(c echo.Context) error {
		return mutationroutes.PileupRoute(c)
	})

	mutationsGroup.POST("/pileup", func(c echo.Context) error {
		return mutationroutes.PileupRoute(c)
	})

	geneConvGroup := moduleGroup.Group("/geneconv")

	geneConvGroup.POST("/convert/:from/:to", func(c echo.Context) error {
		return geneconvroutes.ConvertRoute(c)
	})

	geneConvGroup.POST("/:species", func(c echo.Context) error {
		return geneconvroutes.GeneInfoRoute(c, "")
	})

	//
	// module groups: end
	//

	httpPort := os.Getenv("PORT")
	if httpPort == "" {
		httpPort = "8080"
	}

	e.Logger.Fatal(e.Start(":" + httpPort))
}
