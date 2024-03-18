package auth

import (
	"net/http"
	"net/mail"

	"github.com/antonybholmes/go-auth"
	"github.com/antonybholmes/go-auth/userdb"
	"github.com/antonybholmes/go-edb-api/consts"
	"github.com/antonybholmes/go-edb-api/routes"
	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
)

const SESSION_NAME string = "session"
const SESSION_UUID string = "uuid"

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
	validator, err := routes.NewValidator(c).ReqBind().Ok()

	if err != nil {
		return err
	}

	if validator.Req.Password == "" {
		return PasswordlessEmailRoute(c, validator)
	}

	user := validator.Req.Username

	log.Debug().Msgf("session %s", user)

	authUser, err := userdb.FindUserByUsername(user)

	log.Debug().Msgf("session %s", authUser.Email)

	if err != nil {
		email, err := mail.ParseAddress(user)

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

	sess, err := session.Get(SESSION_NAME, c)
	if err != nil {
		return routes.ErrorReq("error creating session")
	}

	if validator.Req.StaySignedIn {
		sess.Options = SESSION_OPT_MAX_AGE_30D
	} else {
		sess.Options = SESSION_OPT_MAX_AGE_ZERO
	}

	sess.Values[SESSION_UUID] = authUser.Uuid

	sess.Save(c.Request(), c.Response())

	return routes.UserSignedInResp(c)
	//return c.NoContent(http.StatusOK)
}

func SessionSignOutRoute(c echo.Context) error {
	sess, err := session.Get(SESSION_NAME, c)
	if err != nil {
		return routes.ErrorReq("error creating session")
	}

	// invalidate by time
	sess.Options.MaxAge = 0

	sess.Save(c.Request(), c.Response())

	return routes.MakeOkResp(c, "user was signed out")
}

func SessionNewAccessTokenRoute(c echo.Context) error {
	sess, _ := session.Get(SESSION_NAME, c)
	uuid, _ := sess.Values[SESSION_UUID].(string)

	log.Debug().Msgf("session tokens %s", uuid)

	_, err := userdb.FindUserByUuid(uuid)

	if err != nil {
		return routes.UserDoesNotExistReq()
	}

	t, err := auth.AccessToken(c, uuid, consts.JWT_SECRET)

	if err != nil {
		return routes.TokenErrorReq()
	}

	return routes.MakeDataResp(c, "", &routes.AccessTokenResp{AccessToken: t})
}

func SessionUserInfoRoute(c echo.Context) error {
	sess, _ := session.Get(SESSION_NAME, c)
	uuid, _ := sess.Values[SESSION_UUID].(string)

	authUser, err := userdb.FindUserByUuid(uuid)

	if err != nil {
		return routes.UserDoesNotExistReq()
	}

	return routes.MakeDataResp(c, "", *authUser.ToPublicUser())
}

func SessionUpdateAccountRoute(c echo.Context) error {
	sess, _ := session.Get(SESSION_NAME, c)
	uuid, _ := sess.Values[SESSION_UUID].(string)

	return routes.NewValidator(c).ValidEmail().Success(func(validator *routes.Validator) error {

		authUser, err := userdb.FindUserByUuid(uuid)

		if err != nil {
			return routes.UserDoesNotExistReq()
		}

		// if !authUser.CheckPasswords(req.Password) {
		// 	log.Debug().Msgf("%s", routes.InvalidPasswordReq())
		// 	return routes.InvalidPasswordReq()
		// }

		err = userdb.SetUsername(authUser.Uuid, validator.Req.Username)

		if err != nil {
			return routes.ErrorReq(err)
		}

		err = userdb.SetName(authUser.Uuid, validator.Req.Name)

		if err != nil {
			return routes.ErrorReq(err)
		}

		err = userdb.SetEmailAddress(authUser.Uuid, validator.Address)

		if err != nil {
			return routes.ErrorReq(err)
		}

		return routes.MakeOkResp(c, "account updated")
	})
}

func SessionUpdatePasswordRoute(c echo.Context) error {
	sess, _ := session.Get(SESSION_NAME, c)
	uuid, _ := sess.Values[SESSION_UUID].(string)

	req := new(auth.NewPasswordReq)

	err := c.Bind(req)

	if err != nil {
		return routes.ErrorReq("login parameters missing")
	}

	authUser, err := userdb.FindUserByUuid(uuid)

	if err != nil {
		return routes.UserDoesNotExistReq()
	}

	err = authUser.CheckPasswordsMatch(req.Password)

	if err != nil {
		return routes.ErrorReq(err)
	}

	err = userdb.SetPassword(authUser.Uuid, req.NewPassword)

	if err != nil {
		return routes.ErrorReq(err)
	}

	return SendPasswordEmail(c, authUser, req.NewPassword)
}

func SessionPasswordlessSignInRoute(c echo.Context) error {

	return routes.NewValidator(c).AuthUserFromUuid().VerifiedEmail().Success(func(validator *routes.Validator) error {

		if validator.Claims.Type != auth.TOKEN_TYPE_PASSWORDLESS {
			return routes.WrongTokentTypeReq()
		}

		if !validator.AuthUser.CanSignIn {
			return routes.UserNotAllowedToSignIn()
		}

		sess, err := session.Get(SESSION_NAME, c)

		if err != nil {
			return routes.ErrorReq("error creating session")
		}

		sess.Options = SESSION_OPT_MAX_AGE_30D
		sess.Values[SESSION_UUID] = validator.AuthUser.Uuid

		sess.Save(c.Request(), c.Response())

		return routes.UserSignedInResp(c)
	})
}
