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
		Path:     "/",
		MaxAge:   86400,
		HttpOnly: true,
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
		return routes.BadReq("empty password: use passwordless")
	}

	user := validator.Req.Username

	log.Debug().Msgf("session %s", user)

	authUser, err := userdb.FindUserByUsername(user)

	log.Debug().Msgf("session %s", authUser.Email)

	if err != nil {
		email, err := mail.ParseAddress(user)

		if err != nil {
			return routes.BadReq("email address not valid")
		}

		// also check if username is valid email and try to login
		// with that
		authUser, err = userdb.FindUserByEmail(email)

		if err != nil {
			return routes.BadReq("user does not exist")
		}
	}

	if !authUser.EmailVerified {
		return routes.BadReq("email address not verified")
	}

	if !authUser.CanLogin {
		return routes.BadReq("user not allowed to login")
	}

	if !authUser.CheckPasswords(validator.Req.Password) {
		return routes.BadReq("incorrect password")
	}

	sess, err := session.Get(SESSION_NAME, c)
	if err != nil {
		return routes.BadReq("error creating session")
	}
	sess.Options = SESSION_OPT_24H
	sess.Values[SESSION_UUID] = authUser.Uuid

	log.Debug().Msgf("session %s", sess.Values)

	sess.Save(c.Request(), c.Response())

	return routes.MakeOkResp(c, "user was signed in")
	//return c.NoContent(http.StatusOK)
}

func SessionNewAccessTokenRoute(c echo.Context) error {
	sess, err := session.Get(SESSION_NAME, c)
	if err != nil {
		return err
	}

	log.Debug().Msgf("%s", sess.ID)

	uuid, ok := sess.Values[SESSION_UUID].(string)

	if !ok {
		return routes.BadReq("cannot get user id from session")
	}

	_, err = userdb.FindUserByUuid(uuid)

	if err != nil {
		return routes.BadReq("user does not exist")
	}

	t, err := auth.AccessToken(c, uuid, consts.JWT_SECRET)

	if err != nil {
		return routes.BadReq("error creating access token")
	}

	return routes.MakeDataResp(c, "", &routes.AccessTokenResp{AccessToken: t})
}
