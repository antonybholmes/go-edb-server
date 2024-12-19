package trackroutes

import (
	"github.com/antonybholmes/go-dna"
	"github.com/antonybholmes/go-edb-server/routes"
	"github.com/rs/zerolog/log"

	"github.com/antonybholmes/go-tracks"
	"github.com/antonybholmes/go-tracks/tracksdbcache"
	"github.com/labstack/echo/v4"
)

type ReqTracksParams struct {
	Location string         `json:"location"`
	BinWidth uint           `json:"binWidth"`
	Tracks   []tracks.Track `json:"tracks"`
}

type TracksParams struct {
	Location *dna.Location  `json:"location"`
	BinWidth uint           `json:"binWidth"`
	Tracks   []tracks.Track `json:"tracks"`
}

func ParseTrackParamsFromPost(c echo.Context) (*TracksParams, error) {

	var params ReqTracksParams

	err := c.Bind(&params)

	if err != nil {
		return nil, err
	}

	location, err := dna.ParseLocation(params.Location)

	if err != nil {
		return nil, err
	}

	return &TracksParams{Location: location, BinWidth: params.BinWidth, Tracks: params.Tracks}, nil
}

func PlatformRoute(c echo.Context) error {

	return routes.MakeDataPrettyResp(c, "", tracksdbcache.Platforms())
}

func GenomeRoute(c echo.Context) error {
	platform := c.Param("platform")

	genomes, err := tracksdbcache.Genomes(platform)

	if err != nil {
		return routes.ErrorReq(err)
	}

	return routes.MakeDataPrettyResp(c, "", genomes)
}

func TracksRoute(c echo.Context) error {
	platform := c.Param("platform")
	genome := c.Param("genome")

	tracks, err := tracksdbcache.Tracks(platform, genome)

	if err != nil {
		return routes.ErrorReq(err)
	}

	return routes.MakeDataPrettyResp(c, "", tracks)
}

func AllTracksRoute(c echo.Context) error {

	tracks, err := tracksdbcache.AllTracks()

	if err != nil {
		return routes.ErrorReq(err)
	}

	return routes.MakeDataPrettyResp(c, "", tracks)
}

func BinsRoute(c echo.Context) error {

	params, err := ParseTrackParamsFromPost(c)

	log.Debug().Msgf("err %s", err)

	if err != nil {
		return routes.ErrorReq(err)
	}

	ret := make([]*tracks.BinCounts, 0, len(params.Tracks))

	for _, track := range params.Tracks {
		reader := tracksdbcache.Reader(track, params.BinWidth)

		binCounts, err := reader.BinCounts(params.Location)

		if err != nil {
			return routes.ErrorReq(err)
		}

		ret = append(ret, binCounts)
	}

	return routes.MakeDataPrettyResp(c, "", ret)
}
