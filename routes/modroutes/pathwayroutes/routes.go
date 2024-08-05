package pathwayroutes

import (
	"github.com/antonybholmes/go-edb-api/routes"
	pathway "github.com/antonybholmes/go-pathway"
	"github.com/antonybholmes/go-pathway/pathwaydbcache"
	"github.com/labstack/echo/v4"
)

type ReqParams struct {
	Geneset  pathway.Geneset `json:"geneset"`
	Datasets []string        `json:"datasets"`
}

func ParseParamsFromPost(c echo.Context) (*ReqParams, error) {

	params := new(ReqParams)

	err := c.Bind(params)

	if err != nil {
		return nil, err
	}

	return params, nil
}

// If species is the empty string, species will be determined
// from the url parameters
// func GeneInfoRoute(c echo.Context, species string) error {
// 	if species == "" {
// 		species = c.Param("species")
// 	}

// 	params, err := ParseParamsFromPost(c)

// 	if err != nil {
// 		return routes.ErrorReq(err)
// 	}

// 	ret := make([]geneconv.Conversion, len(params.Searches))

// 	for ni, search := range params.Searches {

// 		genes, _ := geneconvdb.GeneInfo(search, species, params.Exact)

// 		ret[ni] = geneconv.Conversion{Search: search, Genes: genes}
// 	}

// 	return routes.MakeDataResp(c, "", ret)
// }

func DatasetsRoute(c echo.Context) error {

	datasets, err := pathwaydbcache.Datasets()

	if err != nil {
		return routes.ErrorReq(err)
	}

	return routes.MakeDataPrettyResp(c, "", datasets)
}

func PathwayOverlapRoute(c echo.Context) error {

	params, err := ParseParamsFromPost(c)

	if err != nil {
		return routes.ErrorReq(err)
	}

	testPathway := params.Geneset.ToPathway()

	tests, err := pathwaydbcache.Overlap(testPathway, params.Datasets)

	if err != nil {
		return routes.ErrorReq(err)
	}

	//var ret geneconv.ConversionResults

	//ret.Conversions = make([]geneconv.Conversion, len(params.Searches))

	// for _, search := range params.Genes {

	// 	// Don't care about the errors, just plug empty list into failures
	// 	//conversion, _ := geneconvdbcache.Convert(search, fromSpecies, toSpecies, params.Exact)

	// 	ret.Conversions = append(ret.Conversions, conversion)
	// }

	return routes.MakeDataResp(c, "", tests, false)

	// return routes.MakeDataResp(c, "", mutationdbcache.GetInstance().List())
}
