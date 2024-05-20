package mutationroutes

import (
	"github.com/antonybholmes/go-edb-api/routes"
	"github.com/antonybholmes/go-edb-api/routes/modroutes/dnaroutes"
	"github.com/antonybholmes/go-microarray"
	"github.com/antonybholmes/go-mutations"
	"github.com/antonybholmes/go-mutations/mutationdbcache"
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

func MafRoute(c echo.Context) error {
	return routes.NewValidator(c).Success(func(validator *routes.Validator) error {
		locations, err := dnaroutes.ParseLocationsFromPost(c)

		if err != nil {
			return routes.ErrorReq(err)
		}

		assembly := c.Param("assembly")
		name := c.Param("name")

		db, err := mutationdbcache.MutationDB(assembly, name)

		if err != nil {
			return routes.ErrorReq(err)
		}

		ret := make([]*mutations.MutationResults, 0, len(locations))

		for _, location := range locations {
			mutations, err := db.FindMutations(&location)

			if err != nil {
				return routes.ErrorReq(err)
			}

			ret = append(ret, mutations)
		}

		return routes.MakeDataResp(c, "", ret)
	})

	//return routes.MakeDataResp(c, "", mutationdbcache.GetInstance().List())
}
