package routes

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"
)

type JwtResp struct {
	Jwt string `json:"jwt"`
}

type JwtInfo struct {
	Uuid string `json:"uuid"`
	//Name  string `json:"name"`
	Type    string `json:"type"`
	IpAddr  string `json:"ipAddr"`
	Expires string `json:"expires"`
}

type ReqJwt struct {
	Jwt string `json:"jwt"`
}

type StatusResp struct {
	Status int `json:"status"`
}

type StatusMessageResp struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

type DataResp struct {
	StatusMessageResp
	Data interface{} `json:"data"`
}

type SuccessResp struct {
	Success bool `json:"success"`
}

type ValidResp struct {
	Valid bool `json:"valid"`
}

func JsonRep[V any](c echo.Context, status int, data V) error {
	return c.JSONPretty(status, data, " ")
}

// func MakeBadResp(c echo.Context, err error) error {
// 	return JsonRep(c, http.StatusBadRequest, StatusResp{StatusResp: StatusResp{Status: http.StatusBadRequest}, Message: err.Error()})
// }

func BadReq(message interface{}) *echo.HTTPError {
	return echo.NewHTTPError(http.StatusBadRequest, message)
}

func MakeDataResp[V any](c echo.Context, message string, data V) error {
	return JsonRep(c, http.StatusOK, DataResp{StatusMessageResp: StatusMessageResp{Status: http.StatusOK, Message: message}, Data: data})
}

func MakeValidResp(c echo.Context, message string, valid bool) error {
	return MakeDataResp(c, message, &ValidResp{Valid: valid})
}

func MakeSuccessResp(c echo.Context, message string, success bool) error {
	return MakeDataResp(c, message, &SuccessResp{Success: success})
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
