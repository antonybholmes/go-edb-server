package routes

import (
	"bytes"
	"html/template"
	"net/url"
	"strings"

	"github.com/antonybholmes/go-auth"
	"github.com/antonybholmes/go-edb-api/consts"
	"github.com/antonybholmes/go-edb-api/users"
	"github.com/antonybholmes/go-mailer/email"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
)

type JwtResp struct {
	Jwt string `json:"jwt"`
}

type LoginResp struct {
	auth.PublicUser
	JwtResp
}

type JwtInfo struct {
	UserId string `json:"userId"`
	//Name  string `json:"name"`
	Type    string `json:"type"`
	IpAddr  string `json:"ipAddr"`
	Expires string `json:"expires"`
}

type ReqJwt struct {
	Jwt string `json:"jwt"`
}

func Signup(c echo.Context, userdb *auth.UserDb, secret string) error {
	req := new(auth.SignupReq)

	err := c.Bind(req)

	if err != nil {
		return err
	}

	randCode := auth.RandCode()

	authUser, err := userdb.CreateUser(req, randCode)

	if err != nil {
		return BadReq(err)
	}

	if authUser.IsVerified {
		return BadReq("user is already verified")
	}

	otpJwt, err := auth.CreateOtpJwt(authUser.UserId, randCode, c.RealIP(), consts.JWT_SECRET)

	log.Debug().Msgf("%s", otpJwt)

	if err != nil {
		return BadReq(err)
	}

	var body bytes.Buffer

	var file string

	if req.Url != "" {
		file = "templates/email/verify/web.html"
	} else {
		file = "templates/email/verify/api.html"
	}

	t, err := template.ParseFiles(file)

	if err != nil {
		return BadReq(err)
	}

	firstName := strings.Split(authUser.Name, " ")[0]

	if req.CallbackUrl != "" {
		callbackUrl, err := url.Parse(req.CallbackUrl)

		if err != nil {
			return BadReq(err)
		}

		params, err := url.ParseQuery(callbackUrl.RawQuery)

		if err != nil {
			return BadReq(err)
		}

		if req.Url != "" {
			params.Set("url", req.Url)
		}

		params.Set("otp", otpJwt)

		callbackUrl.RawQuery = params.Encode()

		link := callbackUrl.String()

		err = t.Execute(&body, struct {
			Name string
			Link string
			From string
		}{
			Name: firstName,
			Link: link,
			From: consts.NAME,
		})

		if err != nil {
			return BadReq(err)
		}
	} else {
		err = t.Execute(&body, struct {
			Name string
			Code string
			From string
		}{
			Name: firstName,
			Code: otpJwt,
			From: consts.NAME,
		})

		if err != nil {
			return BadReq(err)
		}
	}

	//log.Debug().Msgf("%s", body.String())

	err = email.SendHtmlEmail(req.Mailbox(), "Email verification", body.String())

	if err != nil {
		return BadReq(err)
	}

	// if err != nil {
	// 	log.Error().Msgf("%s", err)
	// }

	// // Set custom claims
	// claims := &JwtCustomClaims{
	// 	UserId: authUser.UserId,
	// 	//Email: authUser.Email,
	// 	IpAddr: c.RealIP(),
	// 	RegisteredClaims: jwt.RegisteredClaims{
	// 		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * auth.JWT_TOKEN_EXPIRES_HOURS)),
	// 	},
	// }

	// // Create token with claims
	// token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// // Generate encoded token and send it as response.
	// t, err := token.SignedString([]byte(secret))

	// if err != nil {
	// 	return echo.NewHTTPError(http.StatusBadRequest, err)
	// }

	return MakeSuccessResp(c, "verification email sent", true) //c.JSON(http.StatusOK, JWTResp{t})
}

func EmailVerificationRoute(c echo.Context) error {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*auth.JwtOtpCustomClaims)

	authUser, err := users.FindUserById(claims.UserId)

	if err != nil {
		return MakeSuccessResp(c, "user not found", false)
	}

	// if verified, stop and just return true
	if authUser.IsVerified {
		return MakeSuccessResp(c, "", true)
	}

	if authUser.OTP != claims.OTP {
		return MakeSuccessResp(c, "error: wrong otp code", false)
	}

	err = users.SetIsVerified(authUser.UserId)

	if err != nil {
		return MakeSuccessResp(c, "unable to verify user", false)
	}

	return MakeSuccessResp(c, "", true) //c.JSON(http.StatusOK, JWTResp{t})
}

func LoginRoute(c echo.Context) error {
	req := new(auth.LoginReq)

	err := c.Bind(req)

	if err != nil {
		return err
	}

	//email := c.FormValue("email")
	//password := c.FormValue("password")

	//user := auth.LoginUserFromReq(req)

	authUser, err := users.FindUserByEmail(req.Email)

	if err != nil {
		return BadReq("user does not exist")
	}

	if !authUser.IsVerified {
		return BadReq("email address not verified")
	}

	if !authUser.CanAuth {
		return BadReq("user not allowed tokens")
	}

	if req.Password == "" {
		return BadReq("empty password: use passwordless")
	}

	if !authUser.CheckPasswords(req.Password) {
		return BadReq("incorrect password")
	}

	t, err := auth.CreateRefreshToken(authUser.UserId, c.RealIP(), consts.JWT_SECRET)

	if err != nil {
		return BadReq("error signing token")
	}

	return MakeDataResp(c, "", &LoginResp{
		JwtResp:    JwtResp{Jwt: t},
		PublicUser: auth.PublicUser{UserId: authUser.UserId, User: auth.User{Name: authUser.Name, Email: authUser.Email}}})
}

// Start passwordless login by sending an email
func SendPasswordlessEmail(c echo.Context) error {
	req := new(auth.PasswordlessLoginReq)

	err := c.Bind(req)

	if err != nil {
		return err
	}

	authUser, err := users.FindUserByEmail(req.Email)

	if err != nil {
		return BadReq("user does not exist")
	}

	if !authUser.IsVerified {
		return BadReq("email address not verified")
	}

	randCode := auth.RandCode()

	otpJwt, err := auth.CreateOtpJwt(authUser.UserId, randCode, c.RealIP(), consts.JWT_SECRET)

	if err != nil {
		return BadReq(err)
	}

	var body bytes.Buffer

	var file string

	if req.Url != "" {
		file = "templates/email/passwordless/web.html"
	} else {
		file = "templates/email/passwordless/api.html"
	}

	t, err := template.ParseFiles(file)

	if err != nil {
		return BadReq(err)
	}

	firstName := strings.Split(authUser.Name, " ")[0]

	if req.CallbackUrl != "" {
		callbackUrl, err := url.Parse(req.CallbackUrl)

		if err != nil {
			return BadReq(err)
		}

		params, err := url.ParseQuery(callbackUrl.RawQuery)

		if err != nil {
			return BadReq(err)
		}

		if req.Url != "" {
			params.Set("url", req.Url)
		}

		params.Set("otp", otpJwt)

		callbackUrl.RawQuery = params.Encode()

		link := callbackUrl.String()

		err = t.Execute(&body, struct {
			Name string
			Link string
			From string
		}{
			Name: firstName,
			Link: link,
			From: consts.NAME,
		})

		if err != nil {
			return BadReq(err)
		}
	} else {
		err = t.Execute(&body, struct {
			Name string
			Code string
			From string
		}{
			Name: firstName,
			Code: otpJwt,
			From: consts.NAME,
		})

		if err != nil {
			return BadReq(err)
		}
	}

	//log.Debug().Msgf("%s", body.String())

	users.SetOtp(authUser.UserId, randCode)

	if err != nil {
		return BadReq(err)
	}

	err = email.SendHtmlEmail(req.Mailbox(), "Passwordless Login", body.String())

	if err != nil {
		return BadReq(err)
	}

	return MakeSuccessResp(c, "passwordless email sent", true)
}

func PasswordlessLoginRoute(c echo.Context) error {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*auth.JwtOtpCustomClaims)

	authUser, err := users.FindUserById(claims.UserId)

	if err != nil {
		return BadReq("user does not exist")
	}

	if !authUser.IsVerified {
		return BadReq("email address not verified")
	}

	if !authUser.CanAuth {
		return BadReq("user not allowed tokens")
	}

	t, err := auth.CreateRefreshToken(authUser.UserId, c.RealIP(), consts.JWT_SECRET)

	if err != nil {
		return BadReq("error signing token")
	}

	return MakeDataResp(c, "", &LoginResp{
		JwtResp:    JwtResp{Jwt: t},
		PublicUser: auth.PublicUser{UserId: authUser.UserId, User: auth.User{Name: authUser.Name, Email: authUser.Email}}})
}
