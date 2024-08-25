package adminroutes

import (
	"github.com/antonybholmes/go-auth/userdbcache"
	"github.com/antonybholmes/go-edb-server/routes"
	"github.com/labstack/echo/v4"
)

type UserListReq struct {
	Offset  int
	Records int
}

func ListUsersRoute(c echo.Context) error {

	var req UserListReq

	c.Bind(&req)

	users, err := userdbcache.Users(req.Offset, req.Records)

	if err != nil {
		return routes.ErrorReq(err)
	}

	return routes.MakeDataPrettyResp(c, "", users)

}
