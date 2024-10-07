package pathwayroutes

import (
	"github.com/antonybholmes/go-edb-server/routes"
	pathway "github.com/antonybholmes/go-pathway"
	"github.com/antonybholmes/go-pathway/pathwaydbcache"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
)

type ReqOverlapParams struct {
	Geneset  pathway.Geneset `json:"geneset"`
	Datasets []string        `json:"datasets"`
}

type ReqDatasetParams struct {
	Organization string `json:"organization"`
	Name         string `json:"name"`
}

func ParseOverlapParamsFromPost(c echo.Context) (*ReqOverlapParams, error) {

	var params ReqOverlapParams

	err := c.Bind(&params)

	if err != nil {
		return nil, err
	}

	return &params, nil
}

func ParseDatasetParamsFromPost(c echo.Context) (*ReqDatasetParams, error) {

	var params ReqDatasetParams

	err := c.Bind(&params)

	if err != nil {
		return nil, err
	}

	return &params, nil
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

func GenesRoute(c echo.Context) error {

	return routes.MakeDataPrettyResp(c, "", pathwaydbcache.Genes())
}

func DatasetRoute(c echo.Context) error {

	params, err := ParseDatasetParamsFromPost(c)

	if err != nil {
		return routes.ErrorReq(err)
	}

	log.Debug().Msgf("params %v", params)

	datasets, err := pathwaydbcache.MakePublicDataset(params.Organization, params.Name)

	if err != nil {
		return routes.ErrorReq(err)
	}

	return routes.MakeDataPrettyResp(c, "", datasets)
}

func DatasetsRoute(c echo.Context) error {

	datasets, err := pathwaydbcache.AllDatasetsInfo()

	if err != nil {
		return routes.ErrorReq(err)
	}

	return routes.MakeDataPrettyResp(c, "", datasets)
}

func PathwayOverlapRoute(c echo.Context) error {

	params, err := ParseOverlapParamsFromPost(c)

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
