package authroutes

import (
	"github.com/antonybholmes/go-auth"
	"github.com/antonybholmes/go-auth/userdb"
	"github.com/antonybholmes/go-edb-api/routes"

	"github.com/labstack/echo/v4"
)

type NameReq struct {
	Name string `json:"name"`
}

func AccountUpdatedResp(c echo.Context) error {
	return routes.MakeOkResp(c, "account updated")
}

func UpdateAccountRoute(c echo.Context) error {
	return routes.NewValidator(c).LoadAuthUserFromToken().Success(func(validator *routes.Validator) error {

		authUser := validator.AuthUser

		err := userdb.SetUsername(authUser.Uuid,
			validator.Req.Username)

		if err != nil {
			return routes.ErrorReq(err)
		}

		err = userdb.SetName(authUser.Uuid,
			validator.Req.FirstName,
			validator.Req.LastName)

		if err != nil {
			return routes.ErrorReq(err)
		}

		return SendAccountUpdatedEmail(c, authUser)
	})

	// return routes.ReqBindCB(c, new(auth.UsernameReq), func(c echo.Context, req *auth.UsernameReq) error {
	// 	return routes.IsValidAccessTokenCB(c, func(c echo.Context, claims *auth.JwtCustomClaims) error {
	// 		return routes.AuthUserFromUuidCB(c, claims, func(c echo.Context, claims *auth.JwtCustomClaims, authUser *auth.AuthUser) error {
	// 			err := userdb.SetUsername(authUser.Uuid, req.Username)

	// 			if err != nil {
	// 				return routes.ErrorReq("error setting password")
	// 			}

	// 			return routes.MakeSuccessResp(c, "password updated", true)
	// 		})
	// 	})
	// })

}

// func UpdateNameRoute(c echo.Context) error {
// 	return routes.NewValidator(c).
// 		IsValidAccessToken().
// 		AuthUserFromUuid().
// 		ReqBind().
// 		Success(func(validator *routes.Validator) error {

// 			err := userdb.SetName(validator.AuthUser.Uuid, validator.Req.Name)

// 			if err != nil {
// 				return routes.ErrorReq("error setting password")
// 			}

// 			return routes.MakeOkResp(c, "name updated")
// 		})
// }

func UserInfoRoute(c echo.Context) error {
	return routes.NewValidator(c).
		LoadAuthUserFromToken().
		Success(func(validator *routes.Validator) error {

			return routes.MakeDataResp(c, "", *validator.AuthUser)
		})
}

func SendAccountUpdatedEmail(c echo.Context, authUser *auth.AuthUser) error {

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
