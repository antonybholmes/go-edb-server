package authorization

import (
	"github.com/antonybholmes/go-auth/userdbcache"
	"github.com/antonybholmes/go-edb-server/routes"
	authenticationroutes "github.com/antonybholmes/go-edb-server/routes/authentication"
	"github.com/labstack/echo/v4"
)

func SessionUpdateUserRoute(c echo.Context) error {
	sessionData, err := authenticationroutes.ReadSessionInfo(c)

	if err != nil {
		return routes.ErrorReq(err)
	}

	authUser := sessionData.AuthUser

	return authenticationroutes.NewValidator(c).CheckUsernameIsWellFormed().CheckEmailIsWellFormed().Success(func(validator *authenticationroutes.Validator) error {

		err = userdbcache.SetUserInfo(authUser, validator.LoginBodyReq.Username, validator.LoginBodyReq.FirstName, validator.LoginBodyReq.LastName, false)

		if err != nil {
			return routes.ErrorReq(err)
		}

		return SendUserInfoUpdatedEmail(c, authUser)
	})
}
