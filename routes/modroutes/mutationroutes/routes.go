package mutationroutes

import (
	"sort"

	"github.com/antonybholmes/go-dna"
	"github.com/antonybholmes/go-edb-api/routes"
	"github.com/antonybholmes/go-edb-api/routes/modroutes/dnaroutes"
	"github.com/antonybholmes/go-mutations"
	"github.com/antonybholmes/go-mutations/mutationdbcache"
	"github.com/labstack/echo/v4"
)

type MutationParams struct {
	Locations []*dna.Location
	DBs       []string
}

type ReqMutationParams struct {
	Locations []string `json:"locations"`
	DBs       []string `json:"databases"`
}

func ParseParamsFromPost(c echo.Context) (*MutationParams, error) {

	locs := new(ReqMutationParams)

	err := c.Bind(locs)

	if err != nil {
		return nil, err
	}

	locations, err := dna.ParseLocations(locs.Locations)

	if err != nil {
		return nil, err
	}

	return &MutationParams{locations, locs.DBs}, nil
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
			mutations, err := db.FindMutations(location)

			if err != nil {
				return routes.ErrorReq(err)
			}

			ret = append(ret, mutations)
		}

		return routes.MakeDataResp(c, "", ret)
	})

	//return routes.MakeDataResp(c, "", mutationdbcache.GetInstance().List())
}

type PileupResp struct {
	Location *dna.Location                   `json:"location"`
	Metadata []*mutations.MutationDBMetaData `json:"metadata"`
	//Samples   uint                  `json:"samples"`
	Mutations [][]*mutations.Mutation `json:"mutations"`
}

func PileupRoute(c echo.Context) error {
	return routes.NewValidator(c).Success(func(validator *routes.Validator) error {

		params, err := ParseParamsFromPost(c)

		if err != nil {
			return routes.ErrorReq(err)
		}

		location := params.Locations[0]

		//assembly := c.Param("assembly")
		//name := c.Param("name")

		ret := PileupResp{Location: location,
			Metadata:  make([]*mutations.MutationDBMetaData, len(params.DBs)),
			Mutations: make([][]*mutations.Mutation, location.Len())}

		for i := range location.Len() {
			ret.Mutations[i] = make([]*mutations.Mutation, 0)
		}

		for dbi, id := range params.DBs {

			db, err := mutationdbcache.MutationDBFromId(id)

			if err != nil {
				return routes.ErrorReq(err)
			}

			pileup, err := db.Pileup(location)

			if err != nil {
				return routes.ErrorReq(err)
			}

			for ci := range location.Len() {
				ret.Mutations[ci] = append(ret.Mutations[ci], pileup.Mutations[ci]...)
			}

			ret.Metadata[dbi] = pileup.Metadata

			// for _, location := range locations {
			// 	pileup, err := db.Pileup(location)

			// 	if err != nil {
			// 		return routes.ErrorReq(err)
			// 	}

			// 	ret = append(ret, pileup)
			// }
		}

		// resort pileups

		for ci := range location.Len() {
			sort.Slice(ret.Mutations[ci], func(i, j int) bool {
				if ret.Mutations[ci][i].VariantType < ret.Mutations[ci][j].VariantType {
					return true
				}

				return ret.Mutations[ci][i].Tum < ret.Mutations[ci][j].Tum
			})
		}

		return routes.MakeDataResp(c, "", ret)
	})

	//return routes.MakeDataResp(c, "", mutationdbcache.GetInstance().List())
}
