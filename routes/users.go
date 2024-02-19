package routes

import (
	"bytes"
	"fmt"
	"html/template"
	"time"

	"github.com/antonybholmes/go-auth"
	"github.com/antonybholmes/go-edb-api/consts"
	"github.com/antonybholmes/go-email/email"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

type JWTResp struct {
	Jwt string `json:"jwt"`
}

type JWTValidResp struct {
	JwtIsValid bool `json:"jwtIsValid"`
}

type JWTInfo struct {
	UserId string `json:"userId"`
	//Name  string `json:"name"`
	//Email   string `json:"email"`
	IpAddr  string `json:"ipAddr"`
	Expires string `json:"expires"`
	Expired bool   `json:"expired"`
}

type JwtCustomClaims struct {
	UserId string `json:"userId"`
	//Name  string `json:"name"`
	//Email string `json:"email"`
	IpAddr string `json:"ipAddr"`
	jwt.RegisteredClaims
}

type ReqJwt struct {
	Jwt string `json:"jwt"`
}

func RegisterRoute(c echo.Context, userdb *auth.UserDb, secret string) error {
	req := new(auth.LoginReq)

	err := c.Bind(req)

	if err != nil {
		return err
	}

	//email := c.FormValue("email")
	//password := c.FormValue("password")

	user := auth.NewLoginUser(req.Name, req.Email, req.Password)

	otp := auth.AuthCode()

	_, err = userdb.CreateUser(user, otp)

	if err != nil {
		return BadReq(err)
	}

	mimeHeaders := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	var body bytes.Buffer
	body.Write([]byte(fmt.Sprintf("Subject: Email verification \n%s\n\n", mimeHeaders)))

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

	if req.Url != "" {
		err = t.Execute(&body, struct {
			Name string
			Link string
		}{
			Name: user.Name,
			Link: fmt.Sprintf("%sotp=%s", req.Url, otp),
		})

		if err != nil {
			return BadReq(err)
		}
	} else {
		err = t.Execute(&body, struct {
			Name string
			Code string
		}{
			Name: user.Name,
			Code: otp,
		})

		if err != nil {
			return BadReq(err)
		}
	}

	err = email.SendEmail(user.Email, body.Bytes())

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

func LoginRoute(c echo.Context, userdb *auth.UserDb) error {
	req := new(auth.LoginReq)

	err := c.Bind(req)

	if err != nil {
		return err
	}

	//email := c.FormValue("email")
	//password := c.FormValue("password")

	user := auth.LoginUserFromReq(req)

	authUser, err := userdb.FindUserByEmail(user)

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
	claims := &JwtCustomClaims{
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

	return MakeDataResp(c, "", &JWTResp{t})
}

func ValidateTokenRoute(c echo.Context) error {
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

	return MakeDataResp(c, "", &JWTValidResp{JwtIsValid: true})

}

func RefreshTokenRoute(c echo.Context) error {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*JwtCustomClaims)

	// Throws unauthorized error
	//if username != "edb" || password != "tod4EwVHEyCRK8encuLE" {
	//	return echo.ErrUnauthorized
	//}

	// Set custom claims
	refreshedClaims := &JwtCustomClaims{
		UserId: claims.UserId,
		//Email: authUser.Email,
		IpAddr: claims.IpAddr,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * auth.JWT_TOKEN_EXPIRES_HOURS)),
		},
	}

	// Create token with claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshedClaims)

	// Generate encoded token and send it as response.
	t, err := token.SignedString([]byte(consts.JWT_SECRET))

	if err != nil {
		return BadReq("error signing token")
	}

	return MakeDataResp(c, "", &JWTResp{t})
}

func GetJwtInfoFromRoute(c echo.Context) *JWTInfo {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*JwtCustomClaims)

	t := claims.ExpiresAt.Unix()
	expired := t != 0 && t < time.Now().Unix()

	return &JWTInfo{UserId: claims.UserId, IpAddr: claims.IpAddr, Expires: time.Unix(t, 0).String(), Expired: expired}
}

func JWTInfoRoute(c echo.Context) error {
	info := GetJwtInfoFromRoute(c)

	return MakeDataResp(c, "", info)
}
