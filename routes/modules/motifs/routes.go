package motifroutes

import (
	"github.com/antonybholmes/go-edb-server/routes"
	"github.com/antonybholmes/go-motifs"
	"github.com/antonybholmes/go-motifs/motifsdb"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
)

const MIN_SEARCH_LEN = 3

type ReqParams struct {
	Search     string `json:"search"`
	Exact      bool   `json:"exact"`
	Reverse    bool   `json:"reverse"`
	Complement bool   `json:"complement"`
}

type MotifRes struct {
	Search     string          `json:"search"`
	Motifs     []*motifs.Motif `json:"motifs"`
	Reverse    bool            `json:"reverse"`
	Complement bool            `json:"complement"`
}

func ParseParamsFromPost(c echo.Context) (*ReqParams, error) {

	var params ReqParams

	err := c.Bind(&params)

	if err != nil {
		return nil, err
	}

	return &params, nil
}

func SearchRoute(c echo.Context) error {

	params, err := ParseParamsFromPost(c)

	if err != nil {
		return routes.ErrorReq(err)
	}

	search := params.Search

	if len(search) < MIN_SEARCH_LEN {
		return routes.ErrorReq("Search too short")
	}

	log.Debug().Msgf("motif %v", params)

	// Don't care about the errors, just plug empty list into failures
	motifs, err := motifsdb.Search(search, params.Reverse, params.Complement)

	if err != nil {
		log.Debug().Msgf("motif %s", err)
		return routes.ErrorReq(err)
	}

	return routes.MakeDataPrettyResp(c, "",
		MotifRes{
			Search:     search,
			Motifs:     motifs,
			Reverse:    params.Reverse,
			Complement: params.Complement,
		})

	//return routes.MakeDataResp(c, "", mutationdbcache.GetInstance().List())
}
