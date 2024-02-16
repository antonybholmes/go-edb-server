package main

import (
	"net/http"
	"time"

	"github.com/antonybholmes/go-auth"
	"github.com/antonybholmes/go-edb-api/consts"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

type JWTResp struct {
	JWT string `json:"jwt"`
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

type ReqLogin struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func RegisterRoute(c echo.Context, userdb *auth.UserDb, secret string) error {
	login := new(ReqLogin)

	err := c.Bind(login)

	if err != nil {
		return err
	}

	//email := c.FormValue("email")
	//password := c.FormValue("password")

	user := auth.NewLoginUser(login.Email, login.Password)

	authUser, err := userdb.CreateUser(user)

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

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
	t, err := token.SignedString([]byte(secret))

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	return MakeDataResp(c, &JWTResp{t}) //c.JSON(http.StatusOK, JWTResp{t})
}

func LoginRoute(c echo.Context, userdb *auth.UserDb) error {
	login := new(ReqLogin)

	err := c.Bind(login)

	if err != nil {
		return err
	}

	//email := c.FormValue("email")
	//password := c.FormValue("password")

	user := auth.NewLoginUser(login.Email, login.Password)

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

	return MakeDataResp(c, &JWTResp{t})
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

	return MakeDataResp(c, &JWTResp{t})
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

	return MakeDataResp(c, info)
}
