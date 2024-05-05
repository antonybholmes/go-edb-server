package mutationroutes

import (
	"github.com/antonybholmes/go-edb-api/routes"
	"github.com/antonybholmes/go-microarray"
	"github.com/antonybholmes/go-microarray/microarraydb"
	"github.com/labstack/echo/v4"
)

func ParseSamplesFromPost(c echo.Context) (*microarray.MicroarraySamplesReq, error) {
	var err error
	locs := new(microarray.MicroarraySamplesReq)

	err = c.Bind(locs)

	if err != nil {
		return nil, err
	}

	return locs, nil
}

func MicroarrayExpressionRoute(c echo.Context) error {
	return routes.NewValidator(c).CheckIsValidAccessToken().Success(func(validator *routes.Validator) error {
		samples, err := ParseSamplesFromPost(c)

		if err != nil {
			return routes.ErrorReq(err)
		}

		data, err := microarraydb.Expression(samples)

		if err != nil {
			return routes.ErrorReq(err)
		}

		return routes.MakeDataResp(c, "", data)
	})
}
