package routes

import (
	"fmt"

	"github.com/antonybholmes/go-loctogene"
	"github.com/antonybholmes/go-utils"
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
		return utils.MakeBadResp(c, err)
	}

	query, err := ParseGeneQuery(c, c.Param("assembly"), loctogenedbcache)

	if err != nil {
		return utils.MakeBadResp(c, err)
	}

	data := []*loctogene.GenomicFeatures{}

	for _, location := range locations {
		genes, err := query.Db.WithinGenes(&location, query.Level)

		if err != nil {
			return utils.MakeBadResp(c, err)
		}

		data = append(data, genes)
	}

	return utils.MakeDataResp(c, &data)
}

func ClosestGeneRoute(c echo.Context, loctogenedbcache *loctogene.LoctogeneDbCache) error {
	locations, err := ParseLocationsFromPost(c)

	if err != nil {
		return utils.MakeBadResp(c, err)
	}

	query, err := ParseGeneQuery(c, c.Param("assembly"), loctogenedbcache)

	if err != nil {
		return utils.MakeBadResp(c, err)
	}

	n := utils.ParseN(c, DEFAULT_CLOSEST_N)

	data := []*loctogene.GenomicFeatures{}

	for _, location := range locations {
		genes, err := query.Db.ClosestGenes(&location, n, query.Level)

		if err != nil {
			return utils.MakeBadResp(c, err)
		}

		data = append(data, genes)
	}

	return utils.MakeDataResp(c, &data)
}
