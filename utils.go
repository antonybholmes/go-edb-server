package main

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"strconv"

	"github.com/antonybholmes/go-dna"
	"github.com/antonybholmes/go-loctogene"
	"github.com/labstack/echo/v4"
)

const DEFAULT_ASSEMBLY = "grch38"
const DEFAULT_LEVEL = 1
const DEFAULT_CHR = "chr1"   //"chr3"
const DEFAULT_START = 100000 //187728170
const DEFAULT_END = 100100   //187752257
const DEFAULT_CLOSEST_N = 10

type DNAQuery struct {
	Loc      *dna.Location
	Dir      string
	Assembly string
}

// A GeneQuery contains info from query params.
type GeneQuery struct {
	Loc      *dna.Location
	Level    int
	DB       *sql.DB
	Assembly string
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

	v = c.QueryParam("chr")

	if v != "" {
		chr = v
	} else {
		c.Logger().Warn("chr was not, using default...")
	}

	v = c.QueryParam("start")

	if v != "" {
		start, err = strconv.Atoi(v)

		if err != nil {
			c.Logger().Warn(fmt.Sprintf("%s is an invalid start, using default %d...", v, DEFAULT_START))
			start = DEFAULT_START
		}
	} else {
		c.Logger().Warn(fmt.Sprintf("start was not set, using default %d...", DEFAULT_START))
	}

	v = c.QueryParam("end")

	if v != "" {
		end, err = strconv.Atoi(v)

		if err != nil {
			c.Logger().Warn(fmt.Sprintf("%s is an invalid end, using default %d...", v, DEFAULT_END))
			end = DEFAULT_END
		}
	} else {
		c.Logger().Warn(fmt.Sprintf("end was not set, using default %d...", DEFAULT_END))
	}

	loc := dna.Location{Chr: chr, Start: start, End: end}

	return &loc, nil
}

func parseAssembly(c echo.Context) string {
	assembly := DEFAULT_ASSEMBLY

	v := c.QueryParam("assembly")

	if v != "" {
		assembly = v
	}

	return assembly
}

func parseDNAQuery(c echo.Context) (*DNAQuery, error) {
	loc, err := parseLocation(c)

	if err != nil {
		return nil, err
	}

	assembly := parseAssembly(c)

	dir := filepath.Join(os.Getenv("MODULESDIR"), "dna", assembly)

	return &DNAQuery{Loc: loc, Assembly: assembly, Dir: dir}, nil
}

func parseGeneQuery(c echo.Context) (*GeneQuery, error) {
	loc, err := parseLocation(c)

	if err != nil {
		return nil, err
	}

	assembly := parseAssembly(c)

	level := DEFAULT_LEVEL

	v := c.QueryParam("level")

	if v != "" {
		level = loctogene.GetLevel(v)
	} else {
		c.Logger().Warn(fmt.Sprintf("level was not set, using default %d...", DEFAULT_LEVEL))
	}

	file := filepath.Join(os.Getenv("MODULESDIR"), "loctogene", fmt.Sprintf("%s.db", assembly))
	db, err := loctogene.GetDB(file)

	if err != nil {
		return nil, err
	}

	return &GeneQuery{Loc: loc, Assembly: assembly, DB: db, Level: level}, nil
}
