package authroutes

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

	err = authUser.CheckPasswordsMatch(validator.Req.Password)

	if err != nil {
		return routes.ErrorReq(err)
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

	passwordlessToken, err := auth.PasswordlessToken(c, validator.AuthUser.Uuid, consts.JWT_SECRET)

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
}
