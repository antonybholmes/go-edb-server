package routes

import (
	"net/http"

	"github.com/antonybholmes/go-loctogene"
	"github.com/labstack/echo/v4"
)

type GenesResponse struct {
	Status int                        `json:"status"`
	Data   *loctogene.GenomicFeatures `json:"data"`
}

func WithinGenesRoute(c echo.Context, loctogenedbcache *loctogene.LoctogeneDbCache) error {
	query, err := ParseGeneQuery(c, c.Param("assembly"), loctogenedbcache)

	if err != nil {
		return c.JSON(http.StatusBadRequest, StatusMessage{Status: http.StatusBadRequest, Message: err.Error()})
	}

	genes, err := query.Db.WithinGenes(query.Loc, query.Level)

	if err != nil {
		return c.JSON(http.StatusBadRequest, StatusMessage{Status: http.StatusBadRequest, Message: "there was an error with the database query"})
	}

	return c.JSON(http.StatusOK, GenesResponse{Status: http.StatusOK, Data: genes})
}

func ClosestGeneRoute(c echo.Context, loctogenedbcache *loctogene.LoctogeneDbCache) error {

	query, err := ParseGeneQuery(c, c.Param("assembly"), loctogenedbcache)

	if err != nil {
		return c.JSON(http.StatusBadRequest, StatusMessage{Status: http.StatusBadRequest, Message: err.Error()})
	}

	n := ParseN(c)

	genes, err := query.Db.ClosestGenes(query.Loc, n, query.Level)

	if err != nil {
		return c.JSON(http.StatusBadRequest, StatusMessage{Status: http.StatusBadRequest, Message: "there was an error with the database query"})
	}

	return c.JSON(http.StatusOK, GenesResponse{Status: http.StatusOK, Data: genes})
}
