package modroutes

import (
	"github.com/antonybholmes/go-edb-api/routes"
	"github.com/antonybholmes/go-microarray"
	"github.com/labstack/echo/v4"
)

type ReqSamples struct {
	Samples []string `json:"samples"`
}

func ParseSamplesFromPost(c echo.Context) ([]string, error) {
	var err error
	locs := new(ReqSamples)

	err = c.Bind(locs)

	if err != nil {
		return nil, err
	}

	return locs.Samples, nil
}

func MicroarrayExpressionRoute(c echo.Context) error {
	return routes.NewValidator(c).CheckIsValidAccessToken().Success(func(validator *routes.Validator) error {
		samples, err := ParseSamplesFromPost(c)

		if err != nil {
			return routes.ErrorReq(err)
		}

		mr, err := microarray.NewMicroarrayDb("./data/microarray/hgu133plus2")

		if err != nil {
			return routes.ErrorReq(err)
		}

		data, err := mr.Expression(samples)

		if err != nil {
			return routes.ErrorReq(err)
		}

		return routes.MakeDataResp(c, "", data)
	})
}
