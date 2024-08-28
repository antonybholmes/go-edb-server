package authroutes

import (
	"bytes"
	"fmt"
	"html/template"

	"net/mail"
	"net/url"
	"strings"

	"github.com/antonybholmes/go-auth"
	"github.com/antonybholmes/go-edb-server/consts"
	"github.com/antonybholmes/go-edb-server/routes"
	"github.com/antonybholmes/go-mailer/mailer"
)

const JWT_PARAM = "jwt"
const URL_PARAM = "url"

type EmailBody struct {
	Name       string
	From       string
	Time       string
	Link       string
	DoNotReply string
}

func SendEmailWithToken(subject string,
	authUser *auth.AuthUser,
	file string,
	jwt string,
	callbackUrl string,
	vistUrl string) error {

	address, err := mail.ParseAddress(authUser.Email)

	if err != nil {
		return routes.ErrorReq(err)
	}

	return BaseSendEmailWithToken(subject, authUser, address, file, jwt, callbackUrl, vistUrl)
}

// Generic method for sending an email with a token in it. For APIS this is a token to use in the request, for websites
// it can craft a callback url with the token added as a parameter so that the web app can deal with the response.
func BaseSendEmailWithToken(subject string,
	authUser *auth.AuthUser,
	address *mail.Address,
	file string,
	jwt string,
	callbackUrl string,
	vistUrl string) error {

	var body bytes.Buffer

	t, err := template.ParseFiles(file)

	if err != nil {
		return routes.ErrorReq(err)
	}

	var firstName string = ""

	if len(authUser.FirstName) > 0 {
		firstName = authUser.FirstName
	} else {
		firstName = strings.Split(address.Address, "@")[0]
	}

	firstName = strings.Split(firstName, " ")[0]

	time := fmt.Sprintf("%d minutes", int(auth.TOKEN_TTL_10_MINS.Minutes()))

	if callbackUrl != "" {
		callbackUrl, err := url.Parse(callbackUrl)

		if err != nil {
			return routes.ErrorReq(err)
		}

		params, err := url.ParseQuery(callbackUrl.RawQuery)

		if err != nil {
			return routes.ErrorReq(err)
		}

		if vistUrl != "" {
			params.Set(URL_PARAM, vistUrl)
		}

		params.Set(JWT_PARAM, jwt)

		callbackUrl.RawQuery = params.Encode()

		link := callbackUrl.String()

		err = t.Execute(&body, EmailBody{
			Name:       firstName,
			Link:       link,
			From:       consts.NAME,
			Time:       time,
			DoNotReply: consts.DO_NOT_REPLY,
		})

		if err != nil {
			return routes.ErrorReq(err)
		}
	} else {
		err = t.Execute(&body, EmailBody{
			Name:       firstName,
			Link:       jwt,
			From:       consts.NAME,
			Time:       time,
			DoNotReply: consts.DO_NOT_REPLY,
		})

		if err != nil {
			return routes.ErrorReq(err)
		}
	}

	err = mailer.SendHtmlEmail(address, subject, body.String())

	if err != nil {
		return err
	}

	return nil
}
