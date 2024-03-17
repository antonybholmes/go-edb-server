package auth

import (
	"github.com/antonybholmes/go-auth"
	"github.com/antonybholmes/go-auth/userdb"
	"github.com/antonybholmes/go-edb-api/consts"
	"github.com/antonybholmes/go-edb-api/routes"

	"github.com/labstack/echo/v4"
)

// Start passwordless login by sending an email
func ResetPasswordFromUsernameRoute(c echo.Context) error {
	return routes.NewValidator(c).AuthUserFromUsername().VerifiedEmail().Success(func(validator *routes.Validator) error {
		authUser := validator.AuthUser
		req := validator.Req

		otpJwt, err := auth.ResetPasswordToken(c, authUser, consts.JWT_SECRET)

		if err != nil {
			return routes.ErrorReq(err)
		}

		var file string

		if req.CallbackUrl != "" {
			file = "templates/email/password/reset/web.html"
		} else {
			file = "templates/email/password/reset/api.html"
		}

		err = SendEmailWithToken("Password Reset",
			authUser,
			file,
			otpJwt,
			req.CallbackUrl,
			req.Url)

		if err != nil {
			return routes.ErrorReq(err)
		}

		return routes.MakeOkResp(c, "password reset email sent")
	})

	// return routes.ReqBindCB(c, new(auth.EmailOnlyLoginReq), func(c echo.Context, req *auth.EmailOnlyLoginReq) error {
	// 	return routes.AuthUserFromEmailCB(c, req.Email, func(c echo.Context, authUser *auth.AuthUser) error {
	// 		return routes.VerifiedEmailCB(c, authUser, func(c echo.Context, authUser *auth.AuthUser) error {

	// 			otpJwt, err := auth.ResetPasswordToken(c, authUser.Uuid, consts.JWT_SECRET)

	// 			if err != nil {
	// 				return routes.ErrorReq(err)
	// 			}

	// 			var file string

	// 			if req.Url != "" {
	// 				file = "templates/email/password/reset/web.html"
	// 			} else {
	// 				file = "templates/email/password/reset/api.html"
	// 			}

	// 			err = SendEmailWithToken("Password Reset",
	// 				authUser,
	// 				file,
	// 				otpJwt,
	// 				req.CallbackUrl,
	// 				req.Url)

	// 			if err != nil {
	// 				return routes.ErrorReq(err)
	// 			}

	// 			return routes.MakeSuccessResp(c, "password reset email sent", true)
	// 		})
	// 	})
	// })
}

func UpdatePasswordRoute(c echo.Context) error {
	return routes.NewValidator(c).ReqBind().AuthUserFromUuid().Success(func(validator *routes.Validator) error {

		if validator.Claims.Type != auth.TOKEN_TYPE_RESET_PASSWORD {
			return routes.WrongTokentTypeReq()
		}

		err := userdb.SetPassword(validator.AuthUser.Uuid, validator.Req.Password)

		if err != nil {
			return routes.ErrorReq("error setting password")
		}

		return routes.PasswordUpdatedResp(c)
	})
}
