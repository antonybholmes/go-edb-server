package main

import (
	"io"
	"net/http"
	"os"
	"runtime"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/antonybholmes/go-auth"
	"github.com/antonybholmes/go-dna/dnadbcache"
	"github.com/antonybholmes/go-edb-api/consts"
	"github.com/antonybholmes/go-edb-api/routes"

	"github.com/antonybholmes/go-env"
	"github.com/antonybholmes/go-gene/genedbcache"
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

	env.Ls()

	buildMode := env.GetStr("BUILD", "dev")

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

	userdb, err := auth.NewUserDb("data/users.db")

	if err != nil {
		log.Fatal().Msgf("Error loading user db: %s", err)
	}

	dnadbcache.Dir("data/dna")
	genedbcache.Dir("data/genes")

	e.GET("/about", func(c echo.Context) error {
		return c.JSON(http.StatusOK, AboutResp{Name: consts.NAME, Version: consts.VERSION, Copyright: consts.COPYRIGHT})
	})

	e.GET("/info", func(c echo.Context) error {
		return c.JSON(http.StatusOK, InfoResp{Arch: runtime.GOARCH, IpAddr: c.RealIP()})
	})

	group := e.Group("/users")

	group.POST("/signup", func(c echo.Context) error {
		return routes.Signup(c, userdb, consts.JWT_SECRET)
	})

	authGroup := group.Group("/auth")

	// Configure middleware with the custom claims type
	config := echojwt.Config{
		NewClaimsFunc: func(c echo.Context) jwt.Claims {
			return new(auth.JwtOtpCustomClaims)
		},
		SigningKey: []byte(consts.JWT_SECRET),
	}
	authGroup.Use(echojwt.WithConfig(config))
	//authGroup.Use(JwtOtpCheckMiddleware)

	authGroup.POST("/verify", func(c echo.Context) error {
		return routes.Verification(c)
	})

	group.POST("/login", func(c echo.Context) error {
		return routes.LoginRoute(c)
	})

	// Keep some routes for testing purposes during dev
	if buildMode == "dev" {
		e.POST("/dna/:assembly", func(c echo.Context) error {
			return routes.DNARoute(c)
		})

		e.POST("/genes/within/:assembly", func(c echo.Context) error {
			return routes.WithinGenesRoute(c)
		})

		e.POST("/genes/closest/:assembly", func(c echo.Context) error {
			return routes.ClosestGeneRoute(c)
		})

		e.POST("/annotate/:assembly", func(c echo.Context) error {
			return routes.AnnotationRoute(c)
		})
	}

	group = e.Group("/tokens")

	// Configure middleware with the custom claims type
	config = echojwt.Config{
		NewClaimsFunc: func(c echo.Context) jwt.Claims {
			return new(auth.JwtCustomClaims)
		},
		SigningKey: []byte(consts.JWT_SECRET),
	}
	group.Use(echojwt.WithConfig(config))
	//group.Use(JwtCheckMiddleware)

	group.GET("/info", routes.JWTInfoRoute)

	group.POST("/validate", func(c echo.Context) error {
		return routes.ValidateTokenRoute(c)
	})

	group.POST("/refresh", func(c echo.Context) error {
		return routes.RefreshTokenRoute(c)
	})

	group = e.Group("/modules")

	//authGroup.Use(echojwt.WithConfig(config))
	//authGroup.Use(JwtCheckMiddleware)

	group.Use(echojwt.WithConfig(config))
	//group.Use(JwtCheckMiddleware)

	//authGroup = group.Group("/dna")

	authGroup.POST("/:assembly", func(c echo.Context) error {
		return routes.DNARoute(c)
	})

	authGroup = group.Group("/genes")

	authGroup.POST("/within/:assembly", func(c echo.Context) error {
		return routes.WithinGenesRoute(c)
	})

	authGroup.POST("/closest/:assembly", func(c echo.Context) error {
		return routes.ClosestGeneRoute(c)
	})

	authGroup.POST("/annotation/:assembly", func(c echo.Context) error {
		return routes.AnnotationRoute(c)
	})

	httpPort := os.Getenv("PORT")
	if httpPort == "" {
		httpPort = "8080"
	}

	e.Logger.Fatal(e.Start(":" + httpPort))
}
