package generoutes

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/antonybholmes/go-dna"
	"github.com/antonybholmes/go-edb-server/routes"
	dnaroutes "github.com/antonybholmes/go-edb-server/routes/modules/dna"
	"github.com/antonybholmes/go-genes"
	"github.com/antonybholmes/go-genes/genedbcache"
	basemath "github.com/antonybholmes/go-math"
	"github.com/labstack/echo/v4"
)

const DEFAULT_LEVEL = genes.LEVEL_GENE

const DEFAULT_CLOSEST_N uint16 = 5

// A GeneQuery contains info from query params.
type GeneQuery struct {
	Level    genes.Level
	Db       *genes.GeneDB
	Assembly string
	// only show canonical genes
	Canonical bool
}

type GenesResp struct {
	Location *dna.Location           `json:"location"`
	Features []*genes.GenomicFeature `json:"features"`
}

const MAX_ANNOTATIONS = 1000

type AnnotationResponse struct {
	Status int                     `json:"status"`
	Data   []*genes.GeneAnnotation `json:"data"`
}

func ParseGeneQuery(c echo.Context, assembly string) (*GeneQuery, error) {
	level := genes.LEVEL_GENE

	v := c.QueryParam("level")

	if v != "" {
		level = genes.ParseLevel(v)
	}

	canonical := strings.HasPrefix(strings.ToLower(c.QueryParam("canonical")), "t")

	db, err := genedbcache.GeneDB(assembly)

	if err != nil {
		return nil, fmt.Errorf("unable to open database for assembly %s %s", assembly, err)
	}

	return &GeneQuery{Assembly: assembly, Db: db, Level: level, Canonical: canonical}, nil
}

// func GeneDBInfoRoute(c echo.Context) error {
// 	query, err := ParseGeneQuery(c, c.Param("assembly"))

// 	if err != nil {
// 		return routes.ErrorReq(err)
// 	}

// 	info, _ := query.Db.GeneDBInfo()

// 	// if err != nil {
// 	// 	return routes.ErrorReq(err)
// 	// }

// 	return routes.MakeDataPrettyResp(c, "", &info)
// }

func GenomesRoute(c echo.Context) error {
	infos, err := genedbcache.GetInstance().List()

	if err != nil {
		return routes.ErrorReq(err)
	}

	return routes.MakeDataPrettyResp(c, "", infos)
}

func OverlappingGenesRoute(c echo.Context) error {
	locations, err := dnaroutes.ParseLocationsFromPost(c) // dnaroutes.ParseLocationsFromPost(c)

	if err != nil {
		return routes.ErrorReq(err)
	}

	query, err := ParseGeneQuery(c, c.Param("assembly"))

	if err != nil {
		return routes.ErrorReq(err)
	}

	if len(locations) == 0 {
		return routes.ErrorReq(fmt.Errorf("must supply at least 1 location"))
	}

	ret := make([]*GenesResp, 0, len(locations))

	for _, location := range locations {
		features, err := query.Db.OverlappingGenes(location, query.Canonical)

		if err != nil {
			return routes.ErrorReq(err)
		}

		ret = append(ret, &GenesResp{Location: location, Features: features})

	}

	return routes.MakeDataPrettyResp(c, "", &ret)
}

func GeneInfoRoute(c echo.Context) error {
	search := c.QueryParam("search") // dnaroutes.ParseLocationsFromPost(c)

	if search == "" {
		return routes.ErrorReq(fmt.Errorf("search cannot be empty"))
	}

	query, err := ParseGeneQuery(c, c.Param("assembly"))

	if err != nil {
		return routes.ErrorReq(err)
	}

	features, _ := query.Db.GeneInfo(search, query.Level)

	// if err != nil {
	// 	return routes.ErrorReq(err)
	// }

	return routes.MakeDataPrettyResp(c, "", &features)
}

func WithinGenesRoute(c echo.Context) error {
	locations, err := dnaroutes.ParseLocationsFromPost(c) // dnaroutes.ParseLocationsFromPost(c)

	if err != nil {
		return routes.ErrorReq(err)
	}

	query, err := ParseGeneQuery(c, c.Param("assembly"))

	if err != nil {
		return routes.ErrorReq(err)
	}

	data := make([]*genes.GenomicFeatures, len(locations))

	for li, location := range locations {
		genes, err := query.Db.WithinGenes(location, query.Level)

		if err != nil {
			return routes.ErrorReq(err)
		}

		data[li] = genes
	}

	return routes.MakeDataPrettyResp(c, "", &data)
}

// Find the n closest genes to a location
func ClosestGeneRoute(c echo.Context) error {
	locations, err := dnaroutes.ParseLocationsFromPost(c)

	if err != nil {
		return routes.ErrorReq(err)
	}

	query, err := ParseGeneQuery(c, c.Param("assembly"))

	if err != nil {
		return routes.ErrorReq(err)
	}

	n := routes.ParseN(c, DEFAULT_CLOSEST_N)

	data := make([]*genes.GenomicFeatures, len(locations))

	for li, location := range locations {
		genes, err := query.Db.ClosestGenes(location, n, query.Level)

		if err != nil {
			return routes.ErrorReq(err)
		}

		data[li] = genes
	}

	return routes.MakeDataPrettyResp(c, "", &data)
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

func AnnotateRoute(c echo.Context) error {
	locations, err := dnaroutes.ParseLocationsFromPost(c)

	if err != nil {
		return routes.ErrorReq(err)
	}

	// limit amount of data returned per request to 1000 entries at a time
	locations = locations[0:basemath.Min(len(locations), MAX_ANNOTATIONS)]

	query, err := ParseGeneQuery(c, c.Param("assembly"))

	if err != nil {
		return routes.ErrorReq(err)
	}

	n := routes.ParseN(c, DEFAULT_CLOSEST_N)

	tssRegion := ParseTSSRegion(c)

	output := routes.ParseOutput(c)

	annotationDb := genes.NewAnnotateDb(query.Db, tssRegion, n)

	data := make([]*genes.GeneAnnotation, len(locations))

	for li, location := range locations {

		annotations, err := annotationDb.Annotate(location)

		if err != nil {
			return routes.ErrorReq(err)
		}

		data[li] = annotations
	}

	if output == "text" {
		tsv, err := MakeGeneTable(data, tssRegion)

		if err != nil {
			return routes.ErrorReq(err)
		}

		return c.String(http.StatusOK, tsv)
	} else {

		return c.JSON(http.StatusOK, AnnotationResponse{Status: http.StatusOK, Data: data})
	}
}

func MakeGeneTable(
	data []*genes.GeneAnnotation,
	ts *dna.TSSRegion,
) (string, error) {
	var buffer bytes.Buffer
	wtr := csv.NewWriter(&buffer)
	wtr.Comma = '\t'

	closestN := len(data[0].ClosestGenes)

	headers := make([]string, 6+5*closestN)

	headers[0] = "Location"
	headers[1] = "ID"
	headers[2] = "Gene Symbol"
	headers[3] = fmt.Sprintf(
		"Relative To Gene (prom=-%d/+%dkb)",
		ts.Offset5P()/1000,
		ts.Offset3P()/1000)
	headers[4] = "TSS Distance"
	headers[5] = "Gene Location"

	idx := 6
	for i := 1; i <= closestN; i++ {
		headers[idx] = fmt.Sprintf("#%d Closest ID", i)
		headers[idx] = fmt.Sprintf("#%d Closest Gene Symbols", i)
		headers[idx] = fmt.Sprintf(
			"#%d Relative To Closet Gene (prom=-%d/+%dkb)",
			i,
			ts.Offset5P()/1000,
			ts.Offset3P()/1000)
		headers[idx] = fmt.Sprintf("#%d TSS Closest Distance", i)
		headers[idx] = fmt.Sprintf("#%d Gene Location", i)
	}

	err := wtr.Write(headers)

	if err != nil {
		return "", err
	}

	for _, annotation := range data {
		row := []string{annotation.Location.String(),
			annotation.GeneIds,
			annotation.GeneSymbols,
			annotation.PromLabels,
			annotation.TSSDists,
			annotation.Locations}

		for _, closestGene := range annotation.ClosestGenes {
			row = append(row, closestGene.GeneId)
			row = append(row, genes.GeneWithStrandLabel(closestGene.GeneSymbol, closestGene.Strand))
			row = append(row, closestGene.PromLabel)
			row = append(row, strconv.Itoa(closestGene.TssDist))
			row = append(row, closestGene.Location.String())
		}

		err := wtr.Write(row)

		if err != nil {
			return "", err
		}
	}

	wtr.Flush()

	return buffer.String(), nil
}
