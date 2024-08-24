package authroutes

import (
	"github.com/antonybholmes/go-auth"
	"github.com/antonybholmes/go-auth/userdbcache"
	"github.com/antonybholmes/go-edb-server/consts"
	"github.com/antonybholmes/go-edb-server/routes"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
)

func SignupRoute(c echo.Context) error {

	var req auth.SignupReq

	err := c.Bind(&req)

	if err != nil {
		return err
	}

	authUser, err := userdbcache.CreateStandardUser(&req)

	if err != nil {
		return routes.ErrorReq(err)
	}

	otpJwt, err := auth.VerifyEmailToken(c, authUser.Uuid, consts.JWT_PRIVATE_KEY)

	log.Debug().Msgf("%s", otpJwt)

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
	validator, err := routes.NewValidator(c).LoadAuthUserFromToken().Ok()

	if err != nil {
		return err
	}

	authUser := validator.AuthUser

	// if verified, stop and just return true
	if authUser.EmailVerified {
		return routes.MakeOkPrettyResp(c, "")
	}

	err = userdbcache.SetIsVerified(authUser.Uuid)

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

	return routes.MakeOkPrettyResp(c, "email address verified") //c.JSON(http.StatusOK, JWTResp{t})
}
