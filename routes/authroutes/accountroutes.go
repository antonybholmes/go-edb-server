package authroutes

import (
	"github.com/antonybholmes/go-auth"
	"github.com/antonybholmes/go-auth/userdbcache"
	"github.com/antonybholmes/go-edb-server/routes"

	"github.com/labstack/echo/v4"
)

type NameReq struct {
	Name string `json:"name"`
}

func AccountUpdatedResp(c echo.Context) error {
	return routes.MakeOkPrettyResp(c, "account updated")
}

func UpdateAccountRoute(c echo.Context) error {
	return routes.NewValidator(c).LoadAuthUserFromToken().Success(func(validator *routes.Validator) error {

		authUser := validator.AuthUser

		err := userdbcache.SetUsername(authUser.PublicId,
			validator.Req.Username)

		if err != nil {
			return routes.ErrorReq(err)
		}

		err = userdbcache.SetName(authUser.PublicId,
			validator.Req.FirstName,
			validator.Req.LastName)

		if err != nil {
			return routes.ErrorReq(err)
		}

		return SendUserInfoUpdatedEmail(c, authUser)
	})
}

func UserInfoRoute(c echo.Context) error {
	return routes.NewValidator(c).
		LoadAuthUserFromToken().
		Success(func(validator *routes.Validator) error {

			return routes.MakeDataPrettyResp(c, "", validator.AuthUser)
		})
}

func SendUserInfoUpdatedEmail(c echo.Context, authUser *auth.AuthUser) error {

	var file = "templates/email/account/updated.html"

	go SendEmailWithToken("Account Updated",
		authUser,
		file,
		"",
		"",
		"")

	//if err != nil {
	//	return routes.ErrorReq(err)
	//}

	return AccountUpdatedResp(c)

}
