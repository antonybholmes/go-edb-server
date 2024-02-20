package users

import (
	"github.com/antonybholmes/go-auth"
	"github.com/antonybholmes/go-edb-api/consts"
	"github.com/antonybholmes/go-edb-api/routes"
	"github.com/antonybholmes/go-edb-api/userdb"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
)

func SignupRoute(c echo.Context, userdb *auth.UserDb, secret string) error {
	req := new(auth.SignupReq)

	err := c.Bind(req)

	if err != nil {
		return err
	}

	authUser, err := userdb.CreateUser(req)

	if err != nil {
		return routes.BadReq(err)
	}

	if authUser.IsVerified {
		return routes.BadReq("user is already verified")
	}

	otpJwt, err := auth.VerifyEmailToken(authUser.UserId, c.RealIP(), consts.JWT_SECRET)

	log.Debug().Msgf("%s", otpJwt)

	if err != nil {
		return routes.BadReq(err)
	}

	var file string

	if req.Url != "" {
		file = "templates/email/verify/web.html"
	} else {
		file = "templates/email/verify/api.html"
	}

	err = TokenEmail("Email Verification",
		authUser,
		file,
		otpJwt,
		req.CallbackUrl,
		req.Url)

	if err != nil {
		return routes.BadReq(err)
	}

	return routes.MakeSuccessResp(c, "verification email sent", true) //c.JSON(http.StatusOK, JWTResp{t})
}

func EmailVerificationRoute(c echo.Context) error {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*auth.JwtCustomClaims)

	authUser, err := userdb.FindUserById(claims.UserId)

	if err != nil {
		return routes.MakeSuccessResp(c, "user not found", false)
	}

	// if verified, stop and just return true
	if authUser.IsVerified {
		return routes.MakeSuccessResp(c, "", true)
	}

	err = userdb.SetIsVerified(authUser.UserId)

	if err != nil {
		return routes.MakeSuccessResp(c, "unable to verify user", false)
	}

	return routes.MakeSuccessResp(c, "", true) //c.JSON(http.StatusOK, JWTResp{t})
}
