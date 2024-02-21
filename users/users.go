package users

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
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

type UserInfoResp struct {
	auth.PublicUser
}

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

func UserInfoRoute(c echo.Context) error {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*auth.JwtCustomClaims)

	if claims.Type < auth.TOKEN_TYPE_REFRESH {
		return routes.BadReq("wrong token type")
	}

	authUser, err := userdb.FindUserById(claims.UserId)

	if err != nil {
		return routes.BadReq("user does not exist")
	}

	return routes.MakeDataResp(c, "", &UserInfoResp{
		PublicUser: auth.PublicUser{UserId: authUser.UserId, User: auth.User{Name: authUser.Name, Email: authUser.Email}}})
}

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

	err = email.SendHtmlEmail(authUser.Mailbox(), subject, body.String())

	if err != nil {
		return err
	}

	return nil
}
