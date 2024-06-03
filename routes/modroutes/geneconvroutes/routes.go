package geneconvroutes

import (
	"strings"

	"github.com/antonybholmes/go-edb-api/routes"
	geneconv "github.com/antonybholmes/go-gene-conversion"
	geneconvdb "github.com/antonybholmes/go-gene-conversion/geneconvdb"
	"github.com/labstack/echo/v4"
)

type ReqParams struct {
	Searches []string `json:"searches"`
	Exact    bool     `json:"exact"`
}

func ParseParamsFromPost(c echo.Context) (*ReqParams, error) {

	params := new(ReqParams)

	err := c.Bind(params)

	if err != nil {
		return nil, err
	}

	return params, nil
}

func GenesRoute(c echo.Context) error {
	species := c.Param("species")

	params, err := ParseParamsFromPost(c)

	if err != nil {
		return routes.ErrorReq(err)
	}

	ret := make([]geneconv.Conversion, len(params.Searches))

	for ni, search := range params.Searches {

		genes, _ := geneconvdb.Gene(search, species, params.Exact)

		ret[ni] = geneconv.Conversion{Search: search, Genes: genes}
	}

	return routes.MakeDataResp(c, "", ret)

	//return routes.MakeDataResp(c, "", mutationdbcache.GetInstance().List())
}

func ConvertRoute(c echo.Context) error {
	fromSpecies := c.Param("from")
	toSpecies := c.Param("to")

	params, err := ParseParamsFromPost(c)

	if err != nil {
		return routes.ErrorReq(err)
	}

	var ret geneconv.ConversionResults

	fromSpecies = strings.ToLower(fromSpecies)

	if fromSpecies == geneconv.HUMAN_SPECIES {
		ret.From.TaxId = geneconv.HUMAN_TAXONOMY_ID
		ret.From.Species = geneconv.HUMAN_SPECIES
		ret.To.TaxId = geneconv.MOUSE_TAXONOMY_ID
		ret.To.Species = geneconv.MOUSE_SPECIES
	} else {
		ret.From.TaxId = geneconv.MOUSE_TAXONOMY_ID
		ret.From.Species = geneconv.MOUSE_SPECIES
		ret.To.TaxId = geneconv.HUMAN_TAXONOMY_ID
		ret.To.Species = geneconv.HUMAN_SPECIES
	}

	//ret.Conversions = make([]geneconv.Conversion, len(params.Searches))

	for _, search := range params.Searches {

		// Don't care about the errors, just plug empty list into failures
		conversion, _ := geneconvdb.Convert(search, fromSpecies, toSpecies, params.Exact)

		ret.Conversions = append(ret.Conversions, conversion)
	}

	return routes.MakeDataResp(c, "", ret)

	//return routes.MakeDataResp(c, "", mutationdbcache.GetInstance().List())
}
