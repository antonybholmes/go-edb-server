package routes

import (
	"fmt"
	"net/http"

	"github.com/antonybholmes/go-dna"
	"github.com/labstack/echo/v4"
)

type DNA struct {
	Assembly string        `json:"assembly"`
	Location *dna.Location `json:"location"`
	DNA      string        `json:"dna"`
}

type DNAResponse struct {
	Status int    `json:"status"`
	Data   []*DNA `json:"data"`
}

type DNAQuery struct {
	Rev        bool
	Comp       bool
	Format     string
	RepeatMask string
}

func ParseLocationsFromPost(c echo.Context) ([]dna.Location, error) {
	var err error
	locs := new(ReqLocs)

	err = c.Bind(locs)

	if err != nil {
		return nil, err
	}

	return locs.Locations, nil
}

func DNARoute(c echo.Context, dnadbcache *dna.DNADbCache) error {
	locations, err := ParseLocationsFromPost(c)

	if err != nil {
		return c.JSON(http.StatusBadRequest, StatusMessage{Status: http.StatusBadRequest, Message: err.Error()})
	}

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

	data := []*DNA{}

	for _, location := range locations {
		dna, err := dnadb.DNA(&location, query.Rev, query.Comp)

		if err != nil {
			return c.JSON(http.StatusBadRequest, StatusMessage{Status: http.StatusBadRequest, Message: fmt.Sprintf("%s is not a valid chromosome", location.Chr)})
		}

		data = append(data, &DNA{Assembly: assembly, Location: &location, DNA: dna})
	}

	//c.Logger().Debugf("%s", dna)

	return c.JSON(http.StatusOK, DNAResponse{Status: http.StatusOK, Data: data})
}
