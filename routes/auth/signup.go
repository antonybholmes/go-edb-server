package auth

import (
	"github.com/antonybholmes/go-auth"
	"github.com/antonybholmes/go-auth/userdb"
	"github.com/antonybholmes/go-edb-api/consts"
	"github.com/antonybholmes/go-edb-api/routes"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
)

func SignupRoute(c echo.Context) error {
	req := new(auth.SignupReq)

	err := c.Bind(req)

	if err != nil {
		return err
	}

	authUser, err := userdb.CreateUser(req)

	if err != nil {
		return routes.ErrorReq(err)
	}

	otpJwt, err := auth.VerifyEmailToken(c, authUser.Uuid, consts.JWT_SECRET)

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

	err = SendEmailWithToken("Email Verification",
		authUser,
		file,
		otpJwt,
		req.CallbackUrl,
		req.Url)

	if err != nil {
		return routes.ErrorReq(err)
	}

	return routes.MakeOkResp(c, "verification email sent") //c.JSON(http.StatusOK, JWTResp{t})
}

func EmailVerificationRoute(c echo.Context) error {
	validator, err := routes.NewValidator(c).AuthUserFromUuid().Ok()

	if err != nil {
		return err
	}

	authUser := validator.AuthUser

	// if verified, stop and just return true
	if authUser.EmailVerified {
		return routes.MakeOkResp(c, "")
	}

	err = userdb.SetIsVerified(authUser.Uuid)

	if err != nil {
		return routes.MakeSuccessResp(c, "unable to verify user", false)
	}

	file := "templates/email/verify/verified.html"

	err = SendEmailWithToken("Email Address Verified",
		authUser,
		file,
		"",
		"",
		"")

	if err != nil {
		return routes.ErrorReq(err)
	}

	return routes.MakeOkResp(c, "") //c.JSON(http.StatusOK, JWTResp{t})

	// return routes.AuthUserFromUuidCB(c, nil, func(c echo.Context, claims *auth.JwtCustomClaims, authUser *auth.AuthUser) error {

	// 	// if verified, stop and just return true
	// 	if authUser.EmailVerified {
	// 		return routes.MakeSuccessResp(c, "", true)
	// 	}

	// 	err := userdb.SetIsVerified(authUser.Uuid)

	// 	if err != nil {
	// 		return routes.MakeSuccessResp(c, "unable to verify user", false)
	// 	}

	// 	file := "templates/email/verify/verified.html"

	// 	err = SendEmailWithToken("Email Address Verified",
	// 		authUser,
	// 		file,
	// 		"",
	// 		"",
	// 		"")

	// 	if err != nil {
	// 		return routes.ErrorReq(err)
	// 	}

	// 	return routes.MakeSuccessResp(c, "", true) //c.JSON(http.StatusOK, JWTResp{t})
	// })

}
