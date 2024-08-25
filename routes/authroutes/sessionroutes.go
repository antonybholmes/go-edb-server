package authroutes

import (
	"net/http"
	"net/mail"

	"github.com/antonybholmes/go-auth"
	"github.com/antonybholmes/go-auth/userdbcache"
	"github.com/antonybholmes/go-edb-server/consts"
	"github.com/antonybholmes/go-edb-server/routes"
	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
)

var SESSION_OPT_MAX_AGE_ZERO *sessions.Options
var SESSION_OPT_24H *sessions.Options
var SESSION_OPT_MAX_AGE_30D *sessions.Options

const MONTH_SECONDS = 2592000

func init() {

	SESSION_OPT_MAX_AGE_ZERO = &sessions.Options{
		Path:   "/",
		MaxAge: 0,
		// http only false to allow js to delete etc on the client side
		HttpOnly: false,
		Secure:   true,
		SameSite: http.SameSiteNoneMode,
	}

	SESSION_OPT_24H = &sessions.Options{
		Path:   "/",
		MaxAge: 86400,
		// http only false to allow js to delete etc on the client side
		HttpOnly: false,
		Secure:   true,
		SameSite: http.SameSiteNoneMode,
	}

	SESSION_OPT_MAX_AGE_30D = &sessions.Options{
		Path:   "/",
		MaxAge: MONTH_SECONDS,
		// http only false to allow js to delete etc on the client side
		HttpOnly: false,
		Secure:   true,
		SameSite: http.SameSiteNoneMode,
	}
}

func SessionUsernamePasswordSignInRoute(c echo.Context) error {
	validator, err := routes.NewValidator(c).ParseLoginRequestBody().Ok()

	if err != nil {
		return err
	}

	if validator.Req.Password == "" {
		return PasswordlessEmailRoute(c, validator)
	}

	user := validator.Req.Username

	authUser, err := userdbcache.FindUserByUsername(user)

	if err != nil {
		email, err := mail.ParseAddress(user)

		if err != nil {
			return routes.InvalidEmailReq()
		}

		// also check if username is valid email and try to login
		// with that
		authUser, err = userdbcache.FindUserByEmail(email)

		if err != nil {
			return routes.UserDoesNotExistReq()
		}
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

	sess, err := session.Get(routes.SESSION_NAME, c)
	if err != nil {
		return routes.ErrorReq("error creating session")
	}

	if validator.Req.StaySignedIn {
		sess.Options = SESSION_OPT_MAX_AGE_30D
	} else {
		sess.Options = SESSION_OPT_MAX_AGE_ZERO
	}

	sess.Values[routes.SESSION_UUID] = authUser.Uuid

	sess.Save(c.Request(), c.Response())

	return UserSignedInResp(c)
	//return c.NoContent(http.StatusOK)
}

func SessionSignOutRoute(c echo.Context) error {
	sess, err := session.Get(routes.SESSION_NAME, c)
	if err != nil {
		return routes.ErrorReq("error creating session")
	}

	// invalidate by time
	sess.Options.MaxAge = 0

	sess.Save(c.Request(), c.Response())

	return routes.MakeOkPrettyResp(c, "user was signed out")
}

func SessionNewAccessTokenRoute(c echo.Context) error {
	sess, _ := session.Get(routes.SESSION_NAME, c)
	uuid, _ := sess.Values[routes.SESSION_UUID].(string)

	log.Debug().Msgf("session tokens %s", uuid)

	authUser, err := userdbcache.FindUserByUuid(uuid)

	if err != nil {
		return routes.UserDoesNotExistReq()
	}

	t, err := auth.AccessToken(c, uuid, authUser.Permissions, consts.JWT_PRIVATE_KEY)

	if err != nil {
		return routes.TokenErrorReq()
	}

	return routes.MakeDataPrettyResp(c, "", &routes.AccessTokenResp{AccessToken: t})
}

func SessionUserInfoRoute(c echo.Context) error {
	sess, _ := session.Get(routes.SESSION_NAME, c)
	uuid, _ := sess.Values[routes.SESSION_UUID].(string)

	authUser, err := userdbcache.FindUserByUuid(uuid)

	if err != nil {
		return routes.UserDoesNotExistReq()
	}

	return routes.MakeDataPrettyResp(c, "", *authUser)
}

func SessionUpdateUserInfoRoute(c echo.Context) error {
	sess, _ := session.Get(routes.SESSION_NAME, c)
	uuid, _ := sess.Values[routes.SESSION_UUID].(string)

	authUser, err := userdbcache.FindUserByUuid(uuid)

	if err != nil {
		return routes.UserDoesNotExistReq()
	}

	return routes.NewValidator(c).CheckEmailIsWellFormed().Success(func(validator *routes.Validator) error {

		// if !authUser.CheckPasswords(req.Password) {
		// 	log.Debug().Msgf("%s", routes.InvalidPasswordReq())
		// 	return routes.InvalidPasswordReq()
		// }

		err = userdbcache.SetUserInfo(authUser.Uuid, validator.Req.Username, validator.Req.FirstName, validator.Req.LastName)

		if err != nil {
			return routes.ErrorReq(err)
		}

		// err = userdbcache.SetEmailAddress(authUser.Uuid, validator.Address)

		// if err != nil {
		// 	return routes.ErrorReq(err)
		// }

		return SendUserInfoUpdatedEmail(c, authUser)
	})
}

// func SessionUpdatePasswordRoute(c echo.Context) error {
// 	sess, _ := session.Get(SESSION_NAME, c)
// 	uuid, _ := sess.Values[SESSION_UUID].(string)

// 	req := new(auth.NewPasswordReq)

// 	err := c.Bind(req)

// 	if err != nil {
// 		return routes.ErrorReq("login parameters missing")
// 	}

// 	authUser, err := userdbcache.FindUserByUuid(uuid)

// 	if err != nil {
// 		return routes.UserDoesNotExistReq()
// 	}

// 	log.Debug().Msgf("up %d", authUser.Updated)

// 	err = authUser.CheckPasswordsMatch(req.Password)

// 	if err != nil {
// 		return routes.ErrorReq(err)
// 	}

// 	err = userdbcache.SetPassword(authUser.Uuid, req.NewPassword)

// 	if err != nil {
// 		return routes.ErrorReq(err)
// 	}

// 	return SendPasswordEmail(c, authUser, req.NewPassword)
// }

func SessionPasswordlessSignInRoute(c echo.Context) error {

	return routes.NewValidator(c).LoadAuthUserFromToken().CheckUserHasVerifiedEmailAddress().Success(func(validator *routes.Validator) error {

		if validator.Claims.Type != auth.TOKEN_TYPE_PASSWORDLESS {
			return routes.WrongTokentTypeReq()
		}

		authUser := validator.AuthUser

		log.Debug().Msgf("user %v", authUser)

		if !authUser.CanLogin() {
			return routes.UserNotAllowedToSignIn()
		}

		sess, err := session.Get(routes.SESSION_NAME, c)

		if err != nil {
			return routes.ErrorReq("error creating session")
		}

		sess.Options = SESSION_OPT_MAX_AGE_30D
		sess.Values[routes.SESSION_UUID] = authUser.Uuid

		sess.Save(c.Request(), c.Response())

		return UserSignedInResp(c)
	})
}

// Start passwordless login by sending an email
func SessionSendResetPasswordRoute(c echo.Context) error {

	return routes.NewValidator(c).LoadAuthUserFromSession().CheckUserHasVerifiedEmailAddress().Success(func(validator *routes.Validator) error {
		authUser := validator.AuthUser
		req := validator.Req

		otpJwt, err := auth.ResetPasswordToken(c, authUser, consts.JWT_PRIVATE_KEY)

		if err != nil {
			return routes.ErrorReq(err)
		}

		var file string

		if req.CallbackUrl != "" {
			file = "templates/email/password/reset/web.html"
		} else {
			file = "templates/email/password/reset/api.html"
		}

		go SendEmailWithToken("Password Reset",
			authUser,
			file,
			otpJwt,
			req.CallbackUrl,
			req.Url)

		//if err != nil {
		//	return routes.ErrorReq(err)
		//}

		return routes.MakeOkPrettyResp(c, "check your email for a password reset link")
	})
}

// func SessionUpdatePasswordRoute(c echo.Context) error {

// 	return routes.NewValidator(c).LoadAuthUserFromSession().ParseLoginRequestBody().Success(func(validator *routes.Validator) error {

// 		err := userdbcache.SetPassword(validator.AuthUser.Uuid, validator.Req.Password)

// 		if err != nil {
// 			return routes.ErrorReq(err)
// 		}

// 		return SendEmailChangedEmail(c, validator.AuthUser, validator.Req.Password)
// 	})
// }

func SessionSendChangeEmailRoute(c echo.Context) error {

	return routes.NewValidator(c).LoadAuthUserFromSession().ParseLoginRequestBody().Success(func(validator *routes.Validator) error {

		req := validator.Req

		newEmail, err := mail.ParseAddress(req.Email)

		if err != nil {
			return routes.ErrorReq(err)
		}

		otpJwt, err := auth.ChangeEmailToken(c, validator.AuthUser, newEmail, consts.JWT_PRIVATE_KEY)

		if err != nil {
			return routes.ErrorReq(err)
		}

		file := "templates/email/email/change/web.html"

		go BaseSendEmailWithToken("Change Your Email Address",
			validator.AuthUser,
			newEmail,
			file,
			otpJwt,
			req.CallbackUrl,
			req.Url)

		//if err != nil {
		//	return routes.ErrorReq(err)
		//}

		return routes.MakeOkPrettyResp(c, "check your email for a change email link")
	})
}
