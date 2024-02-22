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
	authroutes "github.com/antonybholmes/go-edb-api/routes/auth"
	modroutes "github.com/antonybholmes/go-edb-api/routes/modules"
	"github.com/antonybholmes/go-genes/genedbcache"
	"github.com/antonybholmes/go-sys/env"
	"github.com/golang-jwt/jwt/v5"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	_ "github.com/mattn/go-sqlite3"
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
	e.Use(middleware.CORS())
	//e.Logger.SetLevel(log.DEBUG)

	err = userdb.Init("data/users.db")

	if err != nil {
		log.Fatal().Msgf("Error loading user db: %s", err)
	}

	dnadbcache.Init("data/dna")
	genedbcache.Init("data/genes")

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

	loginGroup := authGroup.Group("/login")

	loginGroup.POST("/email", func(c echo.Context) error {
		return authroutes.EmailPasswordLoginRoute(c)
	})

	loginGroup.POST("/username", func(c echo.Context) error {
		return authroutes.UsernamePasswordLoginRoute(c)
	})

	authGroup.POST("/login", func(c echo.Context) error {
		return authroutes.EmailPasswordLoginRoute(c)
	})

	authGroup.POST("/verify", func(c echo.Context) error {
		return authroutes.EmailVerificationRoute(c)
	}, jwtMiddleWare)

	authGroup.POST("/info", func(c echo.Context) error {
		return authroutes.UserInfoRoute(c)
	}, jwtMiddleWare)

	passwordGroup := authGroup.Group("/password")

	passwordGroup.POST("/reset", func(c echo.Context) error {
		return authroutes.ResetPasswordEmailRoute(c)
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

	//
	// passwordless groups: end
	//

	//
	// token groups: start
	//

	tokenGroup := e.Group("/tokens")
	tokenGroup.POST("/info", authroutes.TokenInfoRoute)

	tokenAuthGroup := tokenGroup.Group("")
	tokenAuthGroup.Use(jwtMiddleWare)
	tokenAuthGroup.POST("/access", func(c echo.Context) error {
		return authroutes.NewAccessTokenRoute(c)
	})

	//
	// token groups: end
	//

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
