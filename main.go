package main

import (
	"net/http"
	"os"
	"runtime"
	"strings"

	"github.com/antonybholmes/go-dna"
	"github.com/antonybholmes/go-edb-api/routes"
	"github.com/antonybholmes/go-loctogene"
	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
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
		log.Debug("Error loading .env file")
	}

	secret := os.Getenv("JWT_SECRET")

	buildMode := os.Getenv("BUILD")

	e := echo.New()

	e.Use(middleware.Logger())
	//e.Use(loggerMiddleware)
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())
	e.Logger.SetLevel(log.DEBUG)

	userdb, err := routes.NewUserDb("data/users.db")

	if err != nil {
		log.Fatalf("Error loading user db: %s", err)
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
