package routes

import (
	"net/mail"

	"github.com/antonybholmes/go-auth"
	"github.com/antonybholmes/go-auth/userdb"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

func ValidEmailCB(c echo.Context, email string, callback func(c echo.Context, email *mail.Address) error) error {
	address, err := mail.ParseAddress(email)

	if err != nil {
		return InvalidEmailReq()
	}

	return callback(c, address)
}

func EmailUserCB(c echo.Context, email *mail.Address, callback func(c echo.Context, authUser *auth.AuthUser) error) error {
	authUser, err := userdb.FindUserByEmail(email)

	if err != nil {
		return UserDoesNotExistReq()
	}

	return callback(c, authUser)
}

func UuidUserCB(c echo.Context, claims *auth.JwtCustomClaims, callback func(c echo.Context, claims *auth.JwtCustomClaims, authUser *auth.AuthUser) error) error {
	authUser, err := userdb.FindUserByUsername(claims.Uuid)

	if err != nil {
		return UserDoesNotExistReq()
	}

	return callback(c, claims, authUser)
}

func VerifiedEmailCB(c echo.Context, authUser *auth.AuthUser, callback func(c echo.Context, authUser *auth.AuthUser) error) error {

	if !authUser.EmailVerified {
		return BadReq("email address not verified")
	}

	return callback(c, authUser)
}

func JwtCB(c echo.Context, callback func(c echo.Context, claims *auth.JwtCustomClaims) error) error {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*auth.JwtCustomClaims)

	return callback(c, claims)
}

func RefreshTokenCB(c echo.Context,
	callback func(c echo.Context, claims *auth.JwtCustomClaims) error) error {
	return JwtCB(c, func(c echo.Context, claims *auth.JwtCustomClaims) error {
		if claims.Type < auth.TOKEN_TYPE_REFRESH {
			return BadReq("wrong token type")
		}

		return callback(c, claims)
	})
}

func AccessTokenCB(c echo.Context, callback func(c echo.Context, claims *auth.JwtCustomClaims) error) error {
	return JwtCB(c, func(c echo.Context, claims *auth.JwtCustomClaims) error {
		if claims.Type < auth.TOKEN_TYPE_ACCESS {
			return BadReq("wrong token type")
		}

		return callback(c, claims)
	})
}
