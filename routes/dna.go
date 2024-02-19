package routes

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/antonybholmes/go-dna"
	"github.com/antonybholmes/go-dna/dnadbcache"

	"github.com/labstack/echo/v4"
)

const DEFAULT_ASSEMBLY = "grch38"
const DEFAULT_CHR = "chr1"        //"chr3"
const DEFAULT_START uint = 100000 //187728170
const DEFAULT_END uint = 100100   //187752257

type ReqLocs struct {
	Locations []dna.Location `json:"locations"`
}

type DNA struct {
	Location *dna.Location `json:"location"`
	DNA      string        `json:"dna"`
}

type DNAResp struct {
	Assembly     string `json:"assembly"`
	Format       string `json:"format"`
	IsRev        bool   `json:"isRev"`
	IsComplement bool   `json:"isComp"`
	Seqs         []*DNA `json:"seqs"`
}

type DNAQuery struct {
	Rev        bool
	Comp       bool
	Format     string
	RepeatMask string
}

func ParseLocation(c echo.Context) (*dna.Location, error) {
	chr := DEFAULT_CHR
	start := DEFAULT_START
	end := DEFAULT_END

	var v string
	var err error
	var n uint64

	v = c.QueryParam("chr")

	if v != "" {
		chr = v
	}

	v = c.QueryParam("start")

	if v != "" {
		n, err = strconv.ParseUint(v, 10, 0)

		if err != nil {
			return nil, fmt.Errorf("%s is an invalid start", v)
		}

		start = uint(n)
	}

	v = c.QueryParam("end")

	if v != "" {
		n, err = strconv.ParseUint(v, 10, 0)

		if err != nil {
			return nil, fmt.Errorf("%s is an invalid end", v)
		}

		end = uint(n)
	}

	loc := dna.NewLocation(chr, start, end)

	return loc, nil
}

func ParseLocationsFromPost(c echo.Context) ([]dna.Location, error) {
	var err error
	locs := new(ReqLocs)

	err = c.Bind(locs)

	if err != nil {
		return nil, err
	}

	return locs.Locations, nil
}

func ParseDNAQuery(c echo.Context) (*DNAQuery, error) {
	var err error

	rev := false
	v := c.QueryParam("rev")

	if v != "" {
		rev, err = strconv.ParseBool(v)

		if err != nil {
			rev = false
		}
	}

	comp := false
	v = c.QueryParam("comp")

	if v != "" {
		comp, err = strconv.ParseBool(v)

		if err != nil {
			comp = false
		}
	}

	format := ""
	v = c.QueryParam("format")

	if v != "" {
		if strings.Contains(strings.ToLower(v), "lower") {
			format = "lower"
		} else {
			format = "upper"
		}
	}

	repeatMask := ""
	v = c.QueryParam("mask")

	if v != "" {
		if strings.Contains(strings.ToLower(v), "lower") {
			repeatMask = "lower"
		} else {
			repeatMask = "n"
		}
	}

	return &DNAQuery{Rev: rev, Comp: comp, Format: format, RepeatMask: repeatMask}, nil
}

func DNARoute(c echo.Context) error {

	locations, err := ParseLocationsFromPost(c)

	if err != nil {
		return BadReq(err)
	}

	assembly := c.Param("assembly")

	query, err := ParseDNAQuery(c)

	if err != nil {
		return BadReq(err)
	}

	//c.Logger().Debugf("%s %s", query.Loc, query.Dir)

	dnadb, err := dnadbcache.Db(assembly, query.Format, query.RepeatMask)

	if err != nil {
		return BadReq(err)
	}

	seqs := []*DNA{}

	for _, location := range locations {
		dna, err := dnadb.DNA(&location, query.Rev, query.Comp)

		if err != nil {
			return BadReq(err)
		}

		seqs = append(seqs, &DNA{Location: &location, DNA: dna})
	}

	//c.Logger().Debugf("%s", dna)

	return MakeDataResp(c, "", &DNAResp{Assembly: assembly, Format: query.Format, IsRev: query.Rev, IsComplement: query.Comp, Seqs: seqs})
}
