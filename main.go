package main

import (
	"io"
	"net/http"
	"net/mail"
	"os"
	"runtime"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/antonybholmes/go-auth"
	"github.com/antonybholmes/go-auth/userdb"
	"github.com/antonybholmes/go-dna/dnadbcache"
	"github.com/antonybholmes/go-edb-api/consts"
	"github.com/antonybholmes/go-edb-api/routes"
	authroutes "github.com/antonybholmes/go-edb-api/routes/auth"
	modroutes "github.com/antonybholmes/go-edb-api/routes/modules"
	"github.com/antonybholmes/go-genes/genedbcache"
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

func init() {
	var err error
	store, err = sqlitestore.NewSqliteStore("./data/users.db", "sessions", "/", 3600, []byte(consts.SESSION_SECRET))
	if err != nil {
		panic(err)
	}

	dnadbcache.Init("data/dna")
	genedbcache.Init("data/genes")
}

func main() {
	err := env.Load()

	if err != nil {
		log.Error().Msgf("Error loading .env file")
	}

	// list env to see what is loaded
	env.Ls()

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
		AllowMethods:     []string{http.MethodGet, http.MethodPut, http.MethodPost, http.MethodDelete},
		AllowCredentials: true,
	}))

	//e.Logger.SetLevel(log.DEBUG)

	err = userdb.Init("data/users.db")

	if err != nil {
		log.Fatal().Msgf("Error loading user db: %s", err)
	}

	e.GET("/write", func(c echo.Context) error {
		sess, err := session.Get("session", c)
		if err != nil {
			return err
		}
		sess.Options = authroutes.SESSION_OPT_24H
		sess.Values["name"] = "Steve"
		sess.Save(c.Request(), c.Response())

		return c.NoContent(http.StatusOK)
	})

	e.POST("/login", func(c echo.Context) error {
		validator, err := routes.NewValidator(c).ReqBind().Ok()

		if err != nil {
			return err
		}

		if validator.Req.Password == "" {
			return routes.BadReq("empty password: use passwordless")
		}

		authUser, err := userdb.FindUserByUsername(validator.Req.Username)

		if err != nil {
			email, err := mail.ParseAddress(validator.Req.Username)

			if err != nil {
				return routes.BadReq("email address not valid")
			}

			// also check if username is valid email and try to login
			// with that
			authUser, err = userdb.FindUserByEmail(email)

			if err != nil {
				return routes.BadReq("user does not exist")
			}
		}

		if !authUser.EmailVerified {
			return routes.BadReq("email address not verified")
		}

		if !authUser.CanLogin {
			return routes.BadReq("user not allowed to login")
		}

		if !authUser.CheckPasswords(validator.Req.Password) {
			return routes.BadReq("incorrect password")
		}

		sess, err := session.Get("session", c)
		if err != nil {
			return err
		}
		sess.Options = authroutes.SESSION_OPT_24H
		sess.Values["uuid"] = authUser.Uuid
		sess.Save(c.Request(), c.Response())

		return c.NoContent(http.StatusOK)
	})

	e.GET("/read", func(c echo.Context) error {
		sess, err := session.Get("session", c)
		if err != nil {
			return err
		}
		log.Debug().Msgf("%s", sess.ID)

		return c.JSON(http.StatusOK, sess.Values[authroutes.SESSION_UUID])
	})

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
		SigningKey: []byte(consts.JWT_SECRET),
	}
	jwtMiddleWare := echojwt.WithConfig(config)

	//
	// user groups: start
	//

	authGroup := e.Group("/auth")

	authGroup.POST("/signup", func(c echo.Context) error {
		return authroutes.SignupRoute(c)
	})

	authGroup.POST("/login", func(c echo.Context) error {
		return authroutes.UsernamePasswordLoginRoute(c)
	})

	authGroup.POST("/verify", func(c echo.Context) error {
		return authroutes.EmailVerificationRoute(c)
	}, jwtMiddleWare)

	passwordGroup := authGroup.Group("/passwords")

	passwordGroup.POST("/reset", func(c echo.Context) error {
		return authroutes.ResetPasswordFromUsernameRoute(c)
	})

	passwordGroup.POST("/update", func(c echo.Context) error {
		return authroutes.UpdatePasswordRoute(c)
	}, jwtMiddleWare)

	passwordlessGroup := authGroup.Group("/passwordless")

	passwordlessGroup.POST("/email", func(c echo.Context) error {
		return authroutes.PasswordlessEmailRoute(c)
	})

	passwordlessGroup.POST("/login", func(c echo.Context) error {
		return authroutes.PasswordlessLoginRoute(c)
	}, jwtMiddleWare)

	tokenGroup := authGroup.Group("/tokens")
	tokenGroup.Use(jwtMiddleWare)
	tokenGroup.POST("/info", authroutes.TokenInfoRoute)
	tokenGroup.POST("/access", authroutes.NewAccessTokenRoute)

	//
	// Deal with logins where we want a session
	//

	sessionGroup := authGroup.Group("/sessions")

	sessionGroup.POST("/login", func(c echo.Context) error {
		return authroutes.SessionUsernamePasswordLoginRoute(c)
	})

	sessionGroup.POST("/tokens/access", authroutes.SessionNewAccessTokenRoute)

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

	accountsGroup := usersGroup.Group("/accounts")

	accountsGroup.POST("/names/update", func(c echo.Context) error {
		return authroutes.UpdateNameRoute(c)
	})

	//
	// module groups: start
	//

	moduleGroup := e.Group("/modules")
	moduleGroup.Use(jwtMiddleWare)
	moduleGroup.Use(JwtIsAccessTokenMiddleware)

	dnaGroup := moduleGroup.Group("/dna")

	dnaGroup.POST("/:assembly", func(c echo.Context) error {
		return modroutes.DNARoute(c)
	})

	genesGroup := moduleGroup.Group("/genes")

	genesGroup.POST("/within/:assembly", func(c echo.Context) error {
		return modroutes.WithinGenesRoute(c)
	})

	genesGroup.POST("/closest/:assembly", func(c echo.Context) error {
		return modroutes.ClosestGeneRoute(c)
	})

	genesGroup.POST("/annotation/:assembly", func(c echo.Context) error {
		return modroutes.AnnotationRoute(c)
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
