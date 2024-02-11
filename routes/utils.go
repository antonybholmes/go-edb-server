package routes

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/antonybholmes/go-dna"
	"github.com/antonybholmes/go-loctogene"
	"github.com/labstack/echo/v4"
)

const DEFAULT_ASSEMBLY = "grch38"
const DEFAULT_LEVEL = loctogene.Gene
const DEFAULT_CHR = "chr1"        //"chr3"
const DEFAULT_START uint = 100000 //187728170
const DEFAULT_END uint = 100100   //187752257
const DEFAULT_CLOSEST_N uint16 = 5

type StatusMessage struct {
	Message string `json:"message"`
	Status  int    `json:"status"`
}

// A GeneQuery contains info from query params.
type GeneQuery struct {
	Loc      *dna.Location
	Level    loctogene.Level
	Db       *loctogene.LoctogeneDb
	Assembly string
}

// Max returns the larger of x or y.
func Max(x, y int) int {
	if x < y {
		return y
	}
	return x
}

// Min returns the smaller of x or y.
func Min(x, y int) int {
	if x > y {
		return y
	}
	return x
}

// parsedLocation takes an echo context and attempts to extract parameters
// from the query string and return the location to check, the assembly
// (e.g. grch38) to search, the level of detail (1=gene,2=transcript,3=exon).
// If parameters are not provided defaults are used, but if parameters are
// considered invalid, it will throw an error.

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

// func parseAssembly(c echo.Context) string {
// 	assembly := DEFAULT_ASSEMBLY

// 	v := c.QueryParam("assembly")

// 	if v != "" {
// 		assembly = v
// 	}

// 	return assembly
// }

func ParseN(c echo.Context) uint16 {

	v := c.QueryParam("n")

	if v == "" {
		return DEFAULT_CLOSEST_N
	}

	n, err := strconv.ParseUint(v, 10, 0)

	if err != nil {
		return DEFAULT_CLOSEST_N
	}

	return uint16(n)
}

func ParseTSSRegion(c echo.Context) *dna.TSSRegion {

	v := c.QueryParam("tss")

	if v == "" {
		return dna.NewTSSRegion(2000, 1000)
	}

	tokens := strings.Split(v, ",")

	s, err := strconv.ParseUint(tokens[0], 10, 0)

	if err != nil {
		return dna.NewTSSRegion(2000, 1000)
	}

	e, err := strconv.ParseUint(tokens[1], 10, 0)

	if err != nil {
		return dna.NewTSSRegion(2000, 1000)
	}

	return dna.NewTSSRegion(uint(s), uint(e))
}

func ParseOutput(c echo.Context) string {

	v := c.QueryParam("output")

	if strings.Contains(strings.ToLower(v), "text") {
		return "text"
	} else {
		return "json"
	}
}

func ParseDNAQuery(c echo.Context) (*DNAQuery, error) {
	loc, err := ParseLocation(c)

	if err != nil {
		return nil, err
	}

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

	return &DNAQuery{Loc: loc, Rev: rev, Comp: comp, Format: format, RepeatMask: repeatMask}, nil
}

func ParseGeneQuery(c echo.Context, assembly string, loctogenedbcache *loctogene.LoctogeneDbCache) (*GeneQuery, error) {
	loc, err := ParseLocation(c)

	if err != nil {
		return nil, err
	}

	level := loctogene.Gene

	v := c.QueryParam("level")

	if v != "" {
		level = loctogene.ParseLevel(v)
	}

	db, err := loctogenedbcache.Db(assembly)

	if err != nil {
		return nil, fmt.Errorf("unable to open database for assembly %s %s", assembly, err)
	}

	return &GeneQuery{Loc: loc, Assembly: assembly, Db: db, Level: level}, nil
}
