package bedroutes

import (
	"fmt"

	"github.com/antonybholmes/go-dna"
	"github.com/antonybholmes/go-edb-server/routes"
	"github.com/rs/zerolog/log"

	"github.com/antonybholmes/go-beds/bedsdbcache"
	"github.com/labstack/echo/v4"
)

type ReqBedsParams struct {
	Location string   `json:"location"`
	Beds     []string `json:"beds"`
}

type BedsParams struct {
	Location *dna.Location `json:"location"`
	Beds     []string      `json:"beds"`
}

func ParseBedParamsFromPost(c echo.Context) (*BedsParams, error) {

	var params ReqBedsParams

	err := c.Bind(&params)

	if err != nil {
		log.Debug().Msgf("bind err %s", err)
		return nil, err
	}

	location, err := dna.ParseLocation(params.Location)

	if err != nil {
		log.Debug().Msgf("loc err %s", err)
		return nil, err
	}

	return &BedsParams{Location: location, Beds: params.Beds}, nil
}

func PlatformRoute(c echo.Context) error {
	platforms, err := bedsdbcache.Platforms()

	if err != nil {
		return routes.ErrorReq(err)
	}

	return routes.MakeDataPrettyResp(c, "", platforms)
}

func GenomeRoute(c echo.Context) error {
	platform := c.Param("platform")

	genomes, err := bedsdbcache.Genomes(platform)

	if err != nil {
		return routes.ErrorReq(err)
	}

	return routes.MakeDataPrettyResp(c, "", genomes)
}

func AllBedsRoute(c echo.Context) error {
	genome := c.Param("assembly")

	if genome == "" {
		return routes.ErrorReq(fmt.Errorf("must supply a genome"))
	}

	tracks, err := bedsdbcache.AllBeds(genome)

	if err != nil {
		return routes.ErrorReq(err)
	}

	return routes.MakeDataPrettyResp(c, "", tracks)
}

func BedFeaturesRoute(c echo.Context) error {

	params, err := ParseBedParamsFromPost(c)

	if err != nil {
		log.Debug().Msgf("bins param err %s", err)
		return routes.ErrorReq(err)
	}

	if len(params.Beds) == 0 {
		return routes.ErrorReq(fmt.Errorf("at least 1 bed id must be supplied"))
	}

	bed := params.Beds[0]

	log.Debug().Msgf("bed %s", bed)

	reader, err := bedsdbcache.ReaderFromId(bed)

	if err != nil {
		return routes.ErrorReq(err)
	}

	features, err := reader.BedFeatures(params.Location)

	if err != nil {
		return routes.ErrorReq(err)
	}

	return routes.MakeDataPrettyResp(c, "", features)
}
