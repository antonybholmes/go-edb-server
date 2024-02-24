package auth

import (
	"net/mail"

	"github.com/antonybholmes/go-auth"
	"github.com/antonybholmes/go-auth/userdb"
	"github.com/antonybholmes/go-edb-api/consts"
	"github.com/antonybholmes/go-edb-api/routes"
	"github.com/labstack/echo/v4"
)

func UsernamePasswordLoginRoute(c echo.Context) error {
	return routes.ReqBindCB(c, new(auth.UsernamePasswordLoginReq), func(c echo.Context, req *auth.UsernamePasswordLoginReq) error {

		if req.Password == "" {
			return routes.BadReq("empty password: use passwordless")
		}

		authUser, err := userdb.FindUserByUsername(req.Username)

		if err != nil {
			email, err := mail.ParseAddress(req.Username)

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

		if !authUser.CanAuth {
			return routes.BadReq("user not allowed tokens")
		}

		if !authUser.CheckPasswords(req.Password) {
			return routes.BadReq("incorrect password")
		}

		refreshToken, err := auth.RefreshToken(c, authUser.Uuid, consts.JWT_SECRET)

		if err != nil {
			return routes.BadReq("error signing token")
		}

		accessToken, err := auth.AccessToken(c, authUser.Uuid, consts.JWT_SECRET)

		if err != nil {
			return routes.BadReq("error signing token")
		}

		return routes.MakeDataResp(c, "", &routes.LoginResp{RefreshToken: refreshToken, AccessToken: accessToken})
	})
}

// Start passwordless login by sending an email
func PasswordlessEmailRoute(c echo.Context) error {
	return routes.ReqBindCB(c, new(auth.EmailOnlyLoginReq), func(c echo.Context, req *auth.EmailOnlyLoginReq) error {
		return routes.AuthUserFromEmailCB(c, req.Email, func(c echo.Context, authUser *auth.AuthUser) error {
			return routes.VerifiedEmailCB(c, authUser, func(c echo.Context, authUser *auth.AuthUser) error {

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
			})
		})
	})
}

func PasswordlessLoginRoute(c echo.Context) error {
	return routes.AuthUserFromUuidCB(c, nil, func(c echo.Context, claims *auth.JwtCustomClaims, authUser *auth.AuthUser) error {
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
