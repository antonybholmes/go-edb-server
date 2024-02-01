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

func parseLocation(c echo.Context) *dna.Location {
	chr := DEFAULT_CHR
	start := DEFAULT_START
	end := DEFAULT_END

	var v string
	var err error

	v = c.QueryParam("chr")

	if v != "" {
		chr = v
	}

	v = c.QueryParam("start")

	if v != "" {
		start, err = strconv.Atoi(v)

		if err != nil {
			start = DEFAULT_START
		}
	}

	v = c.QueryParam("end")

	if v != "" {
		end, err = strconv.Atoi(v)

		if err != nil {
			end = DEFAULT_END
		}
	}

	loc := dna.Location{Chr: chr, Start: start, End: end}

	return &loc
}

func parseAssembly(c echo.Context) string {
	assembly := DEFAULT_ASSEMBLY

	v := c.QueryParam("assembly")

	if v != "" {
		assembly = v
	}

	return assembly
}

func parseDNAQuery(c echo.Context, modulesDir string) (*DNAQuery, error) {
	loc := parseLocation(c)

	assembly := parseAssembly(c)

	dir := filepath.Join(modulesDir, "dna", assembly)

	_, err := os.Stat(dir)

	if err != nil {
		return nil, fmt.Errorf("%s is not a valid assembly", assembly)
	}

	return &DNAQuery{Loc: loc, Assembly: assembly, Dir: dir}, nil
}

func parseGeneQuery(c echo.Context, modulesDir string) (*GeneQuery, error) {
	loc := parseLocation(c)

	assembly := parseAssembly(c)

	level := DEFAULT_LEVEL

	v := c.QueryParam("level")

	if v != "" {
		level = loctogene.GetLevel(v)
	}

	file := filepath.Join(modulesDir, "loctogene", fmt.Sprintf("%s.db", assembly))
	db, err := loctogene.GetDB(file)

	if err != nil {
		return nil, fmt.Errorf("unable to open database for assembly %s", assembly)
	}

	return &GeneQuery{Loc: loc, Assembly: assembly, DB: db, Level: level}, nil
}
