package toolsroutes

import (
	"github.com/antonybholmes/go-auth"
	"github.com/antonybholmes/go-edb-server/routes"
	"github.com/labstack/echo/v4"
)

type HashResp struct {
	Password string `json:"password"`
	Hash     string `json:"hash"`
}

func HashedPasswordRoute(c echo.Context) error {

	password := c.QueryParam("password")

	if len(password) == 0 {
		return routes.ErrorReq("password cannot be empty")
	}

	hash := auth.HashPassword(password)

	ret := HashResp{Password: password, Hash: hash}

	return routes.MakeDataPrettyResp(c, "", ret)
}
