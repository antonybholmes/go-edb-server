package authorization

import (
	"github.com/antonybholmes/go-auth/userdbcache"
	"github.com/antonybholmes/go-edb-server/consts"
	"github.com/antonybholmes/go-edb-server/routes"
	"github.com/antonybholmes/go-edb-server/routes/authentication"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
)

func SessionUpdateUserRoute(c echo.Context) error {
	sess, _ := session.Get(consts.SESSION_NAME, c)
	publicId, _ := sess.Values[authentication.SESSION_PUBLICID].(string)

	authUser, err := userdbcache.FindUserByPublicId(publicId, nil)

	if err != nil {
		return routes.UserDoesNotExistReq()
	}

	return authentication.NewValidator(c).CheckUsernameIsWellFormed().CheckEmailIsWellFormed().Success(func(validator *authentication.Validator) error {

		err = userdbcache.SetUserInfo(authUser.PublicId, validator.Req.Username, validator.Req.FirstName, validator.Req.LastName, nil)

		if err != nil {
			return routes.ErrorReq(err)
		}

		return SendUserInfoUpdatedEmail(c, authUser)
	})
}
