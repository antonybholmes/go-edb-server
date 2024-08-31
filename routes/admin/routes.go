package adminroutes

import (
	"github.com/antonybholmes/go-auth/userdbcache"
	"github.com/antonybholmes/go-edb-server/routes"
	"github.com/antonybholmes/go-edb-server/routes/authentication"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
)

type UserListReq struct {
	Offset  uint
	Records uint
}

type UserStatResp struct {
	Users uint `json:"users"`
}

func UserStatsRoute(c echo.Context) error {

	var req UserListReq

	c.Bind(&req)

	users, err := userdbcache.NumUsers()

	if err != nil {
		return routes.ErrorReq(err)
	}

	resp := UserStatResp{Users: users}

	return routes.MakeDataPrettyResp(c, "", resp)

}

func UsersRoute(c echo.Context) error {

	var req UserListReq

	c.Bind(&req)

	users, err := userdbcache.Users(req.Offset, req.Records)

	if err != nil {
		return routes.ErrorReq(err)
	}

	return routes.MakeDataPrettyResp(c, "", users)

}

func RolesRoute(c echo.Context) error {

	roles, err := userdbcache.Roles()

	if err != nil {
		return routes.ErrorReq(err)
	}

	return routes.MakeDataPrettyResp(c, "", roles)

}

func UpdateUserRoute(c echo.Context) error {

	return authentication.NewValidator(c).CheckUsernameIsWellFormed().CheckEmailIsWellFormed().LoadAuthUserFromPublicId().Success(func(validator *authentication.Validator) error {

		authUser := validator.AuthUser

		db, err := userdbcache.NewConn()

		if err != nil {
			return routes.ErrorReq(err)
		}

		defer db.Close()

		err = userdbcache.SetUserInfo(authUser.PublicId, validator.Req.Username, validator.Req.FirstName, validator.Req.LastName, db)

		if err != nil {
			return routes.ErrorReq(err)
		}

		err = userdbcache.SetEmailAddress(authUser.PublicId, validator.Address, db)

		if err != nil {
			return routes.ErrorReq(err)
		}

		if validator.Req.Password != "" {
			err = userdbcache.SetPassword(authUser.PublicId, validator.Req.Password, db)

			if err != nil {
				return routes.ErrorReq(err)
			}
		}

		// set roles

		err = userdbcache.SetUserRoles(authUser, validator.Req.Roles, db)

		if err != nil {
			return routes.ErrorReq(err)
		}

		return routes.MakeOkPrettyResp(c, "user updated")
	})
}

func AddUserRoute(c echo.Context) error {

	return authentication.NewValidator(c).CheckUsernameIsWellFormed().CheckEmailIsWellFormed().Success(func(validator *authentication.Validator) error {

		// assume email is not verified
		authUser, err := userdbcache.Instance().CreateUser(validator.Req.Username,
			validator.Address,
			validator.Req.Password,
			validator.Req.FirstName,
			validator.Req.LastName,
			validator.Req.EmailIsVerified)

		if err != nil {
			return routes.ErrorReq(err)
		}

		// tell user their account was created
		go SendAccountCreatedEmail(authUser, validator.Address)

		return routes.MakeOkPrettyResp(c, "account created email sent")
	})
}

func DeleteUserRoute(c echo.Context) error {
	publicId := c.Param("publicId")

	log.Debug().Msgf("delete %s", publicId)

	err := userdbcache.DeleteUser(publicId)

	if err != nil {
		return routes.ErrorReq(err)
	}

	return routes.MakeOkPrettyResp(c, "user deleted")
}
