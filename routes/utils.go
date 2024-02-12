package routes

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/antonybholmes/go-loctogene"
	"github.com/labstack/echo/v4"
)

const DEFAULT_ASSEMBLY = "grch38"
const DEFAULT_LEVEL = loctogene.Gene

const DEFAULT_CLOSEST_N uint16 = 5

type StatusResp struct {
	Status int `json:"status"`
}

type StatusMessageResp struct {
	StatusResp
	Message string `json:"message"`
}

type DataResp struct {
	StatusResp
	Data interface{} `json:"data"`
}

func JsonRep[V interface{}](c echo.Context, status int, data V) error {
	return c.JSONPretty(status, data, "  ")
}

func MakeBadResp(c echo.Context, err error) error {
	return JsonRep(c, http.StatusBadRequest, StatusMessageResp{StatusResp: StatusResp{Status: http.StatusBadRequest}, Message: err.Error()})
}

func MakeDataResp[V interface{}](c echo.Context, data *V) error {
	return JsonRep(c, http.StatusOK, DataResp{StatusResp: StatusResp{Status: http.StatusOK}, Data: data})
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

func ParseOutput(c echo.Context) string {

	v := c.QueryParam("output")

	if strings.Contains(strings.ToLower(v), "text") {
		return "text"
	} else {
		return "json"
	}
}
