package motiftogeneroutes

import (
	"github.com/antonybholmes/go-edb-api/routes"
	"github.com/antonybholmes/go-motiftogene"
	"github.com/antonybholmes/go-motiftogene/motiftogenedb"
	"github.com/labstack/echo/v4"
)

type ReqParams struct {
	Searches []string `json:"searches"`
	Exact    bool     `json:"exact"`
}

type MotifToGeneRes struct {
	Search      string                     `json:"search"`
	Conversions []*motiftogene.MotifToGene `json:"conversions"`
}

func ParseParamsFromPost(c echo.Context) (*ReqParams, error) {

	params := new(ReqParams)

	err := c.Bind(params)

	if err != nil {
		return nil, err
	}

	return params, nil
}

func ConvertRoute(c echo.Context) error {

	params, err := ParseParamsFromPost(c)

	if err != nil {
		return routes.ErrorReq(err)
	}

	ret := make([]*MotifToGeneRes, 0, len(params.Searches))

	for _, search := range params.Searches {

		c := MotifToGeneRes{}

		c.Search = search
		c.Conversions = make([]*motiftogene.MotifToGene, 0, 1)

		// Don't care about the errors, just plug empty list into failures
		conversion, err := motiftogenedb.Convert(search)

		if err == nil {
			c.Conversions = append(c.Conversions, conversion)
		}

		ret = append(ret, &c)
	}

	return routes.MakeDataPrettyResp(c, "", ret)

	//return routes.MakeDataResp(c, "", mutationdbcache.GetInstance().List())
}
