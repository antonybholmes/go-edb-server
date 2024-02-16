package main

import (
	"fmt"

	"github.com/antonybholmes/go-loctogene"
	"github.com/antonybholmes/go-loctogene/loctogenedbcache"
	"github.com/labstack/echo/v4"
)

const DEFAULT_LEVEL = loctogene.Gene

const DEFAULT_CLOSEST_N uint16 = 5

// A GeneQuery contains info from query params.
type GeneQuery struct {
	Level    loctogene.Level
	Db       *loctogene.LoctogeneDb
	Assembly string
}

type GenesResponse struct {
	Genes []*loctogene.GenomicFeatures `json:"genes"`
}

func ParseGeneQuery(c echo.Context, assembly string) (*GeneQuery, error) {
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

func WithinGenesRoute(c echo.Context) error {
	locations, err := ParseLocationsFromPost(c)

	if err != nil {
		return BadReq(err)
	}

	query, err := ParseGeneQuery(c, c.Param("assembly"))

	if err != nil {
		return BadReq(err)
	}

	data := []*loctogene.GenomicFeatures{}

	for _, location := range locations {
		genes, err := query.Db.WithinGenes(&location, query.Level)

		if err != nil {
			return BadReq(err)
		}

		data = append(data, genes)
	}

	return MakeDataResp(c, &data)
}

func ClosestGeneRoute(c echo.Context) error {
	locations, err := ParseLocationsFromPost(c)

	if err != nil {
		return BadReq(err)
	}

	query, err := ParseGeneQuery(c, c.Param("assembly"))

	if err != nil {
		return BadReq(err)
	}

	n := ParseN(c, DEFAULT_CLOSEST_N)

	data := []*loctogene.GenomicFeatures{}

	for _, location := range locations {
		genes, err := query.Db.ClosestGenes(&location, n, query.Level)

		if err != nil {
			return BadReq(err)
		}

		data = append(data, genes)
	}

	return MakeDataResp(c, &data)
}
