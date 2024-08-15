package gexroutes

import (
	"strconv"

	"github.com/antonybholmes/go-edb-api/routes"
	"github.com/antonybholmes/go-gex/gexdbcache"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
)

type GexParams struct {
	GexType  int      `json:"gexType"`
	Genes    []string `json:"genes"`
	Datasets []int    `json:"datasets"`
}

func ParseParamsFromPost(c echo.Context) (*GexParams, error) {

	var params GexParams

	err := c.Bind(&params)

	if err != nil {
		return nil, err
	}

	return &params, nil
}

func GexTypesRoute(c echo.Context) error {

	types, err := gexdbcache.GexTypes()

	if err != nil {
		return err
	}

	return routes.MakeDataPrettyResp(c, "", types)
}

func GexDatasetsRoute(c echo.Context) error {

	gexType, err := strconv.Atoi(c.QueryParam("gex_type"))

	if err != nil {
		return err
	}

	datasets, err := gexdbcache.Datasets(gexType)

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
		log.Debug().Msgf("e1 %s", err)
		return routes.ErrorReq(err)
	}

	if params.GexType == 2 {
		// microarray
		return routes.ErrorReq("microarray not implemented")
	} else {
		// default to rna-seq
		ret, err := gexdbcache.RNASeqValues(gexGenes, params.Datasets)

		if err != nil {
			log.Debug().Msgf("e2 %s", err)
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
