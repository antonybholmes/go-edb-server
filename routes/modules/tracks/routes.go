package trackroutes

import (
	"fmt"

	"github.com/antonybholmes/go-dna"
	"github.com/antonybholmes/go-edb-server/routes"
	"github.com/rs/zerolog/log"

	"github.com/antonybholmes/go-tracks"
	"github.com/antonybholmes/go-tracks/tracksdbcache"
	"github.com/labstack/echo/v4"
)

type ReqTracksParams struct {
	Location string   `json:"location"`
	Scale    float64  `json:"scale"`
	BinWidth uint     `json:"binWidth"`
	Tracks   []string `json:"tracks"`
}

type TracksParams struct {
	Location *dna.Location `json:"location"`
	Scale    float64       `json:"scale"`
	BinWidth uint          `json:"binWidth"`
	Tracks   []string      `json:"tracks"`
}

func ParseTrackParamsFromPost(c echo.Context) (*TracksParams, error) {

	var params ReqTracksParams

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

	log.Debug().Msgf("scale %f", params.Scale)

	return &TracksParams{Location: location, BinWidth: params.BinWidth, Tracks: params.Tracks, Scale: params.Scale}, nil
}

func GenomeRoute(c echo.Context) error {
	platforms, err := tracksdbcache.Genomes()

	if err != nil {
		return routes.ErrorReq(err)
	}

	return routes.MakeDataPrettyResp(c, "", platforms)
}

func PlatformRoute(c echo.Context) error {
	genome := c.Param("assembly")

	platforms, err := tracksdbcache.Platforms(genome)

	if err != nil {
		return routes.ErrorReq(err)
	}

	return routes.MakeDataPrettyResp(c, "", platforms)
}

func TracksRoute(c echo.Context) error {
	platform := c.Param("platform")
	genome := c.Param("assembly")

	tracks, err := tracksdbcache.Tracks(platform, genome)

	if err != nil {
		return routes.ErrorReq(err)
	}

	return routes.MakeDataPrettyResp(c, "", tracks)
}

func SearchTracksRoute(c echo.Context) error {
	genome := c.Param("assembly")

	if genome == "" {
		return routes.ErrorReq(fmt.Errorf("must supply a genome"))
	}

	query := c.QueryParam("search")

	tracks, err := tracksdbcache.Search(genome, query)

	if err != nil {
		return routes.ErrorReq(err)
	}

	return routes.MakeDataPrettyResp(c, "", tracks)
}

func BinsRoute(c echo.Context) error {

	params, err := ParseTrackParamsFromPost(c)

	if err != nil {
		log.Debug().Msgf("bins param err %s", err)
		return routes.ErrorReq(err)
	}

	ret := make([]*tracks.BinCounts, 0, len(params.Tracks))

	for _, track := range params.Tracks {
		log.Debug().Msgf("track %v %f", track, params.Scale)

		reader, err := tracksdbcache.ReaderFromId(track, params.BinWidth, params.Scale)

		if err != nil {
			return routes.ErrorReq(err)
		}

		binCounts, err := reader.BinCounts(params.Location)

		if err != nil {
			return routes.ErrorReq(err)
		}

		ret = append(ret, binCounts)
	}

	return routes.MakeDataPrettyResp(c, "", ret)
}
