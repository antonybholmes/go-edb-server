package routes

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/antonybholmes/go-dna"
	"github.com/antonybholmes/go-edb-api/utils"
	"github.com/antonybholmes/go-gene"
	"github.com/antonybholmes/go-loctogene"
	"github.com/antonybholmes/go-math"
	"github.com/labstack/echo/v4"
)

const MAX_ANNOTATIONS = 1000

type AnnotationResponse struct {
	Status int                    `json:"status"`
	Data   []*gene.GeneAnnotation `json:"data"`
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

func AnnotationRoute(c echo.Context, loctogenedbcache *loctogene.LoctogeneDbCache) error {
	locations, err := ParseLocationsFromPost(c)

	if err != nil {
		return utils.MakeBadResp(c, err)
	}

	// limit amount of data returned per request to 1000 entries at a time
	locations = locations[0:math.IntMin(len(locations), MAX_ANNOTATIONS)]

	query, err := ParseGeneQuery(c, c.Param("assembly"), loctogenedbcache)

	if err != nil {
		return utils.MakeBadResp(c, err)
	}

	n := utils.ParseN(c, DEFAULT_CLOSEST_N)

	tssRegion := ParseTSSRegion(c)

	output := utils.ParseOutput(c)

	annotationDb := gene.NewAnnotateDb(query.Db, tssRegion, n)

	data := []*gene.GeneAnnotation{}

	for _, location := range locations {

		annotations, err := annotationDb.Annotate(&location)

		if err != nil {
			return utils.MakeBadResp(c, err)
		}

		data = append(data, annotations)
	}

	if output == "text" {
		tsv, err := MakeGeneTable(data, tssRegion)

		if err != nil {
			return utils.MakeBadResp(c, err)
		}

		return c.String(http.StatusOK, tsv)
	} else {

		return c.JSON(http.StatusOK, AnnotationResponse{Status: http.StatusOK, Data: data})
	}
}

func MakeGeneTable(
	data []*gene.GeneAnnotation,
	ts *dna.TSSRegion,
) (string, error) {
	buffer := new(bytes.Buffer)
	wtr := csv.NewWriter(buffer)
	wtr.Comma = '\t'

	closestN := len(data[0].ClosestGenes)

	headers := []string{"Location", "ID", "Gene Symbol", fmt.Sprintf(
		"Relative To Gene (prom=-%d/+%dkb)",
		ts.Offset5P()/1000,
		ts.Offset3P()/1000), "TSS Distance", "Gene Location"}

	for i := 1; i <= closestN; i++ {
		headers = append(headers, fmt.Sprintf("#%d Closest ID", i))
		headers = append(headers, fmt.Sprintf("#%d Closest Gene Symbols", i))
		headers = append(headers, fmt.Sprintf(
			"#%d Relative To Closet Gene (prom=-%d/+%dkb)",
			i,
			ts.Offset5P()/1000,
			ts.Offset3P()/1000))
		headers = append(headers, fmt.Sprintf("#%d TSS Closest Distance", i))
		headers = append(headers, fmt.Sprintf("#%d Gene Location", i))
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
			annotation.Dists,
			annotation.Locations}

		for _, closestGene := range annotation.ClosestGenes {
			row = append(row, closestGene.Feature.GeneId)
			row = append(row, gene.LabelGene(closestGene.Feature.GeneSymbol, closestGene.Feature.Strand))
			row = append(row, closestGene.PromLabel)
			row = append(row, strconv.Itoa(closestGene.TssDist))
			row = append(row, closestGene.Feature.ToLocation().String())
		}

		err := wtr.Write(row)

		if err != nil {
			return "", err
		}
	}

	wtr.Flush()

	return buffer.String(), nil
}
