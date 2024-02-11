package main

import (
	"net/http"
	"os"
	"runtime"

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

func main() {
	err := godotenv.Load()

	if err != nil {
		log.Fatal("Error loading .env file")
	}

	secret := os.Getenv("JWT_SECRET")

	e := echo.New()

	e.Use(middleware.Logger())
	//e.Use(loggerMiddleware)
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())
	e.Logger.SetLevel(log.DEBUG)

	dnadbcache := dna.NewDNADbCache("data/dna")
	loctogenedbcache := loctogene.NewLoctogeneDbCache("data/loctogene")

	e.GET("/about", func(c echo.Context) error {
		return c.JSON(http.StatusOK, struct {
			Name      string `json:"name"`
			Copyright string `json:"copyright"`
			Version   string `json:"version"`
			Arch      string `json:"arch"`
		}{Name: "go-edb-api", Version: "1.0.0", Copyright: "Copyright (C) 2024 Antony Holmes", Arch: runtime.GOARCH})
	})

	e.POST("/login", LoginRoute)

	r := e.Group("/restricted")

	// Configure middleware with the custom claims type
	config := echojwt.Config{
		NewClaimsFunc: func(c echo.Context) jwt.Claims {
			return new(JwtCustomClaims)
		},
		SigningKey: []byte(secret),
	}
	r.Use(echojwt.WithConfig(config))
	r.GET("", RestrictedRoute)

	e.GET("/v1/dna/:assembly", func(c echo.Context) error {
		return routes.DNARoute(c, dnadbcache)
	})

	e.GET("/v1/genes/within/:assembly", func(c echo.Context) error {
		return routes.WithinGenesRoute(c, loctogenedbcache)
	})

	e.GET("/v1/genes/closest/:assembly", func(c echo.Context) error {

		return routes.ClosestGeneRoute(c, loctogenedbcache)
	})

	e.POST("/v1/annotate/:assembly", func(c echo.Context) error {
		return routes.AnnotationRoute(c, loctogenedbcache)
	})

	httpPort := os.Getenv("PORT")
	if httpPort == "" {
		httpPort = "8080"
	}

	e.Logger.Fatal(e.Start(":" + httpPort))
}
