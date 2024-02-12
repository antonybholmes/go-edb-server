package routes

import (
	"fmt"
	"net/http"

	"github.com/antonybholmes/go-loctogene"
	"github.com/labstack/echo/v4"
)

// A GeneQuery contains info from query params.
type GeneQuery struct {
	Level    loctogene.Level
	Db       *loctogene.LoctogeneDb
	Assembly string
}

type GenesResponse struct {
	Status int                          `json:"status"`
	Data   []*loctogene.GenomicFeatures `json:"data"`
}

func ParseGeneQuery(c echo.Context, assembly string, loctogenedbcache *loctogene.LoctogeneDbCache) (*GeneQuery, error) {

	level := loctogene.Gene

	v := c.QueryParam("level")

	if v != "" {
		level = loctogene.ParseLevel(v)
	}

	db, err := loctogenedbcache.Db(assembly)

	if err != nil {
		return nil, fmt.Errorf("unable to open database for assembly %s %s", assembly, err)
	}

	return &GeneQuery{Assembly: assembly, Db: db, Level: level}, nil
}

func WithinGenesRoute(c echo.Context, loctogenedbcache *loctogene.LoctogeneDbCache) error {
	locations, err := ParseLocationsFromPost(c)

	if err != nil {
		return c.JSON(http.StatusBadRequest, StatusMessage{Status: http.StatusBadRequest, Message: err.Error()})
	}

	query, err := ParseGeneQuery(c, c.Param("assembly"), loctogenedbcache)

	if err != nil {
		return c.JSON(http.StatusBadRequest, StatusMessage{Status: http.StatusBadRequest, Message: err.Error()})
	}

	data := []*loctogene.GenomicFeatures{}

	for _, location := range locations {
		genes, err := query.Db.WithinGenes(&location, query.Level)

		if err != nil {
			return c.JSON(http.StatusBadRequest, StatusMessage{Status: http.StatusBadRequest, Message: "there was an error with the database query"})
		}

		data = append(data, genes)
	}

	return c.JSON(http.StatusOK, GenesResponse{Status: http.StatusOK, Data: data})
}

func ClosestGeneRoute(c echo.Context, loctogenedbcache *loctogene.LoctogeneDbCache) error {
	locations, err := ParseLocationsFromPost(c)

	if err != nil {
		return c.JSON(http.StatusBadRequest, StatusMessage{Status: http.StatusBadRequest, Message: err.Error()})
	}

	query, err := ParseGeneQuery(c, c.Param("assembly"), loctogenedbcache)

	if err != nil {
		return c.JSON(http.StatusBadRequest, StatusMessage{Status: http.StatusBadRequest, Message: err.Error()})
	}

	n := ParseN(c)

	data := []*loctogene.GenomicFeatures{}

	for _, location := range locations {

		genes, err := query.Db.ClosestGenes(&location, n, query.Level)

		if err != nil {
			return c.JSON(http.StatusBadRequest, StatusMessage{Status: http.StatusBadRequest, Message: "there was an error with the database query"})
		}

		data = append(data, genes)
	}

	return c.JSON(http.StatusOK, GenesResponse{Status: http.StatusOK, Data: data})
}
