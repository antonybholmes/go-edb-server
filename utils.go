package main

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"
)

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
	return c.JSONPretty(status, data, " ")
}

func MakeBadResp(c echo.Context, err error) error {
	return JsonRep(c, http.StatusBadRequest, StatusMessageResp{StatusResp: StatusResp{Status: http.StatusBadRequest}, Message: err.Error()})
}

func BadReq(message ...interface{}) *echo.HTTPError {
	return echo.NewHTTPError(http.StatusBadRequest, message)
}

func MakeDataResp[V interface{}](c echo.Context, data *V) error {
	return JsonRep(c, http.StatusOK, DataResp{StatusResp: StatusResp{Status: http.StatusOK}, Data: data})
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

func ParseN(c echo.Context, defaultN uint16) uint16 {

	v := c.QueryParam("n")

	if v == "" {
		return defaultN
	}

	n, err := strconv.ParseUint(v, 10, 0)

	if err != nil {
		return defaultN
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
