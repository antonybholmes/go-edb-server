package authroutes

import (
	"net/mail"

	"github.com/antonybholmes/go-auth"
	"github.com/antonybholmes/go-auth/userdb"
	"github.com/antonybholmes/go-edb-api/consts"
	"github.com/antonybholmes/go-edb-api/routes"

	"github.com/labstack/echo/v4"
)

func EmailUpdatedResp(c echo.Context) error {
	return routes.MakeOkResp(c, "email updated")
}

// Start passwordless login by sending an email
func SendChangeEmailRoute(c echo.Context) error {
	return routes.NewValidator(c).ParseLoginRequestBody().LoadAuthUserFromToken().Success(func(validator *routes.Validator) error {
		authUser := validator.AuthUser
		req := validator.Req

		email, err := mail.ParseAddress(req.Email)

		if err != nil {
			return routes.ErrorReq(err)
		}

		otpJwt, err := auth.ChangeEmailToken(c, authUser, email, consts.JWT_SECRET)

		if err != nil {
			return routes.ErrorReq(err)
		}

		var file string

		if req.CallbackUrl != "" {
			file = "templates/email/email/change/web.html"
		} else {
			file = "templates/email/email/change/api.html"
		}

		go SendEmailWithToken("Update Email",
			authUser,
			file,
			otpJwt,
			req.CallbackUrl,
			req.Url)

		//if err != nil {
		//	return routes.ErrorReq(err)
		//}

		return routes.MakeOkResp(c, "check your email for a change email link")
	})
}

func ChangePasswordRoute(c echo.Context) error {
	return routes.NewValidator(c).ParseLoginRequestBody().LoadAuthUserFromToken().Success(func(validator *routes.Validator) error {

		if validator.Claims.Type != auth.TOKEN_TYPE_CHANGE_EMAIL {
			return routes.WrongTokentTypeReq()
		}

		err := auth.CheckOtpValid(validator.AuthUser, validator.Claims.Otp)

		if err != nil {
			return routes.ErrorReq(err)
		}

		err = userdb.SetEmail(validator.AuthUser.Uuid, validator.Claims.Data)

		if err != nil {
			return routes.ErrorReq(err)
		}

		return SendEmailChangedEmail(c, validator.AuthUser, validator.Req.Password)
	})
}

func SendEmailChangedEmail(c echo.Context, authUser *auth.AuthUser, password string) error {

	file := "templates/email/email/updated.html"

	go SendEmailWithToken("Email Address Changed",
		authUser,
		file,
		"",
		"",
		"")

	//if err != nil {
	//	return routes.ErrorReq(err)
	//}

	return EmailUpdatedResp(c)

}
