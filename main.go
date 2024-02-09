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
	Message string                 `json:"message"`
	Status  int                    `json:"status"`
	Data    []*gene.GeneAnnotation `json:"data"`
}

type ReqLocs struct {
	Locations []dna.Location `json:"locations"`
}

const MAX_ANNOTATIONS = 1000

func main() {
	//zerolog.SetGlobalLevel(zerolog.DebugLevel)

	e := echo.New()

	e.Use(middleware.Logger())
	//e.Use(loggerMiddleware)
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())
	e.Logger.SetLevel(log.DEBUG)

	dnadbcache := dna.NewDNADbCache("data/dna")
	loctogenedbcache := loctogene.NewLoctogeneDbCache("data/loctogene")

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

	e.GET("/v1/dna/:assembly", func(c echo.Context) error {
		assembly := c.Param("assembly")

		query, err := parseDNAQuery(c)

		if err != nil {
			return c.JSON(http.StatusBadRequest, StatusMessage{Status: http.StatusBadRequest, Message: err.Error()})
		}

		//c.Logger().Debugf("%s %s", query.Loc, query.Dir)

		dnadb, err := dnadbcache.Db(assembly)

		if err != nil {
			return c.JSON(http.StatusBadRequest, StatusMessage{Status: http.StatusBadRequest, Message: err.Error()})
		}

		dna, err := dnadb.DNA(query.Loc, query.Rev, query.Comp, query.Format, query.RepeatMask)

		//c.Logger().Debugf("%s", dna)

		if err != nil {
			return c.JSON(http.StatusBadRequest, StatusMessage{Status: http.StatusBadRequest, Message: fmt.Sprintf("%s is not a valid chromosome", query.Loc.Chr)})
		}

		return c.JSON(http.StatusOK, DNAResponse{Status: http.StatusOK, Message: "", Data: &DNA{Assembly: assembly, Location: query.Loc.String(), DNA: dna}})
	})

	e.GET("/v1/genes/within/:assembly", func(c echo.Context) error {
		query, err := parseGeneQuery(c, c.Param("assembly"), loctogenedbcache)

		if err != nil {
			return c.JSON(http.StatusBadRequest, StatusMessage{Status: http.StatusBadRequest, Message: err.Error()})
		}

		genes, err := query.Db.WithinGenes(query.Loc, query.Level)

		if err != nil {
			return c.JSON(http.StatusBadRequest, StatusMessage{Status: http.StatusBadRequest, Message: "there was an error with the database query"})
		}

		return c.JSON(http.StatusOK, GenesResponse{Status: http.StatusOK, Message: "", Data: genes})
	})

	e.GET("/v1/genes/closest/:assembly", func(c echo.Context) error {

		query, err := parseGeneQuery(c, c.Param("assembly"), loctogenedbcache)

		if err != nil {
			return c.JSON(http.StatusBadRequest, StatusMessage{Status: http.StatusBadRequest, Message: err.Error()})
		}

		n := parseN(c)

		genes, err := query.Db.ClosestGenes(query.Loc, n, query.Level)

		if err != nil {
			return c.JSON(http.StatusBadRequest, StatusMessage{Status: http.StatusBadRequest, Message: "there was an error with the database query"})
		}

		return c.JSON(http.StatusOK, GenesResponse{Status: http.StatusOK, Message: "", Data: genes})
	})

	e.POST("/v1/annotate/:assembly", func(c echo.Context) error {
		var err error
		locs := new(ReqLocs)

		err = c.Bind(locs)

		if err != nil {
			return c.JSON(http.StatusBadRequest, StatusMessage{Status: http.StatusBadRequest, Message: err.Error()})
		}

		if err != nil {
			return c.JSON(http.StatusBadRequest, StatusMessage{Status: http.StatusBadRequest, Message: err.Error()})
		}

		// limit amount of data returned per request to 1000 entries at a time
		locations := locs.Locations[0:MAX_ANNOTATIONS]

		query, err := parseGeneQuery(c, c.Param("assembly"), loctogenedbcache)

		if err != nil {
			return c.JSON(http.StatusBadRequest, StatusMessage{Status: http.StatusBadRequest, Message: err.Error()})
		}

		n := parseN(c)

		tssRegion := parseTSSRegion(c)

		output := parseOutput(c)

		annotationDb := gene.NewAnnotateDb(query.Db, tssRegion, n)

		data := []*gene.GeneAnnotation{}

		for _, location := range locations {

			annotations, err := annotationDb.Annotate(&location)

			if err != nil {
				return c.JSON(http.StatusBadRequest, StatusMessage{Status: http.StatusBadRequest, Message: "there was an error with the database query"})
			}

			data = append(data, annotations)
		}

		if output == "text" {
			tsv, err := makeGeneTable(data, tssRegion)

			if err != nil {
				return c.JSON(http.StatusBadRequest, StatusMessage{Status: http.StatusBadRequest, Message: err.Error()})
			}

			return c.String(http.StatusOK, tsv)
		} else {

			return c.JSON(http.StatusOK, AnnotationResponse{Status: http.StatusOK, Message: "", Data: data})
		}
	})

	httpPort := os.Getenv("PORT")
	if httpPort == "" {
		httpPort = "8080"
	}

	e.Logger.Fatal(e.Start(":" + httpPort))
}
