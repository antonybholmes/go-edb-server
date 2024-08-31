package authroutes

import (
	"github.com/antonybholmes/go-auth"
	"github.com/antonybholmes/go-auth/jwtgen"
	"github.com/antonybholmes/go-auth/userdbcache"
	"github.com/antonybholmes/go-edb-server/routes"
	"github.com/labstack/echo/v4"
)

func SignupRoute(c echo.Context) error {

	var req auth.SignupReq

	err := c.Bind(&req)

	if err != nil {
		return err
	}

	authUser, err := userdbcache.CreateUserFromSignup(&req)

	if err != nil {
		return routes.ErrorReq(err)
	}

	otpJwt, err := jwtgen.VerifyEmailJwt(c, authUser.PublicId)

	//log.Debug().Msgf("%s", otpJwt)

	if err != nil {
		return routes.ErrorReq(err)
	}

	var file string

	if req.CallbackUrl != "" {
		file = "templates/email/verify/web.html"
	} else {
		file = "templates/email/verify/api.html"
	}

	go SendEmailWithToken("Email Verification",
		authUser,
		file,
		otpJwt,
		req.CallbackUrl,
		req.Url)

	//if err != nil {
	//	return routes.ErrorReq(err)
	//}

	return routes.MakeOkPrettyResp(c, "check your email for a verification link") //c.JSON(http.StatusOK, JWTResp{t})
}

func EmailAddressWasVerifiedRoute(c echo.Context) error {
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

		file := "templates/email/verify/verified.html"

		go SendEmailWithToken("Email Address Verified",
			authUser,
			file,
			"",
			"",
			"")

		//if err != nil {
		//	return routes.ErrorReq(err)
		//}

		return routes.MakeOkPrettyResp(c, "email address verified")
	})
}
