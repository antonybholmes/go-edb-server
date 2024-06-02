package geneconroutes

import (
	"github.com/antonybholmes/go-edb-api/routes"
	genecon "github.com/antonybholmes/go-gene-conversion"
	"github.com/antonybholmes/go-gene-conversion/genecondb"
	"github.com/antonybholmes/go-mutations/mutationdbcache"
	"github.com/labstack/echo/v4"
)

type ReqParams struct {
	Names []string `json:"names"`
}

func ParseParamsFromPost(c echo.Context) (*ReqParams, error) {

	params := new(ReqParams)

	err := c.Bind(params)

	if err != nil {
		return nil, err
	}

	return params, nil
}

func MutationDatabaseRoutes(c echo.Context) error {
	// return routes.NewValidator(c).CheckIsValidAccessToken().Success(func(validator *routes.Validator) error {
	// 	samples, err := ParseSamplesFromPost(c)

	// 	if err != nil {
	// 		return routes.ErrorReq(err)
	// 	}

	// 	data, err := microarraydb.Expression(samples)

	// 	if err != nil {
	// 		return routes.ErrorReq(err)
	// 	}

	// 	return routes.MakeDataResp(c, "", data)
	// })

	return routes.MakeDataResp(c, "", mutationdbcache.GetInstance().List())
}

func GenesRoute(c echo.Context) error {
	species := c.Param("species")

	params, err := ParseParamsFromPost(c)

	if err != nil {
		return routes.ErrorReq(err)
	}

	ret := make([]*genecon.Gene, len(params.Names))

	for ni, name := range params.Names {

		gene, _ := genecondb.Gene(name, species)

		ret[ni] = gene
	}

	return routes.MakeDataResp(c, "", ret)

	//return routes.MakeDataResp(c, "", mutationdbcache.GetInstance().List())
}

func ConvertRoute(c echo.Context) error {
	species := c.Param("species")

	params, err := ParseParamsFromPost(c)

	if err != nil {
		return routes.ErrorReq(err)
	}

	ret := make([]*genecon.Conversion, len(params.Names))

	for ni, name := range params.Names {

		gene, _ := genecondb.Convert(name, species)

		ret[ni] = gene
	}

	return routes.MakeDataResp(c, "", ret)

	//return routes.MakeDataResp(c, "", mutationdbcache.GetInstance().List())
}
