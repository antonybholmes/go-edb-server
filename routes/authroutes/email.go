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
	"github.com/rs/zerolog/log"
)

const DO_NOT_REPLY = "Please do not reply to this message. It was sent from a notification-only email address that we don't monitor."
const TOKEN_PARAM = "token"
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
	token string,
	callbackUrl string,
	vistUrl string) error {

	log.Debug().Msgf("SendEmailWithToken %v", authUser.Email)

	address, err := mail.ParseAddress(authUser.Email)

	if err != nil {
		log.Debug().Msgf("asdasd %v %s", authUser.Email, err)
		return routes.ErrorReq(err)
	}

	log.Debug().Msgf("asdasd %v %s", authUser, address)

	return BaseSendEmailWithToken(subject, authUser, address, file, token, callbackUrl, vistUrl)
}

// Generic method for sending an email with a token in it. For APIS this is a token to use in the request, for websites
// it can craft a callback url with the token added as a parameter so that the web app can deal with the response.
func BaseSendEmailWithToken(subject string,
	authUser *auth.AuthUser,
	address *mail.Address,
	file string,
	token string,
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

	time := fmt.Sprintf("%d minutes", int(auth.TOKEN_TYPE_SHORT_TIME_TTL_MINS.Minutes()))

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

		params.Set(TOKEN_PARAM, token)

		callbackUrl.RawQuery = params.Encode()

		link := callbackUrl.String()

		err = t.Execute(&body, EmailBody{
			Name:       firstName,
			Link:       link,
			From:       consts.NAME,
			Time:       time,
			DoNotReply: DO_NOT_REPLY,
		})

		if err != nil {
			return routes.ErrorReq(err)
		}
	} else {
		err = t.Execute(&body, EmailBody{
			Name:       firstName,
			Link:       token,
			From:       consts.NAME,
			Time:       time,
			DoNotReply: DO_NOT_REPLY,
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
