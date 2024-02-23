package auth

import (
	"bytes"
	"fmt"
	"html/template"
	"net/mail"
	"net/url"
	"strings"

	"github.com/antonybholmes/go-auth"
	"github.com/antonybholmes/go-auth/userdb"
	"github.com/antonybholmes/go-edb-api/consts"
	"github.com/antonybholmes/go-edb-api/routes"

	"github.com/antonybholmes/go-mailer/email"
	"github.com/labstack/echo/v4"
)

// type LoginResp struct {
// 	auth.PublicUser
// 	JwtResp
// }

type EmailBody struct {
	Name string
	From string
	Time string
	Link string
}

type PasswordResetReq struct {
	Password string `json:"password"`
}

type UsernameReq struct {
	Username string `json:"username"`
}

// Start passwordless login by sending an email
func ResetPasswordEmailRoute(c echo.Context) error {
	req := new(auth.EmailOnlyLoginReq)

	err := c.Bind(req)

	if err != nil {
		return err
	}

	return routes.ValidEmailCB(c, req.Email, func(c echo.Context, email *mail.Address) error {
		return routes.EmailUserCB(c, email, func(c echo.Context, authUser *auth.AuthUser) error {
			return routes.VerifiedEmailCB(c, authUser, func(c echo.Context, authUser *auth.AuthUser) error {

				otpJwt, err := auth.ResetPasswordToken(authUser.Uuid, c.RealIP(), consts.JWT_SECRET)

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

				if err != nil {
					return routes.BadReq(err)
				}

				return routes.MakeSuccessResp(c, "password reset email sent", true)
			})
		})
	})
}

func UpdatePasswordRoute(c echo.Context) error {
	req := new(PasswordResetReq)

	err := c.Bind(req)

	if err != nil {
		return routes.BadReq(err)
	}

	return routes.JwtCB(c, func(c echo.Context, claims *auth.JwtCustomClaims) error {

		if claims.Type != auth.TOKEN_TYPE_RESET_PASSWORD {
			return routes.BadReq("wrong token type")
		}

		return routes.UuidUserCB(c, claims, func(c echo.Context, claims *auth.JwtCustomClaims, authUser *auth.AuthUser) error {
			err = userdb.SetPassword(authUser.Uuid, req.Password)

			if err != nil {
				return routes.BadReq("error setting password")
			}

			return routes.MakeSuccessResp(c, "password updated", true)
		})
	})
}

func UpdateUsernameRoute(c echo.Context) error {
	req := new(UsernameReq)

	err := c.Bind(req)

	if err != nil {
		return routes.BadReq(err)
	}

	return routes.JwtCB(c, func(c echo.Context, claims *auth.JwtCustomClaims) error {
		return routes.UuidUserCB(c, claims, func(c echo.Context, claims *auth.JwtCustomClaims, authUser *auth.AuthUser) error {

			err = userdb.SetUsername(authUser.Uuid, req.Username)

			if err != nil {
				return routes.BadReq("error setting password")
			}

			return routes.MakeSuccessResp(c, "password updated", true)
		})
	})
}

func UserInfoRoute(c echo.Context) error {
	return routes.RefreshTokenCB(c, func(c echo.Context, claims *auth.JwtCustomClaims) error {
		return routes.UuidUserCB(c, claims, func(c echo.Context, claims *auth.JwtCustomClaims, authUser *auth.AuthUser) error {
			return routes.MakeDataResp(c, "", *authUser.ToPublicUser())
		})
	})
}

// Generic method for sending an email with a token in it. For APIS this is a token to use in the request, for websites
// it can craft a callback url with the token added as a parameter so that the web app can deal with the response.
func TokenEmail(subject string,
	authUser *auth.AuthUser,
	file string,
	token string,
	callbackUrl string,
	vistUrl string) error {

	var body bytes.Buffer

	t, err := template.ParseFiles(file)

	if err != nil {
		return routes.BadReq(err)
	}

	firstName := strings.Split(authUser.Name, " ")[0]

	time := fmt.Sprintf("%d minutes", int(auth.TOKEN_TYPE_OTP_TTL_MINS.Minutes()))

	if callbackUrl != "" {
		callbackUrl, err := url.Parse(callbackUrl)

		if err != nil {
			return routes.BadReq(err)
		}

		params, err := url.ParseQuery(callbackUrl.RawQuery)

		if err != nil {
			return routes.BadReq(err)
		}

		if vistUrl != "" {
			params.Set("url", vistUrl)
		}

		params.Set("otp", token)

		callbackUrl.RawQuery = params.Encode()

		link := callbackUrl.String()

		err = t.Execute(&body, EmailBody{
			Name: firstName,
			Link: link,
			From: consts.NAME,
			Time: time,
		})

		if err != nil {
			return routes.BadReq(err)
		}
	} else {
		err = t.Execute(&body, EmailBody{
			Name: firstName,
			Link: token,
			From: consts.NAME,
			Time: time,
		})

		if err != nil {
			return routes.BadReq(err)
		}
	}

	err = email.SendHtmlEmail(authUser.Address(), subject, body.String())

	if err != nil {
		return err
	}

	return nil
}