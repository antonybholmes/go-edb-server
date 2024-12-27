package authorization

import (
	"github.com/antonybholmes/go-auth"
	"github.com/antonybholmes/go-auth/userdbcache"
	"github.com/antonybholmes/go-edb-server/rdb"
	"github.com/antonybholmes/go-edb-server/routes"
	authenticationroutes "github.com/antonybholmes/go-edb-server/routes/authentication"

	"github.com/antonybholmes/go-mailer"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
)

type NameReq struct {
	Name string `json:"name"`
}

func UserUpdatedResp(c echo.Context) error {
	return routes.MakeOkPrettyResp(c, "account updated")
}

func UpdateUserRoute(c echo.Context) error {

	return authenticationroutes.NewValidator(c).ParseLoginRequestBody().LoadTokenClaims().Success(func(validator *authenticationroutes.Validator) error {

		//db, err := userdbcache.AutoConn(nil) //not clear on what is needed for the user and password

		publicId := validator.Claims.PublicId

		authUser, err := userdbcache.FindUserByPublicId(publicId)

		if err != nil {
			return routes.ErrorReq(err)
		}

		log.Debug().Msgf("update pub %s", publicId)

		err = userdbcache.SetUserInfo(authUser,
			validator.LoginBodyReq.Username,
			validator.LoginBodyReq.FirstName,
			validator.LoginBodyReq.LastName,
			false)

		if err != nil {
			return routes.ErrorReq(err)
		}

		//return SendUserInfoUpdatedEmail(c, authUser)

		// reload user details
		authUser, err = userdbcache.FindUserByPublicId(publicId)

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
	return authenticationroutes.NewValidator(c).
		LoadAuthUserFromToken().
		Success(func(validator *authenticationroutes.Validator) error {
			return routes.MakeDataPrettyResp(c, "", validator.AuthUser)
		})
}

func SendUserInfoUpdatedEmail(c echo.Context, authUser *auth.AuthUser) error {

	file := "templates/email/account/updated.html"

	go authenticationroutes.SendEmailWithToken("Account Updated",
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
