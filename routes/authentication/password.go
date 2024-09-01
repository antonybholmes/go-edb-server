package authentication

import (
	"fmt"

	"github.com/antonybholmes/go-auth"
	"github.com/antonybholmes/go-auth/tokengen"
	"github.com/antonybholmes/go-auth/userdbcache"
	"github.com/antonybholmes/go-edb-server/consts"
	"github.com/antonybholmes/go-edb-server/rdb"
	"github.com/antonybholmes/go-edb-server/routes"

	"github.com/antonybholmes/go-mailer"
	"github.com/labstack/echo/v4"
)

func PasswordUpdatedResp(c echo.Context) error {
	return routes.MakeOkPrettyResp(c, "password updated")
}

// Start passwordless login by sending an email
func SendResetPasswordFromUsernameEmailRoute(c echo.Context) error {
	return NewValidator(c).LoadAuthUserFromUsername().CheckUserHasVerifiedEmailAddress().Success(func(validator *Validator) error {
		authUser := validator.AuthUser
		req := validator.Req

		otpToken, err := tokengen.ResetPasswordToken(c, authUser)

		if err != nil {
			return routes.ErrorReq(err)
		}

		// var file string

		// if req.CallbackUrl != "" {
		// 	file = "templates/email/password/reset/web.html"
		// } else {
		// 	file = "templates/email/password/reset/api.html"
		// }

		// go authentication.SendEmailWithToken("Password Reset",
		// 	authUser,
		// 	file,
		// 	otpToken,
		// 	req.CallbackUrl,
		// 	req.VisitUrl)

		email := mailer.RedisQueueEmail{Name: authUser.FirstName,
			To:          authUser.Email,
			Token:       otpToken,
			EmailType:   mailer.REDIS_EMAIL_TYPE_PASSWORD_RESET,
			Ttl:         fmt.Sprintf("%d minutes", int(consts.SHORT_TTL_MINS.Minutes())),
			CallBackUrl: req.CallbackUrl,
			VisitUrl:    req.VisitUrl}
		rdb.PublishEmail(&email)

		//if err != nil {
		//	return routes.ErrorReq(err)
		//}

		return routes.MakeOkPrettyResp(c, "check your email for a password reset link")
	})
}

func UpdatePasswordRoute(c echo.Context) error {
	return NewValidator(c).ParseLoginRequestBody().LoadAuthUserFromToken().Success(func(validator *Validator) error {

		if validator.Claims.Type != auth.RESET_PASSWORD_TOKEN {
			return routes.WrongTokentTypeReq()
		}

		authUser := validator.AuthUser

		err := auth.CheckOTPValid(authUser, validator.Claims.Otp)

		if err != nil {
			return routes.ErrorReq(err)
		}

		err = userdbcache.SetPassword(authUser.PublicId, validator.Req.Password, nil)

		if err != nil {
			return routes.ErrorReq(err)
		}

		//return SendPasswordEmail(c, validator.AuthUser, validator.Req.Password)

		email := mailer.RedisQueueEmail{Name: authUser.FirstName,
			To:        authUser.Email,
			EmailType: mailer.REDIS_EMAIL_TYPE_PASSWORD_UPDATED}
		rdb.PublishEmail(&email)

		return routes.MakeOkPrettyResp(c, "password updated confirmation email sent")
	})
}

func SendPasswordEmail(c echo.Context, authUser *auth.AuthUser, password string) error {

	var file string

	if password == "" {
		file = "templates/email/password/switch-to-passwordless.html"
	} else {
		file = "templates/email/password/updated.html"
	}

	go SendEmailWithToken("Password Updated",
		authUser,
		file,
		"",
		"",
		"")

	//if err != nil {
	//	return routes.ErrorReq(err)
	//}

	return PasswordUpdatedResp(c)

}
