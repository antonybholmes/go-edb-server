package auth

import (
	"net/mail"

	"github.com/antonybholmes/go-auth"
	"github.com/antonybholmes/go-auth/userdb"
	"github.com/antonybholmes/go-edb-api/consts"
	"github.com/antonybholmes/go-edb-api/routes"
	"github.com/labstack/echo/v4"
)

func EmailPasswordLoginRoute(c echo.Context) error {
	req := new(auth.EmailPasswordLoginReq)

	err := c.Bind(req)

	if err != nil {
		return err
	}

	if req.Password == "" {
		return routes.BadReq("empty password: use passwordless")
	}

	return routes.AuthUserFromEmailCB(c, req.Email, func(c echo.Context, authUser *auth.AuthUser) error {
		return loginRoute(c, authUser, req.Password)
	})
}

func UsernamePasswordLoginRoute(c echo.Context) error {
	req := new(auth.UsernamePasswordLoginReq)

	err := c.Bind(req)

	if err != nil {
		return err
	}

	if req.Password == "" {
		return routes.BadReq("empty password: use passwordless")
	}

	authUser, err := userdb.FindUserByUsername(req.Username)

	if err != nil {
		email, err := mail.ParseAddress(req.Username)

		if err != nil {
			return routes.InvalidEmailReq()
		}

		// also check if username is valid email and try to login
		// with that
		authUser, err = userdb.FindUserByEmail(email)

		if err != nil {
			return routes.BadReq("user does not exist")
		}
	}

	return loginRoute(c, authUser, req.Password)
}

func loginRoute(c echo.Context, authUser *auth.AuthUser, password string) error {

	if !authUser.EmailVerified {
		return routes.BadReq("email address not verified")
	}

	if !authUser.CanAuth {
		return routes.BadReq("user not allowed tokens")
	}

	if !authUser.CheckPasswords(password) {
		return routes.BadReq("incorrect password")
	}

	t, err := auth.RefreshToken(c, authUser.Uuid, consts.JWT_SECRET)

	if err != nil {
		return routes.BadReq("error signing token")
	}

	return routes.MakeDataResp(c, "", &routes.RefreshTokenResp{RefreshToken: t})
}

// Start passwordless login by sending an email
func PasswordlessEmailRoute(c echo.Context) error {
	req := new(auth.EmailOnlyLoginReq)

	err := c.Bind(req)

	if err != nil {
		return err
	}

	email, err := mail.ParseAddress(req.Email)

	if err != nil {
		return routes.InvalidEmailReq()
	}

	authUser, err := userdb.FindUserByEmail(email)

	if err != nil {
		return routes.BadReq("user does not exist")
	}

	if !authUser.EmailVerified {
		return routes.BadReq("email address not verified")
	}

	otpJwt, err := auth.PasswordlessToken(c, authUser.Uuid, consts.JWT_SECRET)

	if err != nil {
		return routes.BadReq(err)
	}

	var file string

	if req.Url != "" {
		file = "templates/email/passwordless/web.html"
	} else {
		file = "templates/email/passwordless/api.html"
	}

	err = SendEmailWithToken("Passwordless Login",
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

	return routes.UserFromUuidCB(c, nil, func(c echo.Context, claims *auth.JwtCustomClaims, authUser *auth.AuthUser) error {
		return routes.VerifiedEmailCB(c, authUser, func(c echo.Context, authUser *auth.AuthUser) error {
			if claims.Type != auth.TOKEN_TYPE_PASSWORDLESS {
				return routes.BadReq("wrong token type")
			}

			if !authUser.CanAuth {
				return routes.BadReq("user not allowed tokens")
			}

			t, err := auth.RefreshToken(c, authUser.Uuid, consts.JWT_SECRET)

			if err != nil {
				return routes.BadReq("error signing token")
			}

			return routes.MakeDataResp(c, "", &routes.RefreshTokenResp{RefreshToken: t})
		})
	})

}
