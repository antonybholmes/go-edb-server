package pathwayroutes

import (
	"github.com/antonybholmes/go-edb-api/routes"
	geneconv "github.com/antonybholmes/go-geneconv"
	"github.com/labstack/echo/v4"
)

type ReqParams struct {
	Genes []string `json:"genes"`
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

func PathwayRoute(c echo.Context) error {

	// if there is no conversion, just use the regular gene info
	// if fromSpecies == toSpecies {
	// 	return GeneInfoRoute(c, fromSpecies)
	// }

	params, err := ParseParamsFromPost(c)

	if err != nil {
		return routes.ErrorReq(err)
	}

	var ret geneconv.ConversionResults

	//ret.Conversions = make([]geneconv.Conversion, len(params.Searches))

	// for _, search := range params.Genes {

	// 	// Don't care about the errors, just plug empty list into failures
	// 	//conversion, _ := geneconvdbcache.Convert(search, fromSpecies, toSpecies, params.Exact)

	// 	ret.Conversions = append(ret.Conversions, conversion)
	// }

	return routes.MakeDataResp(c, "", ret)

	//return routes.MakeDataResp(c, "", mutationdbcache.GetInstance().List())
}
