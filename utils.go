package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"

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

type DNAQuery struct {
	Loc      *dna.Location
	Dir      string
	Assembly string
	Rev      bool
	Comp     bool
}

// A GeneQuery contains info from query params.
type GeneQuery struct {
	Loc      *dna.Location
	Level    loctogene.Level
	DB       *loctogene.LoctogeneDB
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

func parseLocation(c echo.Context) (*dna.Location, error) {
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

func parseAssembly(c echo.Context) string {
	assembly := DEFAULT_ASSEMBLY

	v := c.QueryParam("assembly")

	if v != "" {
		assembly = v
	}

	return assembly
}

func parseN(c echo.Context) uint16 {

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

func parseDNAQuery(c echo.Context) (*DNAQuery, error) {
	loc, err := parseLocation(c)

	if err != nil {
		return nil, err
	}

	assembly := parseAssembly(c)

	dir := filepath.Join("data/dna", assembly)

	_, err = os.Stat(dir)

	if err != nil {
		return nil, fmt.Errorf("%s is not a valid assembly", assembly)
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

	return &DNAQuery{Loc: loc, Assembly: assembly, Dir: dir, Rev: rev, Comp: comp}, nil
}

func parseGeneQuery(c echo.Context, assembly string) (*GeneQuery, error) {
	loc, err := parseLocation(c)

	if err != nil {
		return nil, err
	}

	level := loctogene.Gene

	v := c.QueryParam("level")

	if v != "" {
		level = loctogene.ParseLevel(v)
	}

	file := filepath.Join("data/loctogene", fmt.Sprintf("%s.db", assembly))
	c.Logger().Debugf("%s", file)
	db, err := loctogene.NewLoctogeneDB(file)

	if err != nil {
		return nil, fmt.Errorf("unable to open database for assembly %s %s", assembly, err)
	}

	return &GeneQuery{Loc: loc, Assembly: assembly, DB: db, Level: level}, nil
}
