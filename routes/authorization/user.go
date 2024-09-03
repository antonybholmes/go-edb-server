package authorization

import (
	"github.com/antonybholmes/go-auth"
	"github.com/antonybholmes/go-auth/userdbcache"
	"github.com/antonybholmes/go-edb-server/rdb"
	"github.com/antonybholmes/go-edb-server/routes"
	"github.com/antonybholmes/go-edb-server/routes/authentication"
	"github.com/antonybholmes/go-mailer"
	"github.com/labstack/echo/v4"
)

type NameReq struct {
	Name string `json:"name"`
}

func UserUpdatedResp(c echo.Context) error {
	return routes.MakeOkPrettyResp(c, "account updated")
}

func UpdateUserRoute(c echo.Context) error {
	return authentication.NewValidator(c).ParseLoginRequestBody().LoadAuthUserFromToken().Success(func(validator *authentication.Validator) error {

		authUser := validator.AuthUser

		err := userdbcache.SetUserInfo(authUser.PublicId,
			validator.Req.Username,
			validator.Req.FirstName,
			validator.Req.LastName, nil)

		if err != nil {
			return routes.ErrorReq(err)
		}

		//return SendUserInfoUpdatedEmail(c, authUser)

		// reload user details
		authUser, err = userdbcache.FindUserByPublicId(authUser.PublicId)

		if err != nil {
			return routes.ErrorReq(err)
		}

		// send email notification of change
		email := mailer.RedisQueueEmail{Name: authUser.FirstName,
			To:        authUser.Email,
			EmailType: mailer.REDIS_EMAIL_TYPE_ACCOUNT_UPDATED}
		rdb.PublishEmail(&email)

		// send back updated user to having to do a separate call to get the new data
		return routes.MakeDataPrettyResp(c, "account updated confirmation email sent", authUser)
	})
}

func UserRoute(c echo.Context) error {
	return authentication.NewValidator(c).
		LoadAuthUserFromToken().
		Success(func(validator *authentication.Validator) error {
			return routes.MakeDataPrettyResp(c, "", validator.AuthUser)
		})
}

func SendUserInfoUpdatedEmail(c echo.Context, authUser *auth.AuthUser) error {

	file := "templates/email/account/updated.html"

	go authentication.SendEmailWithToken("Account Updated",
		authUser,
		file,
		"",
		"",
		"")

	//if err != nil {
	//	return routes.ErrorReq(err)
	//}

	return UserUpdatedResp(c)

}
