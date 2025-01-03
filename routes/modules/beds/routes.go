package bedroutes

import (
	"fmt"

	"github.com/antonybholmes/go-beds"
	"github.com/antonybholmes/go-dna"
	"github.com/antonybholmes/go-edb-server/routes"

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
		return nil, err
	}

	location, err := dna.ParseLocation(params.Location)

	if err != nil {
		return nil, err
	}

	return &BedsParams{Location: location, Beds: params.Beds}, nil
}

func GenomeRoute(c echo.Context) error {
	platforms, err := bedsdbcache.Genomes()

	if err != nil {
		return routes.ErrorReq(err)
	}

	return routes.MakeDataPrettyResp(c, "", platforms)
}

func PlatformRoute(c echo.Context) error {
	genome := c.Param("assembly")

	platforms, err := bedsdbcache.Platforms(genome)

	if err != nil {
		return routes.ErrorReq(err)
	}

	return routes.MakeDataPrettyResp(c, "", platforms)
}

func SearchBedsRoute(c echo.Context) error {
	genome := c.Param("assembly")

	if genome == "" {
		return routes.ErrorReq(fmt.Errorf("must supply a genome"))
	}

	query := c.QueryParam("search")

	tracks, err := bedsdbcache.Search(genome, query)

	if err != nil {
		return routes.ErrorReq(err)
	}

	return routes.MakeDataPrettyResp(c, "", tracks)
}

func BedRegionsRoute(c echo.Context) error {

	params, err := ParseBedParamsFromPost(c)

	if err != nil {
		return routes.ErrorReq(err)
	}

	if len(params.Beds) == 0 {
		return routes.ErrorReq(fmt.Errorf("at least 1 bed id must be supplied"))
	}

	ret := make([][]*beds.BedRegion, 0, len(params.Beds))

	for _, bed := range params.Beds {

		//log.Debug().Msgf("bed id %s", bed)

		reader, err := bedsdbcache.ReaderFromId(bed)

		if err != nil {
			return routes.ErrorReq(err)
		}

		features, _ := reader.OverlappingRegions(params.Location)

		ret = append(ret, features)
	}

	return routes.MakeDataPrettyResp(c, "", ret)
}
