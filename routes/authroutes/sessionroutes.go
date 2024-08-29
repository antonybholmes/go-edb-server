package authroutes

import (
	"net/http"
	"net/mail"

	"github.com/antonybholmes/go-auth"
	"github.com/antonybholmes/go-auth/jwtgen"
	"github.com/antonybholmes/go-auth/userdbcache"
	"github.com/antonybholmes/go-edb-server/consts"
	"github.com/antonybholmes/go-edb-server/routes"
	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
)

var SESSION_OPT_MAX_AGE_ZERO *sessions.Options
var SESSION_OPT_24H *sessions.Options
var SESSION_OPT_MAX_AGE_30D *sessions.Options

const MONTH_SECONDS = 2592000

func init() {

	// HttpOnly and Secure are disabled so we can use them
	// cross domain for testing
	// http only false to allow js to delete etc on the client side

	// For sessions that should end when browser closes
	SESSION_OPT_MAX_AGE_ZERO = &sessions.Options{
		Path:     "/",
		MaxAge:   0,
		HttpOnly: false,
		Secure:   true,
		SameSite: http.SameSiteNoneMode,
	}

	SESSION_OPT_24H = &sessions.Options{
		Path:     "/",
		MaxAge:   86400,
		HttpOnly: false,
		Secure:   true,
		SameSite: http.SameSiteNoneMode,
	}

	SESSION_OPT_MAX_AGE_30D = &sessions.Options{
		Path:     "/",
		MaxAge:   MONTH_SECONDS,
		HttpOnly: false,
		Secure:   true,
		SameSite: http.SameSiteNoneMode,
	}
}

// Validate the passwordless token we generated and create
// a user session. The session acts as a refresh token and
// can be used to generate access tokens to use resources
func SessionPasswordlessSignInRoute(c echo.Context) error {

	return routes.NewValidator(c).LoadAuthUserFromToken().CheckUserHasVerifiedEmailAddress().Success(func(validator *routes.Validator) error {

		if validator.Claims.Type != auth.TOKEN_TYPE_PASSWORDLESS {
			return routes.WrongTokentTypeReq()
		}

		authUser := validator.AuthUser

		//log.Debug().Msgf("user %v", authUser)

		if !authUser.CanLogin() {
			return routes.UserNotAllowedToSignIn()
		}

		sess, err := session.Get(consts.SESSION_NAME, c)

		if err != nil {
			return routes.ErrorReq("error creating session")
		}

		// set session options such as if cookie secure and how long it
		// persists
		sess.Options = SESSION_OPT_MAX_AGE_30D
		sess.Values[routes.SESSION_PUBLICID] = authUser.PublicId
		sess.Values[routes.SESSION_ROLES] = auth.MakeClaim(authUser.Roles)

		sess.Save(c.Request(), c.Response())

		return UserSignedInResp(c)
	})
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

	sess, err := session.Get(consts.SESSION_NAME, c)
	if err != nil {
		return routes.ErrorReq("error creating session")
	}

	if validator.Req.StaySignedIn {
		sess.Options = SESSION_OPT_MAX_AGE_30D
	} else {
		sess.Options = SESSION_OPT_MAX_AGE_ZERO
	}

	sess.Values[routes.SESSION_PUBLICID] = authUser.PublicId
	sess.Values[routes.SESSION_ROLES] = auth.MakeClaim(authUser.Roles)

	sess.Save(c.Request(), c.Response())

	return UserSignedInResp(c)
	//return c.NoContent(http.StatusOK)
}

func SessionSignOutRoute(c echo.Context) error {
	sess, err := session.Get(consts.SESSION_NAME, c)
	if err != nil {
		return routes.ErrorReq("error creating session")
	}

	// invalidate by time
	sess.Options.MaxAge = 0

	sess.Save(c.Request(), c.Response())

	return routes.MakeOkPrettyResp(c, "user was signed out")
}

func SessionNewAccessJwtRoute(c echo.Context) error {
	sess, _ := session.Get(consts.SESSION_NAME, c)
	publicId, _ := sess.Values[routes.SESSION_PUBLICID].(string)
	roles, _ := sess.Values[routes.SESSION_ROLES].(string)

	t, err := jwtgen.AccessToken(c, publicId, roles)

	if err != nil {
		return routes.TokenErrorReq()
	}

	return routes.MakeDataPrettyResp(c, "", &routes.AccessTokenResp{AccessToken: t})
}

func SessionUserRoute(c echo.Context) error {
	sess, _ := session.Get(consts.SESSION_NAME, c)
	publicId, _ := sess.Values[routes.SESSION_PUBLICID].(string)

	authUser, err := userdbcache.FindUserByPublicId(publicId)

	if err != nil {
		return routes.UserDoesNotExistReq()
	}

	return routes.MakeDataPrettyResp(c, "", *authUser)
}

func SessionUpdateUserRoute(c echo.Context) error {
	sess, _ := session.Get(consts.SESSION_NAME, c)
	publicId, _ := sess.Values[routes.SESSION_PUBLICID].(string)

	authUser, err := userdbcache.FindUserByPublicId(publicId)

	if err != nil {
		return routes.UserDoesNotExistReq()
	}

	return routes.NewValidator(c).CheckUsernameIsWellFormed().CheckEmailIsWellFormed().Success(func(validator *routes.Validator) error {

		// if !authUser.CheckPasswords(req.Password) {
		// 	log.Debug().Msgf("%s", routes.InvalidPasswordReq())
		// 	return routes.InvalidPasswordReq()
		// }

		err = userdbcache.SetUserInfo(authUser.PublicId, validator.Req.Username, validator.Req.FirstName, validator.Req.LastName, nil)

		if err != nil {
			return routes.ErrorReq(err)
		}

		// err = userdbcache.SetEmailAddress(authUser.PublicId, validator.Address)

		// if err != nil {
		// 	return routes.ErrorReq(err)
		// }

		return SendUserInfoUpdatedEmail(c, authUser)
	})
}

// Start passwordless login by sending an email
func SessionSendResetPasswordEmailRoute(c echo.Context) error {

	return routes.NewValidator(c).LoadAuthUserFromSession().CheckUserHasVerifiedEmailAddress().Success(func(validator *routes.Validator) error {

		authUser := validator.AuthUser
		req := validator.Req

		otpJwt, err := jwtgen.ResetPasswordToken(c, authUser)

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

// 		err := userdbcache.SetPassword(validator.AuthUser.PublicId, validator.Req.Password)

// 		if err != nil {
// 			return routes.ErrorReq(err)
// 		}

// 		return SendEmailChangedEmail(c, validator.AuthUser, validator.Req.Password)
// 	})
// }

func SessionSendResetEmailEmailRoute(c echo.Context) error {

	return routes.NewValidator(c).LoadAuthUserFromSession().ParseLoginRequestBody().Success(func(validator *routes.Validator) error {

		req := validator.Req

		newEmail, err := mail.ParseAddress(req.Email)

		if err != nil {
			return routes.ErrorReq(err)
		}

		otpJwt, err := jwtgen.ResetEmailToken(c, validator.AuthUser, newEmail)

		if err != nil {
			return routes.ErrorReq(err)
		}

		file := "templates/email/email/reset/web.html"

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
