package routes

import (
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

// See https://echo.labstack.com/docs/cookbook/jwt#login

const JWT_TOKEN_EXPIRES_HOURS time.Duration = 24

const INVALID_JWT_MESSAGE string = "Invalid JWT"

type JWTInfo struct {
	Id    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

type JwtCustomClaims struct {
	Id    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	jwt.RegisteredClaims
}

func LoginRoute(c echo.Context, secret string) error {
	username := c.FormValue("username")
	password := c.FormValue("password")

	// Throws unauthorized error
	if username != "edb" || password != "tod4EwVHEyCRK8encuLE" {
		return echo.ErrUnauthorized
	}

	// Set custom claims
	claims := &JwtCustomClaims{
		"05c80ad1913248d4880dbc2f496cb151",
		"edb",
		"antony@antonyholmes.dev",
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24)),
		},
	}

	// Create token with claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Generate encoded token and send it as response.
	t, err := token.SignedString([]byte(secret))
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, echo.Map{
		"jwt": t,
	})
}

func GetJwtInfoFromRoute(c echo.Context) *JWTInfo {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*JwtCustomClaims)

	return &JWTInfo{Id: claims.Id, Name: claims.Name, Email: claims.Email}
}

func IsValidJwtInfo(jwtInfo *JWTInfo) bool {
	return jwtInfo.Name == "edb"
}

func JWTInfoRoute(c echo.Context) error {
	info := GetJwtInfoFromRoute(c)

	return c.JSONPretty(http.StatusOK, info, "  ")
}
