package authorization

import (
	"github.com/antonybholmes/go-auth/userdbcache"
	"github.com/antonybholmes/go-edb-server/consts"
	"github.com/antonybholmes/go-edb-server/routes"
	authenticationroutes "github.com/antonybholmes/go-edb-server/routes/authentication"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
)

func SessionUpdateUserRoute(c echo.Context) error {
	sess, _ := session.Get(consts.SESSION_NAME, c)
	publicId, _ := sess.Values[authenticationroutes.SESSION_PUBLICID].(string)

	authUser, err := userdbcache.FindUserByPublicId(publicId)

	if err != nil {
		return routes.UserDoesNotExistReq()
	}

	return authenticationroutes.NewValidator(c).CheckUsernameIsWellFormed().CheckEmailIsWellFormed().Success(func(validator *authenticationroutes.Validator) error {

		err = userdbcache.SetUserInfo(authUser.PublicId, validator.LoginBodyReq.Username, validator.LoginBodyReq.FirstName, validator.LoginBodyReq.LastName)

		if err != nil {
			return routes.ErrorReq(err)
		}

		return SendUserInfoUpdatedEmail(c, authUser)
	})
}
