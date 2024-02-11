package routes

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"net/http"
	"strconv"

	"github.com/antonybholmes/go-dna"
	"github.com/antonybholmes/go-gene"
	"github.com/antonybholmes/go-loctogene"
	"github.com/antonybholmes/go-utils"
	"github.com/labstack/echo/v4"
)

const MAX_ANNOTATIONS = 1000

type ReqLocs struct {
	Locations []dna.Location `json:"locations"`
}

type AnnotationResponse struct {
	Message string                 `json:"message"`
	Status  int                    `json:"status"`
	Data    []*gene.GeneAnnotation `json:"data"`
}

func AnnotationRoute(c echo.Context, loctogenedbcache *loctogene.LoctogeneDbCache) error {
	var err error
	locs := new(ReqLocs)

	err = c.Bind(locs)

	if err != nil {
		return c.JSON(http.StatusBadRequest, StatusMessage{Status: http.StatusBadRequest, Message: err.Error()})
	}

	if err != nil {
		return c.JSON(http.StatusBadRequest, StatusMessage{Status: http.StatusBadRequest, Message: err.Error()})
	}

	// limit amount of data returned per request to 1000 entries at a time
	locations := locs.Locations[0:utils.IntMin(len(locs.Locations), MAX_ANNOTATIONS)]

	query, err := ParseGeneQuery(c, c.Param("assembly"), loctogenedbcache)

	if err != nil {
		return c.JSON(http.StatusBadRequest, StatusMessage{Status: http.StatusBadRequest, Message: err.Error()})
	}

	n := ParseN(c)

	tssRegion := ParseTSSRegion(c)

	output := ParseOutput(c)

	annotationDb := gene.NewAnnotateDb(query.Db, tssRegion, n)

	data := []*gene.GeneAnnotation{}

	for _, location := range locations {

		annotations, err := annotationDb.Annotate(&location)

		if err != nil {
			return c.JSON(http.StatusBadRequest, StatusMessage{Status: http.StatusBadRequest, Message: "there was an error with the database query"})
		}

		data = append(data, annotations)
	}

	if output == "text" {
		tsv, err := MakeGeneTable(data, tssRegion)

		if err != nil {
			return c.JSON(http.StatusBadRequest, StatusMessage{Status: http.StatusBadRequest, Message: err.Error()})
		}

		return c.String(http.StatusOK, tsv)
	} else {

		return c.JSON(http.StatusOK, AnnotationResponse{Status: http.StatusOK, Message: "", Data: data})
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
			row = append(row, strconv.Itoa(closestGene.Dist))
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
