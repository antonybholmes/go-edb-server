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

var SESSION_OPT_24H *sessions.Options

func init() {

	SESSION_OPT_24H = &sessions.Options{
		Path:   "/",
		MaxAge: 86400,
		// http only false to allow js to delete etc on the client side
		HttpOnly: false,
		Secure:   true,
		SameSite: http.SameSiteNoneMode,
	}
}

func SessionUsernamePasswordLoginRoute(c echo.Context) error {
	validator, err := routes.NewValidator(c).ReqBind().Ok()

	if err != nil {
		return err
	}

	if validator.Req.Password == "" {
		return routes.ErrorReq("password required")
	}

	user := validator.Req.Username

	log.Debug().Msgf("session %s", user)

	authUser, err := userdb.FindUserByUsername(user)

	log.Debug().Msgf("session %s", authUser.Email)

	if err != nil {
		email, err := mail.ParseAddress(user)

		if err != nil {
			return routes.ErrorReq("email address not valid")
		}

		// also check if username is valid email and try to login
		// with that
		authUser, err = userdb.FindUserByEmail(email)

		if err != nil {
			return routes.ErrorReq("user does not exist")
		}
	}

	if !authUser.EmailVerified {
		return routes.ErrorReq("email address not verified")
	}

	if !authUser.CanLogin {
		return routes.ErrorReq("user not allowed to login")
	}

	if !authUser.CheckPasswords(validator.Req.Password) {
		return routes.InvalidPasswordReq()
	}

	sess, err := session.Get(SESSION_NAME, c)
	if err != nil {
		return routes.ErrorReq("error creating session")
	}
	sess.Options = SESSION_OPT_24H
	sess.Values[SESSION_UUID] = authUser.Uuid

	sess.Save(c.Request(), c.Response())

	return routes.MakeOkResp(c, "user was signed in")
	//return c.NoContent(http.StatusOK)
}

func SessionNewAccessTokenRoute(c echo.Context) error {
	sess, _ := session.Get(SESSION_NAME, c)
	uuid, _ := sess.Values[SESSION_UUID].(string)

	log.Debug().Msgf("session tokens %s", uuid)

	_, err := userdb.FindUserByUuid(uuid)

	if err != nil {
		return routes.ErrorReq("user does not exist")
	}

	t, err := auth.AccessToken(c, uuid, consts.JWT_SECRET)

	if err != nil {
		return routes.ErrorReq("error creating access token")
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

func SessionUpdateUserInfoRoute(c echo.Context) error {
	sess, _ := session.Get(SESSION_NAME, c)
	uuid, _ := sess.Values[SESSION_UUID].(string)

	req := new(auth.LoginReq)

	err := c.Bind(req)

	if err != nil {
		return routes.ErrorReq("login parameters missing")
	}

	authUser, err := userdb.FindUserByUuid(uuid)

	if err != nil {
		return routes.UserDoesNotExistReq()
	}

	if !authUser.CheckPasswords(req.Password) {
		log.Debug().Msgf("%s", routes.InvalidPasswordReq())
		return routes.InvalidPasswordReq()
	}

	err = userdb.SetUsername(authUser.Uuid, req.Username)

	if err != nil {
		return routes.ErrorReq("error changing username")
	}

	err = userdb.SetName(authUser.Uuid, req.Name)

	if err != nil {
		return routes.ErrorReq("error changing name")
	}

	address, err := mail.ParseAddress(req.Email)

	if err != nil {
		return routes.InvalidEmailReq()
	}

	err = userdb.SetEmail(authUser.Uuid, address)

	if err != nil {
		return routes.ErrorReq("error changing email")
	}

	return routes.MakeOkResp(c, "user info updated")
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

	if !authUser.CheckPasswords(req.Password) {
		return routes.InvalidPasswordReq()
	}

	err = userdb.SetPassword(authUser.Uuid, req.NewPassword)

	if err != nil {
		return routes.ErrorReq("error setting password")
	}

	return routes.PasswordUpdatedResp(c)
}
