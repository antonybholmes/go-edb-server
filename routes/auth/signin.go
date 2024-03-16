package auth

import (
	"net/mail"

	"github.com/antonybholmes/go-auth"
	"github.com/antonybholmes/go-auth/userdb"
	"github.com/antonybholmes/go-edb-api/consts"
	"github.com/antonybholmes/go-edb-api/routes"
	"github.com/labstack/echo/v4"
)

func UsernamePasswordSignInRoute(c echo.Context) error {
	validator, err := routes.NewValidator(c).ReqBind().Ok()

	if err != nil {
		return err
	}

	if validator.Req.Password == "" {
		return PasswordlessEmailRoute(c, validator)
	}

	authUser, err := userdb.FindUserByUsername(validator.Req.Username)

	if err != nil {
		email, err := mail.ParseAddress(validator.Req.Username)

		if err != nil {
			return routes.InvalidEmailReq()
		}

		// also check if username is valid email and try to login
		// with that
		authUser, err = userdb.FindUserByEmail(email)

		if err != nil {
			return routes.UserDoesNotExistReq()
		}
	}

	if !authUser.EmailVerified {
		return routes.EmailNotVerifiedReq()
	}

	if !authUser.CanSignIn {
		return routes.UserNotAllowedToSignIn()
	}

	if !authUser.CheckPasswords(validator.Req.Password) {
		return routes.InvalidPasswordReq()
	}

	refreshToken, err := auth.RefreshToken(c, authUser.Uuid, consts.JWT_SECRET)

	if err != nil {
		return routes.TokenErrorReq()
	}

	accessToken, err := auth.AccessToken(c, authUser.Uuid, consts.JWT_SECRET)

	if err != nil {
		return routes.TokenErrorReq()
	}

	return routes.MakeDataResp(c, "", &routes.LoginResp{RefreshToken: refreshToken, AccessToken: accessToken})

	// return routes.ReqBindCB(c, new(auth.UsernamePasswordLoginReq), func(c echo.Context, req *auth.UsernamePasswordLoginReq) error {

	// 	if req.Password == "" {
	// 		return routes.ErrorReq("empty password: use passwordless")
	// 	}

	// 	authUser, err := userdb.FindUserByUsername(req.Username)

	// 	if err != nil {
	// 		email, err := mail.ParseAddress(req.Username)

	// 		if err != nil {
	// 			return routes.ErrorReq("email address not valid")
	// 		}

	// 		// also check if username is valid email and try to login
	// 		// with that
	// 		authUser, err = userdb.FindUserByEmail(email)

	// 		if err != nil {
	// 			return routes.ErrorReq("user does not exist")
	// 		}
	// 	}

	// 	if !authUser.EmailVerified {
	// 		return routes.ErrorReq("email address not verified")
	// 	}

	// 	if !authUser.CanAuth {
	// 		return routes.ErrorReq("user not allowed tokens")
	// 	}

	// 	if !authUser.CheckPasswords(req.Password) {
	// 		return routes.ErrorReq("incorrect password")
	// 	}

	// 	refreshToken, err := auth.RefreshToken(c, authUser.Uuid, consts.JWT_SECRET)

	// 	if err != nil {
	// 		return routes.ErrorReq("error signing token")
	// 	}

	// 	accessToken, err := auth.AccessToken(c, authUser.Uuid, consts.JWT_SECRET)

	// 	if err != nil {
	// 		return routes.ErrorReq("error signing token")
	// 	}

	// 	return routes.MakeDataResp(c, "", &routes.LoginResp{RefreshToken: refreshToken, AccessToken: accessToken})
	// })
}

// Start passwordless login by sending an email
func PasswordlessEmailRoute(c echo.Context, validator *routes.Validator) error {
	if validator == nil {
		validator = routes.NewValidator(c)
	}

	validator, err := validator.AuthUserFromUsername().VerifiedEmail().Ok()

	if err != nil {
		return err
	}

	otpJwt, err := auth.PasswordlessToken(c, validator.AuthUser.Uuid, consts.JWT_SECRET)

	if err != nil {
		return routes.ErrorReq(err)
	}

	var file string

	if validator.Req.CallbackUrl != "" {
		file = "templates/email/passwordless/web.html"
	} else {
		file = "templates/email/passwordless/api.html"
	}

	err = SendEmailWithToken("Passwordless Sign In",
		validator.AuthUser,
		file,
		otpJwt,
		validator.Req.CallbackUrl,
		validator.Req.Url)

	if err != nil {
		return routes.ErrorReq(err)
	}

	return routes.MakeOkResp(c, "passwordless email sent")

	// return routes.ReqBindCB(c, new(auth.EmailOnlyLoginReq), func(c echo.Context, req *auth.EmailOnlyLoginReq) error {
	// 	return routes.AuthUserFromEmailCB(c, req.Email, func(c echo.Context, authUser *auth.AuthUser) error {
	// 		return routes.VerifiedEmailCB(c, authUser, func(c echo.Context, authUser *auth.AuthUser) error {

	// 			otpJwt, err := auth.PasswordlessToken(c, authUser.Uuid, consts.JWT_SECRET)

	// 			if err != nil {
	// 				return routes.ErrorReq(err)
	// 			}

	// 			var file string

	// 			if req.Url != "" {
	// 				file = "templates/email/passwordless/web.html"
	// 			} else {
	// 				file = "templates/email/passwordless/api.html"
	// 			}

	// 			err = SendEmailWithToken("Passwordless Login",
	// 				authUser,
	// 				file,
	// 				otpJwt,
	// 				req.CallbackUrl,
	// 				req.Url)

	// 			if err != nil {
	// 				return routes.ErrorReq(err)
	// 			}

	// 			return routes.MakeSuccessResp(c, "passwordless email sent", true)
	// 		})
	// 	})
	// })
}

func PasswordlessSignInRoute(c echo.Context) error {
	return routes.NewValidator(c).AuthUserFromUuid().VerifiedEmail().Success(func(validator *routes.Validator) error {

		if validator.Claims.Type != auth.TOKEN_TYPE_PASSWORDLESS {
			return routes.WrongTokentTypeReq()
		}

		if !validator.AuthUser.CanSignIn {
			return routes.UserNotAllowedToSignIn()
		}

		t, err := auth.RefreshToken(c, validator.AuthUser.Uuid, consts.JWT_SECRET)

		if err != nil {
			return routes.TokenErrorReq()
		}

		return routes.MakeDataResp(c, "", &routes.RefreshTokenResp{RefreshToken: t})
	})

	// return routes.AuthUserFromUuidCB(c, nil, func(c echo.Context, claims *auth.JwtCustomClaims, authUser *auth.AuthUser) error {
	// 	return routes.VerifiedEmailCB(c, authUser, func(c echo.Context, authUser *auth.AuthUser) error {
	// 		if claims.Type != auth.TOKEN_TYPE_PASSWORDLESS {
	// 			return routes.ErrorReq("wrong token type")
	// 		}

	// 		if !authUser.CanAuth {
	// 			return routes.ErrorReq("user not allowed tokens")
	// 		}

	// 		t, err := auth.RefreshToken(c, authUser.Uuid, consts.JWT_SECRET)

	// 		if err != nil {
	// 			return routes.ErrorReq("error signing token")
	// 		}

	// 		return routes.MakeDataResp(c, "", &routes.RefreshTokenResp{RefreshToken: t})
	// 	})
	// })

}
