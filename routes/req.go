package routes

import (
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"
)

type JwtInfo struct {
	Uuid string `json:"uuid"`
	//Name  string `json:"name"`
	Type string `json:"type"`
	//IpAddr  string `json:"ipAddr"`
	Expires string `json:"expires"`
}

type ReqJwt struct {
	Jwt string `json:"jwt"`
}

func InvalidEmailReq() *echo.HTTPError {
	return ErrorReq("invalid email address")
}

func EmailNotVerifiedReq() *echo.HTTPError {
	return ErrorReq("email address not verified")
}

func UserDoesNotExistReq() *echo.HTTPError {
	return ErrorReq("user does not exist")
}

func UserNotAllowedToSignIn() *echo.HTTPError {
	return ErrorReq("user not allowed to sign in")
}

func InvalidUsernameReq() *echo.HTTPError {
	return ErrorReq("invalid username")
}

func PasswordsDoNotMatchReq() *echo.HTTPError {
	return ErrorReq("passwords do not match")
}

func WrongTokentTypeReq() *echo.HTTPError {
	return ErrorReq("wrong token type")
}

func TokenErrorReq() *echo.HTTPError {
	return ErrorReq("token not generated")
}

func ErrorReq(message interface{}) *echo.HTTPError {
	return echo.NewHTTPError(http.StatusBadRequest, message)
}

func AuthErrorReq(message interface{}) *echo.HTTPError {
	return echo.NewHTTPError(http.StatusUnauthorized, message)
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

// get the auth token from the header
func HeaderAuthToken(c echo.Context) (string, error) {

	h := c.Request().Header.Get("Authorization")

	if h == "" {
		return "", errors.New("authorization header not present")
	}

	if !strings.Contains(h, "Bearer") {
		return "", errors.New("bearer not present")
	}

	tokens := strings.Split(h, " ")

	if len(tokens) < 2 {
		return "", errors.New("jwt not present")
	}

	return tokens[1], nil
}
