package adminroutes

import (
	"github.com/antonybholmes/go-auth/userdbcache"
	"github.com/antonybholmes/go-edb-server/consts"
	"github.com/antonybholmes/go-edb-server/rdb"
	"github.com/antonybholmes/go-edb-server/routes"
	"github.com/antonybholmes/go-edb-server/routes/authentication"
	"github.com/antonybholmes/go-mailer"
	"github.com/labstack/echo/v4"
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

	users, err := userdbcache.Users(req.Records, req.Offset)

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

		//db, err := userdbcache.NewConn()

		// if err != nil {
		// 	return routes.ErrorReq(err)
		// }

		//defer db.Close()

		//authUser, err := userdbcache.FindUserByPublicId(validator.Req.PublicId)

		// if err != nil {
		// 	return routes.ErrorReq(err)
		// }

		authUser := validator.AuthUser

		err := userdbcache.SetUserInfo(authUser.PublicId, validator.Req.Username, validator.Req.FirstName, validator.Req.LastName)

		if err != nil {
			return routes.ErrorReq(err)
		}

		err = userdbcache.SetEmailAddress(authUser.PublicId, validator.Address)

		if err != nil {
			return routes.ErrorReq(err)
		}

		if validator.Req.Password != "" {
			err = userdbcache.SetPassword(authUser.PublicId, validator.Req.Password)

			if err != nil {
				return routes.ErrorReq(err)
			}
		}

		// set roles
		err = userdbcache.SetUserRoles(authUser, validator.Req.Roles)

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
		//go SendAccountCreatedEmail(authUser, validator.Address)

		email := mailer.RedisQueueEmail{Name: authUser.FirstName,
			To: authUser.Email,

			EmailType: mailer.REDIS_EMAIL_TYPE_ACCOUNT_CREATED,

			CallBackUrl: consts.APP_URL}
		rdb.PublishEmail(&email)

		return routes.MakeOkPrettyResp(c, "account created email sent")
	})
}

func DeleteUserRoute(c echo.Context) error {
	publicId := c.Param("publicId")

	err := userdbcache.DeleteUser(publicId)

	if err != nil {
		return routes.ErrorReq(err)
	}

	return routes.MakeOkPrettyResp(c, "user deleted")
}
