package mutationroutes

import (
	"sort"

	"github.com/antonybholmes/go-dna"
	"github.com/antonybholmes/go-edb-api/routes"
	"github.com/antonybholmes/go-edb-api/routes/modroutes/dnaroutes"
	"github.com/antonybholmes/go-mutations"
	"github.com/antonybholmes/go-mutations/mutationdbcache"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
)

type MutationParams struct {
	Locations []*dna.Location
	Databases []string
}

type ReqMutationParams struct {
	Locations []string `json:"locations"`
	Databases []string `json:"databases"`
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

	return &MutationParams{locations, locs.Databases}, nil
}

func MutationDatabasesRoute(c echo.Context) error {

	return routes.MakeDataResp(c, "", mutationdbcache.GetInstance().List())
}

func MutationsRoute(c echo.Context) error {
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

		ret := make([]*mutations.MutationResults, len(locations))

		for i, location := range locations {
			mutations, err := db.FindMutations(location)

			if err != nil {
				return routes.ErrorReq(err)
			}

			ret[i] = mutations
		}

		return routes.MakeDataResp(c, "", ret)
	})

	//return routes.MakeDataResp(c, "", mutationdbcache.GetInstance().List())
}

type MafResp struct {
	Location  *dna.Location               `json:"location"`
	Info      []*mutations.MutationDBInfo `json:"info"`
	Samples   int                         `json:"samples"`
	Mutations [][]string                  `json:"mutations"`
}

func MafRoute(c echo.Context) error {
	return routes.NewValidator(c).Success(func(validator *routes.Validator) error {

		params, err := ParseParamsFromPost(c)

		if err != nil {
			return routes.ErrorReq(err)
		}

		location := params.Locations[0]

		//assembly := c.Param("assembly")
		//name := c.Param("name")

		ret := MafResp{Location: location,
			Info:      make([]*mutations.MutationDBInfo, len(params.Databases)),
			Mutations: make([][]string, location.Len())}

		for i := range location.Len() {
			ret.Mutations[i] = make([]string, 0, 10)
		}

		sampleMap := make([]map[string]struct{}, location.Len())

		for i := range location.Len() {
			sampleMap[i] = make(map[string]struct{})
		}

		for dbi, id := range params.Databases {
			db, err := mutationdbcache.MutationDBFromId(id)

			if err != nil {
				return routes.ErrorReq(err)
			}

			results, err := db.FindMutations(location)

			if err != nil {
				return routes.ErrorReq(err)
			}

			// sum the total number of samples involved
			ret.Samples += len(db.Info.Samples)

			for _, mutation := range results.Mutations {
				offset := mutation.Start - location.Start
				sample := mutation.Sample

				_, ok := sampleMap[offset][sample]

				if !ok {
					sampleMap[offset][sample] = struct{}{}
				}

			}

			ret.Info[dbi] = db.Info
		}

		// sort each pileup
		for ci := range location.Len() {
			if len(sampleMap[ci]) > 0 {
				samples := make([]string, 0, len(sampleMap[ci]))

				for sample := range sampleMap[ci] {
					samples = append(samples, sample)
				}

				// sort the samples for ease of use
				sort.Strings(samples)

				ret.Mutations[ci] = append(ret.Mutations[ci], samples...)
			}

		}

		return routes.MakeDataResp(c, "", ret)
	})

	//return routes.MakeDataResp(c, "", mutationdbcache.GetInstance().List())
}

type PileupResp struct {
	Location  *dna.Location               `json:"location"`
	Info      []*mutations.MutationDBInfo `json:"info"`
	Samples   int                         `json:"samples"`
	Mutations [][]*mutations.Mutation     `json:"mutations"`
}

func PileupRoute(c echo.Context) error {
	return routes.NewValidator(c).Success(func(validator *routes.Validator) error {

		params, err := ParseParamsFromPost(c)

		if err != nil {
			return routes.ErrorReq(err)
		}

		log.Debug().Msgf("pileup: %v", params)

		location := params.Locations[0]

		//assembly := c.Param("assembly")
		//name := c.Param("name")

		ret := PileupResp{Location: location,
			// one metadata file for each database requested
			Info:      make([]*mutations.MutationDBInfo, len(params.Databases)),
			Mutations: make([][]*mutations.Mutation, location.Len())}

		for i := range location.Len() {
			ret.Mutations[i] = make([]*mutations.Mutation, 0, 10)
		}

		for dbi, id := range params.Databases {
			db, err := mutationdbcache.MutationDBFromId(id)

			if err != nil {
				return routes.ErrorReq(err)
			}

			pileup, err := db.Pileup(location)

			if err != nil {
				return routes.ErrorReq(err)
			}

			// sum the total number of samples involved
			ret.Samples += len(db.Info.Samples)

			for ci := range location.Len() {
				ret.Mutations[ci] = append(ret.Mutations[ci], pileup.Mutations[ci]...)
			}

			ret.Info[dbi] = db.Info
		}

		// sort each pileup
		for ci := range location.Len() {
			sort.Slice(ret.Mutations[ci], func(i, j int) bool {
				// sort by variant type and then tumor
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
