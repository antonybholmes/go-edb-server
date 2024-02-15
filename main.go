package main

import (
	"io"
	"net/http"
	"os"
	"runtime"
	"strings"

	"github.com/antonybholmes/go-dna"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/antonybholmes/go-auth"
	"github.com/antonybholmes/go-edb-api/routes"
	"github.com/antonybholmes/go-loctogene"
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
		log.Panic().Msgf("Error loading .env file")
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
	f := os.Getenv("LOG_FILE")
	if f == "" {
		f = "logs/app.log"
	}

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

	dnadbcache := dna.NewDNADbCache("data/dna")
	loctogenedbcache := loctogene.NewLoctogeneDbCache("data/loctogene")

	e.GET("/about", func(c echo.Context) error {
		return c.JSON(http.StatusOK, AboutResp{Name: "go-edb-api", Version: "1.0.0", Copyright: "Copyright (C) 2024 Antony Holmes", Arch: runtime.GOARCH})
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
			return routes.DNARoute(c, dnadbcache)
		})

		e.POST("/genes/within/:assembly", func(c echo.Context) error {
			return routes.WithinGenesRoute(c, loctogenedbcache)
		})

		e.POST("/genes/closest/:assembly", func(c echo.Context) error {
			return routes.ClosestGeneRoute(c, loctogenedbcache)
		})

		e.POST("/annotate/:assembly", func(c echo.Context) error {
			return routes.AnnotationRoute(c, loctogenedbcache)
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
		return routes.DNARoute(c, dnadbcache)
	})

	r.POST("/genes/within/:assembly", func(c echo.Context) error {
		return routes.WithinGenesRoute(c, loctogenedbcache)
	})

	r.POST("/genes/closest/:assembly", func(c echo.Context) error {
		return routes.ClosestGeneRoute(c, loctogenedbcache)
	})

	r.POST("/annotate/:assembly", func(c echo.Context) error {
		return routes.AnnotationRoute(c, loctogenedbcache)
	})

	httpPort := os.Getenv("PORT")
	if httpPort == "" {
		httpPort = "8080"
	}

	e.Logger.Fatal(e.Start(":" + httpPort))
}
