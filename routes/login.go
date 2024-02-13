package routes

import (
	"fmt"
	"time"

	"github.com/antonybholmes/go-edb-api/auth"
	"github.com/antonybholmes/go-utils"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

type JWTResp struct {
	JWT string `json:"jwt"`
}

type JWTInfo struct {
	Id string `json:"id"`
	//Name  string `json:"name"`
	Email   string `json:"email"`
	Expires string `json:"expires"`
	Expired bool   `json:"expired"`
}

type JwtCustomClaims struct {
	Id string `json:"id"`
	//Name  string `json:"name"`
	Email string `json:"email"`
	jwt.RegisteredClaims
}

func RegisterRoute(c echo.Context, userdb *auth.UserDb, secret string) error {
	email := c.FormValue("email")
	password := c.FormValue("password")

	user := auth.NewLoginUser(email, password)

	authUser, err := userdb.CreateUser(user)

	if err != nil {
		return utils.MakeBadResp(c, err)
	}

	// Set custom claims
	claims := &JwtCustomClaims{
		authUser.UserId,
		authUser.Email,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * auth.JWT_TOKEN_EXPIRES_HOURS)),
		},
	}

	// Create token with claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Generate encoded token and send it as response.
	t, err := token.SignedString([]byte(secret))

	if err != nil {
		return utils.MakeBadResp(c, err)
	}

	return utils.MakeDataResp(c, &JWTResp{t}) //c.JSON(http.StatusOK, JWTResp{t})
}

func LoginRoute(c echo.Context, userdb *auth.UserDb, secret string) error {
	email := c.FormValue("email")
	password := c.FormValue("password")

	user := auth.NewLoginUser(email, password)

	authUser, err := userdb.FindUserByEmail(user)

	if err != nil {
		return utils.MakeBadResp(c, fmt.Errorf("user does not exist"))
	}

	if !authUser.CheckPasswords(user.Password) {
		return utils.MakeBadResp(c, fmt.Errorf("incorrect password"))
	}

	// Throws unauthorized error
	//if username != "edb" || password != "tod4EwVHEyCRK8encuLE" {
	//	return echo.ErrUnauthorized
	//}

	// Set custom claims
	claims := &JwtCustomClaims{
		authUser.UserId,
		authUser.Email,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * auth.JWT_TOKEN_EXPIRES_HOURS)),
		},
	}

	// Create token with claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Generate encoded token and send it as response.
	t, err := token.SignedString([]byte(secret))

	if err != nil {
		return utils.MakeBadResp(c, fmt.Errorf("error signing token"))
	}

	return utils.MakeDataResp(c, &JWTResp{t})
}

func GetJwtInfoFromRoute(c echo.Context) *JWTInfo {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*JwtCustomClaims)

	t := claims.ExpiresAt.Unix()
	expired := t != 0 && t < time.Now().Unix()

	return &JWTInfo{Id: claims.Id, Email: claims.Email, Expires: time.Unix(t, 0).String(), Expired: expired}
}

func JWTInfoRoute(c echo.Context) error {
	info := GetJwtInfoFromRoute(c)

	return utils.MakeDataResp(c, info)
}
