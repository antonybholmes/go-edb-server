package auth

import (
	"bytes"
	"fmt"
	"html/template"
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

type NameReq struct {
	Name string `json:"name"`
}

// Start passwordless login by sending an email
func ResetPasswordEmailRoute(c echo.Context) error {

	return routes.ReqBindCB(c, new(auth.EmailOnlyLoginReq), func(c echo.Context, req *auth.EmailOnlyLoginReq) error {
		return routes.AuthUserFromEmailCB(c, req.Email, func(c echo.Context, authUser *auth.AuthUser) error {
			return routes.VerifiedEmailCB(c, authUser, func(c echo.Context, authUser *auth.AuthUser) error {

				otpJwt, err := auth.ResetPasswordToken(c, authUser.Uuid, consts.JWT_SECRET)

				if err != nil {
					return routes.BadReq(err)
				}

				var file string

				if req.Url != "" {
					file = "templates/email/password/reset/web.html"
				} else {
					file = "templates/email/password/reset/api.html"
				}

				err = SendEmailWithToken("Password Reset",
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

	return routes.ReqBindCB(c, new(PasswordResetReq), func(c echo.Context, req *PasswordResetReq) error {
		return routes.AuthUserFromUuidCB(c, nil, func(c echo.Context, claims *auth.JwtCustomClaims, authUser *auth.AuthUser) error {
			if claims.Type != auth.TOKEN_TYPE_RESET_PASSWORD {
				return routes.BadReq("wrong token type")
			}

			err := userdb.SetPassword(authUser.Uuid, req.Password)

			if err != nil {
				return routes.BadReq("error setting password")
			}

			return routes.MakeSuccessResp(c, "password updated", true)
		})
	})
}

func UpdateUsernameRoute(c echo.Context) error {

	return routes.ReqBindCB(c, new(UsernameReq), func(c echo.Context, req *UsernameReq) error {
		return routes.IsValidAccessTokenCB(c, func(c echo.Context, claims *auth.JwtCustomClaims) error {
			return routes.AuthUserFromUuidCB(c, claims, func(c echo.Context, claims *auth.JwtCustomClaims, authUser *auth.AuthUser) error {
				err := userdb.SetUsername(authUser.Uuid, req.Username)

				if err != nil {
					return routes.BadReq("error setting password")
				}

				return routes.MakeSuccessResp(c, "password updated", true)
			})
		})
	})

}

func UpdateNameRoute(c echo.Context) error {
	return routes.ReqBindCB(c, new(NameReq), func(c echo.Context, req *NameReq) error {
		return routes.IsValidAccessTokenCB(c, func(c echo.Context, claims *auth.JwtCustomClaims) error {
			return routes.AuthUserFromUuidCB(c, claims, func(c echo.Context, claims *auth.JwtCustomClaims, authUser *auth.AuthUser) error {

				err := userdb.SetName(authUser.Uuid, req.Name)

				if err != nil {
					return routes.BadReq("error setting password")
				}

				return routes.MakeSuccessResp(c, "name updated", true)
			})
		})
	})
}

func UserInfoRoute(c echo.Context) error {
	return routes.IsValidAccessTokenCB(c, func(c echo.Context, claims *auth.JwtCustomClaims) error {
		return routes.AuthUserFromUuidCB(c, claims, func(c echo.Context, claims *auth.JwtCustomClaims, authUser *auth.AuthUser) error {
			return routes.MakeDataResp(c, "", *authUser.ToPublicUser())
		})
	})
}

// Generic method for sending an email with a token in it. For APIS this is a token to use in the request, for websites
// it can craft a callback url with the token added as a parameter so that the web app can deal with the response.
func SendEmailWithToken(subject string,
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

	var firstName string = ""

	if len(authUser.Name) > 0 {
		firstName = authUser.Name
	} else {
		firstName = authUser.Email.Address
	}

	firstName = strings.Split(firstName, " ")[0]

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

	err = email.SendHtmlEmail(authUser.Email, subject, body.String())

	if err != nil {
		return err
	}

	return nil
}
