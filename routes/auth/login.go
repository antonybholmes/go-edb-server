package auth

import (
	"github.com/antonybholmes/go-auth"
	"github.com/antonybholmes/go-auth/userdb"
	"github.com/antonybholmes/go-edb-api/consts"
	"github.com/antonybholmes/go-edb-api/routes"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

func LoginRoute(c echo.Context) error {
	req := new(auth.LoginReq)

	err := c.Bind(req)

	if err != nil {
		return err
	}

	//email := c.FormValue("email")
	//password := c.FormValue("password")

	//user := auth.LoginUserFromReq(req)

	authUser, err := userdb.FindUserByEmail(req.Email)

	if err != nil {
		return routes.BadReq("user does not exist")
	}

	if !authUser.IsVerified {
		return routes.BadReq("email address not verified")
	}

	if !authUser.CanAuth {
		return routes.BadReq("user not allowed tokens")
	}

	if req.Password == "" {
		return routes.BadReq("empty password: use passwordless")
	}

	if !authUser.CheckPasswords(req.Password) {
		return routes.BadReq("incorrect password")
	}

	t, err := auth.RefreshToken(authUser.Uuid, c.RealIP(), consts.JWT_SECRET)

	if err != nil {
		return routes.BadReq("error signing token")
	}

	return routes.MakeDataResp(c, "", &routes.JwtResp{Jwt: t})
}

// Start passwordless login by sending an email
func PasswordlessEmailRoute(c echo.Context) error {
	req := new(auth.EmailOnlyLoginReq)

	err := c.Bind(req)

	if err != nil {
		return err
	}

	authUser, err := userdb.FindUserByEmail(req.Email)

	if err != nil {
		return routes.BadReq("user does not exist")
	}

	if !authUser.IsVerified {
		return routes.BadReq("email address not verified")
	}

	otpJwt, err := auth.PasswordlessToken(authUser.Uuid, c.RealIP(), consts.JWT_SECRET)

	if err != nil {
		return routes.BadReq(err)
	}

	var file string

	if req.Url != "" {
		file = "templates/email/passwordless/web.html"
	} else {
		file = "templates/email/passwordless/api.html"
	}

	err = TokenEmail("Passwordless Login",
		authUser,
		file,
		otpJwt,
		req.CallbackUrl,
		req.Url)

	if err != nil {
		return routes.BadReq(err)
	}

	return routes.MakeSuccessResp(c, "passwordless email sent", true)
}

func PasswordlessLoginRoute(c echo.Context) error {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*auth.JwtCustomClaims)

	if claims.Type != auth.TOKEN_TYPE_PASSWORDLESS {
		return routes.BadReq("wrong token type")
	}

	authUser, err := userdb.FindUserByUuid(claims.Uuid)

	if err != nil {
		return routes.BadReq("user does not exist")
	}

	if !authUser.IsVerified {
		return routes.BadReq("email address not verified")
	}

	if !authUser.CanAuth {
		return routes.BadReq("user not allowed tokens")
	}

	t, err := auth.RefreshToken(authUser.Uuid, c.RealIP(), consts.JWT_SECRET)

	if err != nil {
		return routes.BadReq("error signing token")
	}

	return routes.MakeDataResp(c, "", &routes.JwtResp{Jwt: t})
}
