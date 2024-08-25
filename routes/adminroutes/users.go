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

type UserStatResp struct {
	Users uint `json:"users"`
}

func UserStatsRoute(c echo.Context) error {

	var req UserListReq

	c.Bind(&req)

	users, err := userdbcache.NumUsers()

	if err != nil {
		return routes.ErrorReq(err)
	}

	resp := UserStatResp{Users: users}

	return routes.MakeDataPrettyResp(c, "", resp)

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
