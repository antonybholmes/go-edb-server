package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"runtime"
	"strings"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/xyproto/randomstring"

	"github.com/antonybholmes/go-auth"
	"github.com/antonybholmes/go-auth/email"
	"github.com/antonybholmes/go-dna/dnadbcache"
	"github.com/antonybholmes/go-edb-api/consts"
	"github.com/antonybholmes/go-edb-api/routes"
	"github.com/antonybholmes/go-env"
	"github.com/antonybholmes/go-loctogene/loctogenedbcache"
	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	_ "github.com/mattn/go-sqlite3"
)

type AboutResp struct {
	Name      string `json:"name"`
	Copyright string `json:"copyright"`
	Version   string `json:"version"`
	Arch      string `json:"arch"`
}

func main() {
	err := godotenv.Load()

	if err != nil {
		log.Error().Msgf("Error loading .env file")
	}

	secret := os.Getenv("JWT_SECRET")
	buildMode := os.Getenv("BUILD")

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
	loctogenedbcache.Dir("data/loctogene")

	e.GET("/about", func(c echo.Context) error {
		return c.JSON(http.StatusOK, AboutResp{Name: consts.NAME, Version: consts.VERSION, Copyright: consts.COPYRIGHT, Arch: runtime.GOARCH})
	})

	e.POST("/register", func(c echo.Context) error {
		return routes.RegisterRoute(c, userdb, secret)
	})

	e.POST("/login", func(c echo.Context) error {
		return routes.LoginRoute(c, userdb, secret)
	})

	// Keep some routes for testing purposes during dev
	if strings.Contains(buildMode, "prod") {
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

	r := e.Group("/auth")

	// Configure middleware with the custom claims type
	config := echojwt.Config{
		NewClaimsFunc: func(c echo.Context) jwt.Claims {
			return new(routes.JwtCustomClaims)
		},
		SigningKey: []byte(secret),
	}
	r.Use(echojwt.WithConfig(config))
	r.GET("/info", routes.JWTInfoRoute)

	r.POST("/dna/:assembly", func(c echo.Context) error {
		return routes.DNARoute(c)
	})

	r.POST("/genes/within/:assembly", func(c echo.Context) error {
		return routes.WithinGenesRoute(c)
	})

	r.POST("/genes/closest/:assembly", func(c echo.Context) error {
		return routes.ClosestGeneRoute(c)
	})

	r.POST("/annotate/:assembly", func(c echo.Context) error {
		return routes.AnnotationRoute(c)
	})

	httpPort := os.Getenv("PORT")
	if httpPort == "" {
		httpPort = "8080"
	}

	randomstring.Seed()

	email.SetName(os.Getenv("NAME")).SetUser(env.GetStr("SMTP_USER", ""), env.GetStr("SMTP_PASSWORD", "")).SetHost(env.GetStr("SMTP_HOST", ""), env.GetUint32("SMTP_PORT", 587)).SetFrom(env.GetStr("SMTP_FROM", ""))

	log.Debug().Msgf("dd %s", email.From())
	log.Debug().Msgf("dd %s", env.GetStr("SMTP_FROM", ""))

	code := randomstring.CookieFriendlyString(32)

	err = email.Compose("antony@antonyholmes.dev", "OTP code", fmt.Sprintf("Your one time code is: %s", code))

	if err != nil {
		log.Error().Msgf("%s", err)
	}

	e.Logger.Fatal(e.Start(":" + httpPort))
}
