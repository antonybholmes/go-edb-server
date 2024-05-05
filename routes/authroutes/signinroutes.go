package authroutes

import (
	"github.com/antonybholmes/go-auth"
	"github.com/antonybholmes/go-auth/userdb"
	"github.com/antonybholmes/go-edb-api/consts"
	"github.com/antonybholmes/go-edb-api/routes"

	"github.com/labstack/echo/v4"
)

func UserSignedInResp(c echo.Context) error {
	return routes.MakeOkResp(c, "user signed in")
}

func PasswordlessEmailSentResp(c echo.Context) error {
	return routes.MakeOkResp(c, "passwordless email sent")
}

func UsernamePasswordSignInRoute(c echo.Context) error {
	return routes.NewValidator(c).ParseLoginRequestBody().Success(func(validator *routes.Validator) error {

		if validator.Req.Password == "" {
			return PasswordlessEmailRoute(c, validator)
		}

		authUser, err := userdb.FindUserById(validator.Req.Username)

		if err != nil {
			return routes.UserDoesNotExistReq()
		}

		if !authUser.EmailVerified {
			return routes.EmailNotVerifiedReq()
		}

		if !authUser.CanSignIn {
			return routes.UserNotAllowedToSignIn()
		}

		err = authUser.CheckPasswordsMatch(validator.Req.Password)

		if err != nil {
			return routes.ErrorReq(err)
		}

		refreshToken, err := auth.RefreshToken(c, authUser.Uuid, consts.JWT_PRIVATE_KEY)

		if err != nil {
			return routes.TokenErrorReq()
		}

		accessToken, err := auth.AccessToken(c, authUser.Uuid, consts.JWT_PRIVATE_KEY)

		if err != nil {
			return routes.TokenErrorReq()
		}

		return routes.MakeDataResp(c, "", &routes.LoginResp{RefreshToken: refreshToken, AccessToken: accessToken})
	})
}

// Start passwordless login by sending an email
func PasswordlessEmailRoute(c echo.Context, validator *routes.Validator) error {
	if validator == nil {
		validator = routes.NewValidator(c)
	}

	return validator.LoadAuthUserFromId().CheckUserHasVerifiedEmailAddress().Success(func(validator *routes.Validator) error {

		passwordlessToken, err := auth.PasswordlessToken(c, validator.AuthUser.Uuid, consts.JWT_PRIVATE_KEY)

		if err != nil {
			return routes.ErrorReq(err)
		}

		var file string

		if validator.Req.CallbackUrl != "" {
			file = "templates/email/passwordless/web.html"
		} else {
			file = "templates/email/passwordless/api.html"
		}

		go SendEmailWithToken("Passwordless Sign In",
			validator.AuthUser,
			file,
			passwordlessToken,
			validator.Req.CallbackUrl,
			validator.Req.Url)

		//if err != nil {
		//	return routes.ErrorReq(err)
		//}

		return routes.MakeOkResp(c, "check your email for a passwordless sign in link")
	})
}

func PasswordlessSignInRoute(c echo.Context) error {
	return routes.NewValidator(c).LoadAuthUserFromToken().CheckUserHasVerifiedEmailAddress().Success(func(validator *routes.Validator) error {

		if validator.Claims.Type != auth.TOKEN_TYPE_PASSWORDLESS {
			return routes.WrongTokentTypeReq()
		}

		if !validator.AuthUser.CanSignIn {
			return routes.UserNotAllowedToSignIn()
		}

		t, err := auth.RefreshToken(c, validator.AuthUser.Uuid, consts.JWT_PRIVATE_KEY)

		if err != nil {
			return routes.TokenErrorReq()
		}

		return routes.MakeDataResp(c, "", &routes.RefreshTokenResp{RefreshToken: t})
	})
}
