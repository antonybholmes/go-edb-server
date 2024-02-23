package routes

import (
	"net/mail"

	"github.com/antonybholmes/go-auth"
	"github.com/antonybholmes/go-auth/userdb"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
)

//
// Standardized data checkers for checking header and body contain
// the correct data for a route
//

func ValidEmailCB(c echo.Context,
	email string,
	callback func(c echo.Context, email *mail.Address) error) error {
	address, err := mail.ParseAddress(email)

	if err != nil {
		return InvalidEmailReq()
	}

	return callback(c, address)
}

func AuthUserFromEmailCB(c echo.Context,
	email string,
	callback func(c echo.Context, authUser *auth.AuthUser) error) error {
	return ValidEmailCB(c, email, func(c echo.Context, email *mail.Address) error {

		authUser, err := userdb.FindUserByEmail(email)

		if err != nil {
			return UserDoesNotExistReq()
		}

		return callback(c, authUser)
	})
}

func VerifiedEmailCB(c echo.Context,
	authUser *auth.AuthUser,
	callback func(c echo.Context, authUser *auth.AuthUser) error) error {

	if !authUser.EmailVerified {
		return BadReq("email address not verified")
	}

	return callback(c, authUser)
}

func JwtCB(c echo.Context,
	callback func(c echo.Context, claims *auth.JwtCustomClaims) error) error {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*auth.JwtCustomClaims)

	return callback(c, claims)
}

// Extracts uuid from token, checks user exists and calls callback function.
// If claims argument is nil, function will search for claims automatically.
// If claims are supplied, this step is skipped. This is so this function can
// be nested in other call backs that may have already extracted the claims
// without having to repeat this part.
func UserFromUuidCB(c echo.Context,
	claims *auth.JwtCustomClaims,
	callback func(c echo.Context, claims *auth.JwtCustomClaims, authUser *auth.AuthUser) error) error {
	if claims == nil {
		// if no claims specified, extract the claims and run function with claims
		return JwtCB(c, func(c echo.Context, claims *auth.JwtCustomClaims) error {
			return UserFromUuidCB(c, claims, callback)
		})

	}

	log.Debug().Msgf("from uuiid %s", claims.Uuid)

	authUser, err := userdb.FindUserByUuid(claims.Uuid)

	if err != nil {
		return UserDoesNotExistReq()
	}

	return callback(c, claims, authUser)
}

func IsValidRefreshTokenCB(c echo.Context,
	callback func(c echo.Context, claims *auth.JwtCustomClaims) error) error {
	return JwtCB(c, func(c echo.Context, claims *auth.JwtCustomClaims) error {
		if claims.Type != auth.TOKEN_TYPE_REFRESH {
			return BadReq("wrong token type")
		}

		return callback(c, claims)
	})
}

func IsValidAccessTokenCB(c echo.Context,
	callback func(c echo.Context, claims *auth.JwtCustomClaims) error) error {
	return JwtCB(c, func(c echo.Context, claims *auth.JwtCustomClaims) error {
		if claims.Type != auth.TOKEN_TYPE_ACCESS {
			return BadReq("wrong token type")
		}

		return callback(c, claims)
	})
}
