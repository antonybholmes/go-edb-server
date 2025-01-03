package authenticationroutes

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

func UserSignedInResp(c echo.Context) error {
	return routes.MakeOkPrettyResp(c, "user signed in")
}

func PasswordlessEmailSentResp(c echo.Context) error {
	return routes.MakeOkPrettyResp(c, "passwordless email sent")
}

func UsernamePasswordSignInRoute(c echo.Context) error {
	return NewValidator(c).ParseLoginRequestBody().Success(func(validator *Validator) error {

		if validator.LoginBodyReq.Password == "" {
			return PasswordlessSigninEmailRoute(c, validator)
		}

		authUser, err := userdbcache.FindUserByUsername(validator.LoginBodyReq.Username)

		if err != nil {
			return routes.UserDoesNotExistReq()
		}

		if authUser.EmailVerifiedAt == 0 {
			return routes.EmailNotVerifiedReq()
		}

		roles, err := userdbcache.UserRoleList(authUser)

		if err != nil {
			return routes.AuthErrorReq("could not get user roles")
		}

		roleClaim := auth.MakeClaim(roles)

		if !auth.CanSignin(roleClaim) {
			return routes.UserNotAllowedToSignIn()
		}

		err = authUser.CheckPasswordsMatch(validator.LoginBodyReq.Password)

		if err != nil {
			return routes.ErrorReq(err)
		}

		refreshToken, err := tokengen.RefreshToken(c, authUser) //auth.MakeClaim(authUser.Roles))

		if err != nil {
			return routes.TokenErrorReq()
		}

		accessToken, err := tokengen.AccessToken(c, authUser.Uuid, roleClaim) //auth.MakeClaim(authUser.Roles))

		if err != nil {
			return routes.TokenErrorReq()
		}

		return routes.MakeDataPrettyResp(c, "", &routes.LoginResp{RefreshToken: refreshToken, AccessToken: accessToken})
	})
}

// Start passwordless login by sending an email
func PasswordlessSigninEmailRoute(c echo.Context, validator *Validator) error {
	if validator == nil {
		validator = NewValidator(c)
	}

	return validator.LoadAuthUserFromUsername().CheckUserHasVerifiedEmailAddress().Success(func(validator *Validator) error {

		authUser := validator.AuthUser

		passwordlessToken, err := tokengen.PasswordlessToken(c, authUser.Uuid, validator.LoginBodyReq.VisitUrl)

		if err != nil {
			return routes.ErrorReq(err)
		}

		// var file string

		// if validator.Req.CallbackUrl != "" {
		// 	file = "templates/email/passwordless/web.html"
		// } else {
		// 	file = "templates/email/passwordless/api.html"
		// }

		// go SendEmailWithToken("Passwordless Sign In",
		// 	authUser,
		// 	file,
		// 	passwordlessToken,
		// 	validator.Req.CallbackUrl,
		// 	validator.Req.VisitUrl)

		email := mailer.RedisQueueEmail{Name: authUser.FirstName,
			To:          authUser.Email,
			Token:       passwordlessToken,
			EmailType:   mailer.REDIS_EMAIL_TYPE_PASSWORDLESS,
			Ttl:         fmt.Sprintf("%d minutes", int(consts.PASSWORDLESS_TOKEN_TTL_MINS.Minutes())),
			CallBackUrl: validator.LoginBodyReq.CallbackUrl,
			//VisitUrl:    validator.Req.VisitUrl
		}
		rdb.PublishEmail(&email)

		//if err != nil {
		//	return routes.ErrorReq(err)
		//}

		return routes.MakeOkPrettyResp(c, "check your email for a passwordless sign in link")
	})
}

func PasswordlessSignInRoute(c echo.Context) error {
	return NewValidator(c).LoadAuthUserFromToken().CheckUserHasVerifiedEmailAddress().Success(func(validator *Validator) error {

		if validator.Claims.Type != auth.PASSWORDLESS_TOKEN {
			return routes.WrongTokentTypeReq()
		}

		authUser := validator.AuthUser

		roles, err := userdbcache.UserRoleList(authUser)

		if err != nil {
			return routes.AuthErrorReq("could not get user roles")
		}

		roleClaim := auth.MakeClaim(roles)

		if !auth.CanSignin(roleClaim) {
			return routes.UserNotAllowedToSignIn()
		}

		t, err := tokengen.RefreshToken(c, authUser) //auth.MakeClaim(authUser.Roles))

		if err != nil {
			return routes.TokenErrorReq()
		}

		return routes.MakeDataPrettyResp(c, "", &routes.RefreshTokenResp{RefreshToken: t})
	})
}
