package routes

import (
	"fmt"
	"net/http"

	"github.com/antonybholmes/go-dna"
	"github.com/labstack/echo/v4"
)

type DNA struct {
	Assembly string `json:"assembly"`
	Location string `json:"location"`
	DNA      string `json:"dna"`
}

type DNAResponse struct {
	Status int  `json:"status"`
	Data   *DNA `json:"data"`
}

type DNAQuery struct {
	Loc        *dna.Location
	Rev        bool
	Comp       bool
	Format     string
	RepeatMask string
}

func DNARoute(c echo.Context, dnadbcache *dna.DNADbCache) error {
	assembly := c.Param("assembly")

	query, err := ParseDNAQuery(c)

	if err != nil {
		return c.JSON(http.StatusBadRequest, StatusMessage{Status: http.StatusBadRequest, Message: err.Error()})
	}

	//c.Logger().Debugf("%s %s", query.Loc, query.Dir)

	dnadb, err := dnadbcache.Db(assembly, query.Format, query.RepeatMask)

	if err != nil {
		return c.JSON(http.StatusBadRequest, StatusMessage{Status: http.StatusBadRequest, Message: err.Error()})
	}

	dna, err := dnadb.DNA(query.Loc, query.Rev, query.Comp)

	//c.Logger().Debugf("%s", dna)

	if err != nil {
		return c.JSON(http.StatusBadRequest, StatusMessage{Status: http.StatusBadRequest, Message: fmt.Sprintf("%s is not a valid chromosome", query.Loc.Chr)})
	}

	return c.JSON(http.StatusOK, DNAResponse{Status: http.StatusOK, Data: &DNA{Assembly: assembly, Location: query.Loc.String(), DNA: dna}})
}
