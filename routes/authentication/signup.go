package authentication

import (
	"fmt"

	"github.com/antonybholmes/go-auth/tokengen"
	"github.com/antonybholmes/go-auth/userdbcache"
	"github.com/antonybholmes/go-edb-server/consts"
	"github.com/antonybholmes/go-edb-server/rdb"
	"github.com/antonybholmes/go-edb-server/routes"
	"github.com/antonybholmes/go-mailer"
	"github.com/labstack/echo/v4"
)

func SignupRoute(c echo.Context) error {
	return NewValidator(c).CheckEmailIsWellFormed().Success(func(validator *Validator) error {
		req := validator.Req

		authUser, err := userdbcache.CreateUserFromSignup(req)

		if err != nil {
			return routes.ErrorReq(err)
		}

		otpToken, err := tokengen.VerifyEmailToken(c, authUser.PublicId, req.VisitUrl)

		//log.Debug().Msgf("%s", otpJwt)

		if err != nil {
			return routes.ErrorReq(err)
		}

		// var file string

		// if req.CallbackUrl != "" {
		// 	file = "templates/email/verify/web.html"
		// } else {
		// 	file = "templates/email/verify/api.html"
		// }

		// go SendEmailWithToken("Email Verification",
		// 	authUser,
		// 	file,
		// 	otpToken,
		// 	req.CallbackUrl,
		// 	req.VisitUrl)

		//if err != nil {
		//	return routes.ErrorReq(err)
		//}

		email := mailer.RedisQueueEmail{Name: authUser.FirstName,
			To:          authUser.Email,
			Token:       otpToken,
			EmailType:   mailer.REDIS_EMAIL_TYPE_VERIFY,
			Ttl:         fmt.Sprintf("%d minutes", int(consts.SHORT_TTL_MINS.Minutes())),
			CallBackUrl: req.CallbackUrl,
			//VisitUrl:    req.VisitUrl
		}
		rdb.PublishEmail(&email)

		return routes.MakeOkPrettyResp(c, "check your email for a verification link")
	})
}

func EmailAddressVerifiedRoute(c echo.Context) error {
	return NewValidator(c).LoadAuthUserFromToken().Success(func(validator *Validator) error {

		authUser := validator.AuthUser

		// if verified, stop and just return true
		if authUser.EmailIsVerified {
			return routes.MakeOkPrettyResp(c, "")
		}

		err := userdbcache.SetIsVerified(authUser.PublicId)

		if err != nil {
			return routes.MakeSuccessPrettyResp(c, "unable to verify user", false)
		}

		// file := "templates/email/verify/verified.html"

		// go SendEmailWithToken("Email Address Verified",
		// 	authUser,
		// 	file,
		// 	"",
		// 	"",
		// 	"")

		email := mailer.RedisQueueEmail{Name: authUser.FirstName,
			To:        authUser.Email,
			EmailType: mailer.REDIS_EMAIL_TYPE_VERIFIED}
		rdb.PublishEmail(&email)

		return routes.MakeOkPrettyResp(c, "email address verified")
	})
}
