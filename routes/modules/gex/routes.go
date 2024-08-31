package gexroutes

import (
	"github.com/antonybholmes/go-edb-server/routes"
	"github.com/antonybholmes/go-gex"
	"github.com/antonybholmes/go-gex/gexdbcache"
	"github.com/labstack/echo/v4"
)

type GexParams struct {
	Platform     *gex.Platform     `json:"platform"`
	GexValueType *gex.GexValueType `json:"gexValueType"`
	Genes        []string          `json:"genes"`
	Datasets     []int             `json:"datasets"`
}

func ParseParamsFromPost(c echo.Context) (*GexParams, error) {

	var params GexParams

	err := c.Bind(&params)

	if err != nil {
		return nil, err
	}

	return &params, nil
}

func PlaformsRoute(c echo.Context) error {

	types, err := gexdbcache.Platforms()

	if err != nil {
		return err
	}

	return routes.MakeDataPrettyResp(c, "", types)
}

func GexValueTypesRoute(c echo.Context) error {

	params, err := ParseParamsFromPost(c)

	if err != nil {
		return err
	}

	datasets, err := gexdbcache.GexValueTypes(params.Platform)

	if err != nil {
		return err
	}

	return routes.MakeDataPrettyResp(c, "", datasets)
}

func GexDatasetsRoute(c echo.Context) error {

	params, err := ParseParamsFromPost(c)

	if err != nil {
		return err
	}

	datasets, err := gexdbcache.Datasets(params.Platform)

	if err != nil {
		return err
	}

	return routes.MakeDataPrettyResp(c, "", datasets)
}

func GexGeneExpRoute(c echo.Context) error {
	params, err := ParseParamsFromPost(c)

	if err != nil {
		return routes.ErrorReq(err)
	}

	// convert search genes into actual genes in the database
	gexGenes, err := gexdbcache.GetGenes(params.Genes)

	if err != nil {
		return routes.ErrorReq(err)
	}

	if params.Platform.Id == 2 {
		// microarray
		ret, err := gexdbcache.MicroarrayValues(gexGenes, params.Platform, params.GexValueType, params.Datasets)

		if err != nil {

			return routes.ErrorReq(err)
		}

		return routes.MakeDataPrettyResp(c, "", ret)
	} else {
		// default to rna-seq
		ret, err := gexdbcache.RNASeqValues(gexGenes, params.Platform, params.GexValueType, params.Datasets)

		if err != nil {

			return routes.ErrorReq(err)
		}

		return routes.MakeDataPrettyResp(c, "", ret)
	}
}

// func GexRoute(c echo.Context) error {
// 	gexType := c.Param("type")

// 	params, err := ParseParamsFromPost(c)

// 	if err != nil {
// 		return routes.ErrorReq(err)
// 	}

// 	search, err := gexdbcache.GetInstance().Search(gexType, params.Datasets, params.Genes)

// 	if err != nil {
// 		return routes.ErrorReq(err)
// 	}

// 	return routes.MakeDataPrettyResp(c, "", search)

// 	//return routes.MakeDataResp(c, "", mutationdbcache.GetInstance().List())
// }
