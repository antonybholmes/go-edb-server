package main

import (
	"fmt"
	"net/http"
	"os"
	"runtime"

	"github.com/antonybholmes/go-dna"
	"github.com/antonybholmes/go-gene"
	"github.com/antonybholmes/go-loctogene"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
	_ "github.com/mattn/go-sqlite3"
)

type StatusMessage struct {
	Message string `json:"message"`
	Status  int    `json:"status"`
}

type DNA struct {
	Assembly string `json:"assembly"`
	Location string `json:"location"`
	DNA      string `json:"dna"`
}

type DNAResponse struct {
	Message string `json:"message"`
	Status  int    `json:"status"`
	Data    *DNA   `json:"data"`
}

type GenesResponse struct {
	Message string                     `json:"message"`
	Status  int                        `json:"status"`
	Data    *loctogene.GenomicFeatures `json:"data"`
}

type AnnotationResponse struct {
	Message string               `json:"message"`
	Status  int                  `json:"status"`
	Data    *gene.GeneAnnotation `json:"data"`
}

type ReqLocs struct {
	Locations []dna.Location `json:"locations"`
}

func main() {
	//zerolog.SetGlobalLevel(zerolog.DebugLevel)

	e := echo.New()

	e.Use(middleware.Logger())
	//e.Use(loggerMiddleware)
	e.Use(middleware.Recover())
	//e.Use(middleware.CORS())
	e.Logger.SetLevel(log.DEBUG)

	e.GET("/", func(c echo.Context) error {
		return c.JSON(http.StatusOK, struct {
			Name    string `json:"name"`
			Version string `json:"version"`
		}{Name: "go-edb-api", Version: "1.0.0"})
	})

	e.GET("/about", func(c echo.Context) error {
		return c.JSON(http.StatusOK, struct {
			Name      string `json:"name"`
			Copyright string `json:"copyright"`
			Version   string `json:"version"`
			Arch      string `json:"arch"`
		}{Name: "go-edb-api", Version: "1.0.0", Copyright: "Copyright (C) 2024 Antony Holmes", Arch: runtime.GOARCH})
	})

	e.GET("dna", func(c echo.Context) error {
		query, err := parseDNAQuery(c)

		if err != nil {
			return c.JSON(http.StatusBadRequest, StatusMessage{Status: http.StatusBadRequest, Message: err.Error()})
		}

		//c.Logger().Debugf("%s %s", query.Loc, query.Dir)

		dnadb := dna.NewDNADB(query.Dir)

		dna, err := dnadb.GetDNA(query.Loc, query.Rev, query.Comp)

		//c.Logger().Debugf("%s", dna)

		if err != nil {
			return c.JSON(http.StatusBadRequest, StatusMessage{Status: http.StatusBadRequest, Message: fmt.Sprintf("%s is not a valid chromosome", query.Loc.Chr)})
		}

		return c.JSON(http.StatusOK, DNAResponse{Status: http.StatusOK, Message: "", Data: &DNA{Assembly: query.Assembly, Location: query.Loc.String(), DNA: dna}})
	})

	e.GET("genes/within/:assembly", func(c echo.Context) error {
		query, err := parseGeneQuery(c, c.Param("assembly"))

		if err != nil {
			return c.JSON(http.StatusBadRequest, StatusMessage{Status: http.StatusBadRequest, Message: err.Error()})
		}

		genes, err := query.DB.WithinGenes(query.Loc, query.Level)

		if err != nil {
			return c.JSON(http.StatusBadRequest, StatusMessage{Status: http.StatusBadRequest, Message: "there was an error with the database query"})
		}

		return c.JSON(http.StatusOK, GenesResponse{Status: http.StatusOK, Message: "", Data: genes})
	})

	e.GET("genes/closest/:assembly", func(c echo.Context) error {

		query, err := parseGeneQuery(c, c.Param("assembly"))

		if err != nil {
			return c.JSON(http.StatusBadRequest, StatusMessage{Status: http.StatusBadRequest, Message: err.Error()})
		}

		n := parseN(c)

		genes, err := query.DB.ClosestGenes(query.Loc, n, query.Level)

		if err != nil {
			return c.JSON(http.StatusBadRequest, StatusMessage{Status: http.StatusBadRequest, Message: "there was an error with the database query"})
		}

		return c.JSON(http.StatusOK, GenesResponse{Status: http.StatusOK, Message: "", Data: genes})
	})

	e.POST("annotation/:assembly", func(c echo.Context) error {
		var err error
		locs := new(ReqLocs)

		err = c.Bind(locs)

		if err != nil {
			return c.JSON(http.StatusBadRequest, StatusMessage{Status: http.StatusBadRequest, Message: err.Error()})
		}

		if err != nil {
			return c.JSON(http.StatusBadRequest, StatusMessage{Status: http.StatusBadRequest, Message: err.Error()})
		}

		locations := locs.Locations

		query, err := parseGeneQuery(c, c.Param("assembly"))

		if err != nil {
			return c.JSON(http.StatusBadRequest, StatusMessage{Status: http.StatusBadRequest, Message: err.Error()})
		}

		n := parseN(c)

		tssRegion := dna.NewTSSRegion(2000, 1000)

		annotationDB := gene.NewAnnotate(query.DB, tssRegion, n)

		annotations, err := annotationDB.Annotate(&locations[0])

		if err != nil {
			c.Logger().Debugf("%s", err)
			return c.JSON(http.StatusBadRequest, StatusMessage{Status: http.StatusBadRequest, Message: "there was an error with the database query"})
		}

		return c.JSON(http.StatusOK, AnnotationResponse{Status: http.StatusOK, Message: "", Data: annotations})
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
