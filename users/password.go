package users

import (
	"github.com/antonybholmes/go-auth"
	"github.com/antonybholmes/go-edb-api/consts"
	"github.com/antonybholmes/go-edb-api/routes"
	"github.com/antonybholmes/go-edb-api/userdb"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

type PasswordResetReq struct {
	Password string `json:"password"`
}

// Start passwordless login by sending an email
func ResetPasswordEmailRoute(c echo.Context) error {
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

	otpJwt, err := auth.ResetPasswordToken(authUser.UserId, c.RealIP(), consts.JWT_SECRET)

	if err != nil {
		return routes.BadReq(err)
	}

	var file string

	if req.Url != "" {
		file = "templates/email/password/reset/web.html"
	} else {
		file = "templates/email/password/reset/api.html"
	}

	err = TokenEmail("Password Reset",
		authUser,
		file,
		otpJwt,
		req.CallbackUrl,
		req.Url)

	return routes.MakeSuccessResp(c, "password reset email sent", true)
}

func ResetPasswordRoute(c echo.Context) error {
	req := new(PasswordResetReq)

	err := c.Bind(req)

	if err != nil {
		return err
	}

	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*auth.JwtCustomClaims)

	if claims.Type != auth.TOKEN_TYPE_RESET_PASSWORD {
		return routes.BadReq("wrong token type")
	}

	authUser, err := userdb.FindUserById(claims.UserId)

	if err != nil {
		return routes.BadReq("user does not exist")
	}

	err = userdb.SetPassword(authUser.UserId, req.Password)

	if err != nil {
		return routes.BadReq("error setting password")
	}

	return routes.MakeSuccessResp(c, "password reset", true)
}
