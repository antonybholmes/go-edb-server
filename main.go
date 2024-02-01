package main

import (
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/antonybholmes/go-dna"
	"github.com/antonybholmes/go-loctogene"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
)

type StatusMessage struct {
	Message string `json:"message"`
	Status  string `json:"status"`
}

func main() {
	//zerolog.SetGlobalLevel(zerolog.DebugLevel)

	e := echo.New()

	e.Use(middleware.Logger())
	//e.Use(loggerMiddleware)
	e.Use(middleware.Recover())
	e.Logger.SetLevel(log.DEBUG)

	modulesDir := os.Getenv("MODULESDIR")
	if modulesDir == "" {
		modulesDir = "data/"
	}

	e.GET("/", func(c echo.Context) error {
		return c.JSON(http.StatusOK, struct {
			Name    string `json:"name"`
			Version string `json:"version"`
		}{Name: "go-edb-api", Version: "1.0.0"})
	})

	e.GET("/about", func(c echo.Context) error {
		return c.JSON(http.StatusOK, struct {
			Name    string `json:"name"`
			Version string `json:"version"`
		}{Name: "go-edb-api", Version: "1.0.0"})
	})

	e.GET("/health", func(c echo.Context) error {
		return c.JSON(http.StatusOK, StatusMessage{Status: "OK"})
	})

	e.GET("/dna", func(c echo.Context) error {

		query, err := parseDNAQuery(c, modulesDir)

		if err != nil {
			return c.JSON(http.StatusBadRequest, StatusMessage{Status: "Error", Message: err.Error()})
		}

		dna, err := dna.GetDNA(query.Dir, query.Loc, query.RevComp)

		if err != nil {
			return c.JSON(http.StatusBadRequest, StatusMessage{Status: "Error", Message: fmt.Sprintf("%s is not a valid chromosome", query.Loc.Chr)})
		}

		return c.JSON(http.StatusOK, dna)
	})

	e.GET("/genes/within", func(c echo.Context) error {
		query, err := parseGeneQuery(c, modulesDir)

		if err != nil {
			return c.JSON(http.StatusBadRequest, StatusMessage{Status: "Error", Message: err.Error()})
		}

		genes, err := loctogene.GetGenesWithin(query.DB, query.Loc, query.Level)

		if err != nil {
			return c.JSON(http.StatusBadRequest, StatusMessage{Status: "Error", Message: err.Error()})
		}

		return c.JSON(http.StatusOK, genes)
	})

	e.GET("/genes/closest", func(c echo.Context) error {

		query, err := parseGeneQuery(c, modulesDir)

		if err != nil {
			return c.JSON(http.StatusBadRequest, StatusMessage{Status: "Error", Message: err.Error()})
		}

		n := DEFAULT_CLOSEST_N

		v := c.QueryParam("n")

		if v != "" {
			n, err = strconv.Atoi(v)

			if err != nil {
				n = DEFAULT_CLOSEST_N
			}
		}

		genes, err := loctogene.GetClosestGenes(query.DB, query.Loc, n, query.Level)

		if err != nil {
			return c.JSON(http.StatusBadRequest, StatusMessage{Status: "Error", Message: err.Error()})
		}

		return c.JSON(http.StatusOK, genes)
	})

	// e.POST("/genes/within", func(c echo.Context) error {
	// 	jsonBody := make(map[string]interface{})

	// 	json.NewDecoder(c.Request().Body).Decode(&jsonBody)

	// 	chr := DEFAULT_CHR
	// 	var ok bool

	// 	_, ok = jsonBody["chr"]

	// 	if ok {
	// 		c.Logger().Info("chr set through body.")
	// 		chr = jsonBody["chr"].(string)
	// 	} else {
	// 		c.Logger().Warn("chr not set through body, using default.")
	// 	}

	// 	start := DEFAULT_START

	// 	_, ok = jsonBody["start"]

	// 	if ok {
	// 		c.Logger().Info("start set through body.")
	// 		start = jsonBody["start"].(int)
	// 	} else {
	// 		c.Logger().Warn(fmt.Sprintf("start not set, using default %d.", DEFAULT_START))
	// 	}

	// 	end := DEFAULT_END

	// 	_, ok = jsonBody["end"]

	// 	if ok {
	// 		c.Logger().Info("end set through body.")
	// 		end = jsonBody["end"].(int)
	// 	} else {
	// 		c.Logger().Warn(fmt.Sprintf("end not set, using default %d.", DEFAULT_END))
	// 	}

	// 	assembly := DEFAULT_ASSEMBLY

	// 	_, ok = jsonBody["assembly"]

	// 	if ok {
	// 		c.Logger().Info("assembly set through body.")
	// 		assembly = jsonBody["assembly"].(string)
	// 	} else {
	// 		c.Logger().Warn(fmt.Sprintf("assembly not set, using default %s.", DEFAULT_ASSEMBLY))
	// 	}

	// 	level := DEFAULT_LEVEL

	// 	_, ok = jsonBody["level"]

	// 	if ok {
	// 		c.Logger().Info("level set through body.")
	// 		level = jsonBody["level"].(int)
	// 	} else {
	// 		c.Logger().Warn(fmt.Sprintf("level not set, using default %d.", DEFAULT_LEVEL))
	// 	}

	// 	c.Logger().Info(fmt.Sprintf("loc: %s:%d-%d on %s", chr, start, end, assembly))

	// 	db, err := loctogene.GetDB(fmt.Sprintf("data/modules/loctogene/%s.db", assembly))

	// 	if err != nil {
	// 		return c.JSON(http.StatusInternalServerError, struct{ Status string }{Status: "DB error"})
	// 	}

	// 	loc := loctogene.Location{Chr: chr, Start: start, End: end}

	// 	genes, err := loctogene.GetGenesWithin(db, &loc, level)

	// 	if err != nil {
	// 		return c.JSON(http.StatusInternalServerError, struct{ Status string }{Status: "Error"})
	// 	}

	// 	return c.JSON(http.StatusOK, genes)
	// })

	httpPort := os.Getenv("PORT")
	if httpPort == "" {
		httpPort = "8080"
	}

	e.Logger.Fatal(e.Start(":" + httpPort))
}
