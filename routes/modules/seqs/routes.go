package seqroutes

import (
	"fmt"

	"github.com/antonybholmes/go-dna"
	"github.com/antonybholmes/go-edb-server/routes"
	seq "github.com/antonybholmes/go-seqs"
	"github.com/rs/zerolog/log"

	"github.com/antonybholmes/go-seqs/seqsdbcache"
	"github.com/labstack/echo/v4"
)

type ReqSeqParams struct {
	Locations []string `json:"locations"`
	Scale     float64  `json:"scale"`
	BinSizes  []uint   `json:"binSizes"`
	Tracks    []string `json:"tracks"`
}

type SeqParams struct {
	Locations []*dna.Location
	Scale     float64
	BinSizes  []uint
	Tracks    []string
}

type SeqResp struct {
	Location  *dna.Location    `json:"location"`
	BinCounts []*seq.BinCounts `json:"binCounts"`
}

func ParseSeqParamsFromPost(c echo.Context) (*SeqParams, error) {

	var params ReqSeqParams

	err := c.Bind(&params)

	if err != nil {
		return nil, err
	}

	locations := make([]*dna.Location, 0, len(params.Locations))

	for _, loc := range params.Locations {
		location, err := dna.ParseLocation(loc)

		if err != nil {
			return nil, err
		}

		locations = append(locations, location)
	}

	return &SeqParams{Locations: locations, BinSizes: params.BinSizes, Tracks: params.Tracks, Scale: params.Scale}, nil
}

func GenomeRoute(c echo.Context) error {
	platforms, err := seqsdbcache.Genomes()

	if err != nil {
		return routes.ErrorReq(err)
	}

	return routes.MakeDataPrettyResp(c, "", platforms)
}

func PlatformRoute(c echo.Context) error {
	genome := c.Param("assembly")

	platforms, err := seqsdbcache.Platforms(genome)

	if err != nil {
		return routes.ErrorReq(err)
	}

	return routes.MakeDataPrettyResp(c, "", platforms)
}

func TracksRoute(c echo.Context) error {
	platform := c.Param("platform")
	genome := c.Param("assembly")

	tracks, err := seqsdbcache.Tracks(platform, genome)

	if err != nil {
		return routes.ErrorReq(err)
	}

	return routes.MakeDataPrettyResp(c, "", tracks)
}

func SearchSeqRoute(c echo.Context) error {
	genome := c.Param("assembly")

	if genome == "" {
		return routes.ErrorReq(fmt.Errorf("must supply a genome"))
	}

	query := c.QueryParam("search")

	tracks, err := seqsdbcache.Search(genome, query)

	if err != nil {
		return routes.ErrorReq(err)
	}

	return routes.MakeDataPrettyResp(c, "", tracks)
}

func BinsRoute(c echo.Context) error {

	params, err := ParseSeqParamsFromPost(c)

	if err != nil {
		log.Debug().Msgf("err %s", err)
		return routes.ErrorReq(err)
	}

	//log.Debug().Msgf("bin %v %v", params.Locations, params.BinSizes)

	ret := make([]*SeqResp, 0, len(params.Locations)) //make([]*seq.BinCounts, 0, len(params.Tracks))

	for li, location := range params.Locations {
		resp := SeqResp{Location: location, BinCounts: make([]*seq.BinCounts, 0, len(params.Tracks))}

		for _, track := range params.Tracks {

			reader, err := seqsdbcache.ReaderFromId(track, params.BinSizes[li], params.Scale)

			if err != nil {
				//log.Debug().Msgf("stupid err %s", err)
				return routes.ErrorReq(err)
			}

			// guarantees something is returned even with error
			// so we can ignore the errors for now to make the api
			// more robus
			binCounts, _ := reader.BinCounts(location)

			// if err != nil {
			// 	return routes.ErrorReq(err)
			// }

			resp.BinCounts = append(resp.BinCounts, binCounts)
		}

		ret = append(ret, &resp)
	}

	//log.Debug().Msgf("ret %v", len(ret))

	return routes.MakeDataPrettyResp(c, "", ret)
}
