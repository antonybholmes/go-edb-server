package routes

import (
	"bytes"
	"html/template"
	"net/url"
	"strings"
	"time"

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

type JwtValidResp struct {
	JwtIsValid bool `json:"jwtIsValid"`
}

type JwtInfo struct {
	UserId string `json:"userId"`
	//Name  string `json:"name"`
	//Email   string `json:"email"`
	IpAddr  string `json:"ipAddr"`
	Expires string `json:"expires"`
	Expired bool   `json:"expired"`
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

	//email := c.FormValue("email")
	//password := c.FormValue("password")

	loginUser := auth.NewSignupUser(req.Name, req.Email, req.Password)

	log.Debug().Msgf("%s", loginUser)

	randCode := auth.RandCode()

	authUser, err := userdb.CreateUser(loginUser, randCode)

	if err != nil {
		return BadReq(err)
	}

	if authUser.IsVerified {
		return BadReq("user is already verified")
	}

	otpJwt, err := auth.CreateOtpJwt(authUser, randCode, c.RealIP(), consts.JWT_SECRET)

	log.Debug().Msgf("%s", otpJwt)

	if err != nil {
		return BadReq(err)
	}

	var body bytes.Buffer

	var file string

	if req.Url != "" {
		file = "templates/verification/web.html"
	} else {
		file = "templates/verification/api.html"
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

	log.Debug().Msgf("%s", body.String())

	err = email.SendHtmlEmail(loginUser.Mailbox(), "Email verification", body.String())

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

	return MakeDataResp(c, "verification email sent", &[]string{}) //c.JSON(http.StatusOK, JWTResp{t})
}

func Verification(c echo.Context) error {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*auth.JwtOtpCustomClaims)

	log.Debug().Msgf("%s", claims.UserId)

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

	user := auth.LoginUserFromReq(req)

	authUser, err := users.FindUserByEmail(user.Email)

	if err != nil {
		return BadReq("user does not exist")
	}

	if !authUser.CheckPasswords(user.Password) {
		return BadReq("incorrect password")
	}

	// Throws unauthorized error
	//if username != "edb" || password != "tod4EwVHEyCRK8encuLE" {
	//	return echo.ErrUnauthorized
	//}

	// Set custom claims
	claims := &auth.JwtCustomClaims{
		UserId: authUser.UserId,
		//Email: authUser.Email,
		IpAddr: c.RealIP(),
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * auth.JWT_TOKEN_EXPIRES_HOURS)),
		},
	}

	// Create token with claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Generate encoded token and send it as response.
	t, err := token.SignedString([]byte(consts.JWT_SECRET))

	if err != nil {
		return BadReq("error signing token")
	}

	return MakeDataResp(c, "", &LoginResp{JwtResp: JwtResp{Jwt: t},
		PublicUser: auth.PublicUser{UserId: authUser.UserId,
			User: auth.User{Name: authUser.Name, Email: authUser.Email}}})
}

func ValidateToken(c echo.Context) error {
	// jwtReq := new(ReqJwt)

	// err := c.Bind(jwtReq)

	// if err != nil {
	// 	return err
	// }

	// token, err := jwt.ParseWithClaims(jwtReq.Jwt, &JwtCustomClaims{}, func(token *jwt.Token) (interface{}, error) {
	// 	return []byte(consts.JWT_SECRET), nil
	// })

	// if err != nil {
	// 	return MakeDataResp(c, &JWTValidResp{JwtIsValid: false})
	// }

	// claims := token.Claims.(*JwtCustomClaims)

	// user := c.Get("user").(*jwt.Token)
	// claims := user.Claims.(*JwtCustomClaims)

	// IpAddr := c.RealIP()

	// log.Debug().Msgf("ip: %s, %s", IpAddr, claims.IpAddr)

	// //t := claims.ExpiresAt.Unix()
	// //expired := t != 0 && t < time.Now().Unix()

	// if IpAddr != claims.IpAddr {
	// 	return MakeDataResp(c, &JWTValidResp{JwtIsValid: false})
	// }

	return MakeDataResp(c, "", &JwtValidResp{JwtIsValid: true})

}

func RenewToken(c echo.Context) error {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*auth.JwtCustomClaims)

	// Throws unauthorized error
	//if username != "edb" || password != "tod4EwVHEyCRK8encuLE" {
	//	return echo.ErrUnauthorized
	//}

	// Set custom claims
	renewClaims := auth.JwtCustomClaims{
		UserId: claims.UserId,
		//Email: authUser.Email,
		IpAddr: claims.IpAddr,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * auth.JWT_TOKEN_EXPIRES_HOURS)),
		},
	}

	// Create token with claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, renewClaims)

	// Generate encoded token and send it as response.
	t, err := token.SignedString([]byte(consts.JWT_SECRET))

	if err != nil {
		return BadReq("error signing token")
	}

	return MakeDataResp(c, "", &JwtResp{t})
}

func GetJwtInfoFromRoute(c echo.Context) *JwtInfo {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*auth.JwtCustomClaims)

	t := claims.ExpiresAt.Unix()
	expired := t != 0 && t < time.Now().Unix()

	return &JwtInfo{UserId: claims.UserId, IpAddr: claims.IpAddr, Expires: time.Unix(t, 0).String(), Expired: expired}
}

func JWTInfoRoute(c echo.Context) error {
	info := GetJwtInfoFromRoute(c)

	return MakeDataResp(c, "", info)
}
