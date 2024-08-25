package authroutes

import (
	"github.com/antonybholmes/go-auth"
	"github.com/antonybholmes/go-auth/jwtgen"
	"github.com/antonybholmes/go-auth/userdbcache"
	"github.com/antonybholmes/go-edb-server/routes"

	"github.com/labstack/echo/v4"
)

func UserSignedInResp(c echo.Context) error {
	return routes.MakeOkPrettyResp(c, "user signed in")
}

func PasswordlessEmailSentResp(c echo.Context) error {
	return routes.MakeOkPrettyResp(c, "passwordless email sent")
}

func UsernamePasswordSignInRoute(c echo.Context) error {
	return routes.NewValidator(c).ParseLoginRequestBody().Success(func(validator *routes.Validator) error {

		if validator.Req.Password == "" {
			return PasswordlessEmailRoute(c, validator)
		}

		authUser, err := userdbcache.FindUserByUsername(validator.Req.Username)

		if err != nil {
			return routes.UserDoesNotExistReq()
		}

		if !authUser.EmailIsVerified {
			return routes.EmailNotVerifiedReq()
		}

		if !authUser.CanLogin() {
			return routes.UserNotAllowedToSignIn()
		}

		err = authUser.CheckPasswordsMatch(validator.Req.Password)

		if err != nil {
			return routes.ErrorReq(err)
		}

		refreshToken, err := jwtgen.RefreshToken(c, authUser.PublicId, authUser.Roles)

		if err != nil {
			return routes.TokenErrorReq()
		}

		accessToken, err := jwtgen.AccessToken(c, authUser.PublicId, authUser.Roles)

		if err != nil {
			return routes.TokenErrorReq()
		}

		return routes.MakeDataPrettyResp(c, "", &routes.LoginResp{RefreshToken: refreshToken, AccessToken: accessToken})
	})
}

// Start passwordless login by sending an email
func PasswordlessEmailRoute(c echo.Context, validator *routes.Validator) error {
	if validator == nil {
		validator = routes.NewValidator(c)
	}

	return validator.LoadAuthUserFromUsername().CheckUserHasVerifiedEmailAddress().Success(func(validator *routes.Validator) error {

		authUser := validator.AuthUser

		passwordlessToken, err := jwtgen.PasswordlessToken(c, authUser.PublicId)

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
			authUser,
			file,
			passwordlessToken,
			validator.Req.CallbackUrl,
			validator.Req.Url)

		//if err != nil {
		//	return routes.ErrorReq(err)
		//}

		return routes.MakeOkPrettyResp(c, "check your email for a passwordless sign in link")
	})
}

func PasswordlessSignInRoute(c echo.Context) error {
	return routes.NewValidator(c).LoadAuthUserFromToken().CheckUserHasVerifiedEmailAddress().Success(func(validator *routes.Validator) error {

		if validator.Claims.Type != auth.TOKEN_TYPE_PASSWORDLESS {
			return routes.WrongTokentTypeReq()
		}

		authUser := validator.AuthUser

		//log.Debug().Msgf("user %v", authUser)

		if !authUser.CanLogin() {
			return routes.UserNotAllowedToSignIn()
		}

		t, err := jwtgen.RefreshToken(c, authUser.PublicId, authUser.Roles)

		if err != nil {
			return routes.TokenErrorReq()
		}

		return routes.MakeDataPrettyResp(c, "", &routes.RefreshTokenResp{RefreshToken: t})
	})
}
