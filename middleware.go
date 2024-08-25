package main

import (
	"errors"
	"fmt"

	"strings"

	"github.com/antonybholmes/go-auth"
	"github.com/antonybholmes/go-edb-server/consts"
	"github.com/antonybholmes/go-edb-server/routes"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
)

// func JwtOtpCheckMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
// 	return func(c echo.Context) error {

// 		user := c.Get("user").(*jwt.Token)
// 		claims := user.Claims.(*auth.JwtOtpCustomClaims)

// 		IpAddr := c.RealIP()

// 		//log.Debug().Msgf("ip: %s, %s", IpAddr, claims.IpAddr)

// 		//t := claims.ExpiresAt.Unix()
// 		//expired := t != 0 && t < time.Now().Unix()

// 		if IpAddr != claims.IpAddr {
// 			return routes.ErrorReq("ip address invalid")
// 		}

// 		return next(c)
// 	}
// }

// func JwtCheckMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
// 	return func(c echo.Context) error {

// 		user := c.Get("user").(*jwt.Token)
// 		claims := user.Claims.(*auth.JwtCustomClaims)

// 		IpAddr := c.RealIP()

// 		log.Debug().Msgf("ip: %s, %s", IpAddr, claims.IpAddr)

// 		//t := claims.ExpiresAt.Unix()
// 		//expired := t != 0 && t < time.Now().Unix()

// 		if IpAddr != claims.IpAddr {
// 			return routes.ErrorReq("ip address invalid")
// 		}

// 		return next(c)
// 	}
// }

type CustomContext struct {
	echo.Context
}

// func JwtLoadTokenClaimsMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
// 	return func(c echo.Context) error {
// 		_, err := routes.NewValidator(c).CheckIsValidAccessToken().Ok()

// 		if err != nil {
// 			return err
// 		}

// 		return next(c)
// 	}
// }

func JwtIsRefreshTokenMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		user := c.Get("user").(*jwt.Token)
		claims := user.Claims.(*auth.JwtCustomClaims)

		if claims.Type != auth.TOKEN_TYPE_REFRESH {
			routes.AuthErrorReq("not a refresh token")
		}

		return next(c)
	}
}

func JwtIsAccessTokenMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		user := c.Get("user").(*jwt.Token)
		claims := user.Claims.(*auth.JwtCustomClaims)

		if claims.Type != auth.TOKEN_TYPE_ACCESS {
			routes.AuthErrorReq("not an access token")
		}

		return next(c)
	}
}

func JwtHasAdminPermissionMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		user := c.Get("user").(*jwt.Token)
		claims := user.Claims.(*auth.JwtCustomClaims)

		if !auth.IsAdmin((claims.Roles)) {
			return routes.AuthErrorReq("user is not an admin")
		}

		return next(c)
	}
}

func JwtHasLoginPermissionMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		user := c.Get("user").(*jwt.Token)
		claims := user.Claims.(*auth.JwtCustomClaims)

		if !auth.CanLogin((claims.Roles)) {
			return routes.AuthErrorReq("user is not allowed to login")
		}

		return next(c)
	}
}

func SessionIsValidMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		sess, err := session.Get(routes.SESSION_NAME, c)
		if err != nil {
			return err
		}

		//log.Debug().Msgf("validate session %s", sess.ID)

		_, ok := sess.Values[routes.SESSION_PUBLICID].(string)

		if !ok {
			return errors.New("cannot get user id from session")
		}

		return next(c)
	}
}

func ValidateJwtMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		authorizationHeader := c.Request().Header.Get("authorization")

		if len(authorizationHeader) == 0 {
			return routes.AuthErrorReq("missing Authentication header")

		}

		log.Debug().Msgf("parsing authentication header")

		authPair := strings.SplitN(authorizationHeader, " ", 2)

		if len(authPair) != 2 {
			return routes.AuthErrorReq("wrong Authentication header definiton")
		}

		headerAuthScheme := authPair[0]
		headerAuthToken := authPair[1]

		if headerAuthScheme != "Bearer" {
			return routes.AuthErrorReq("wrong Authentication header definiton")
		}

		log.Debug().Msgf("validating JWT token")

		token, err := validateJwtToken(headerAuthToken)

		if err != nil {
			return routes.AuthErrorReq(err)
		}

		log.Debug().Msgf("JWT token is valid")
		c.Set("user", token)
		return next(c)

	}
}

// you can add your implementation here.
func validateJwtToken(tokenString string) (*jwt.Token, error) {

	token, err := jwt.ParseWithClaims(tokenString, &auth.JwtCustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return consts.JWT_PUBLIC_KEY, nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	return token, nil
}

// Create a permissions middleware to verify jwt permissions on a token
func NewJwtRoleMiddleware(validRoles ...string) echo.MiddlewareFunc {

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {

			user := c.Get("user").(*jwt.Token)

			// if user == nil {
			// 	return routes.AuthErrorReq("no jwt available")
			// }

			claims := user.Claims.(*auth.JwtCustomClaims)

			// shortcut for admin, as we allow this for everything
			if auth.IsAdmin(claims.Roles) {
				return next(c)
			}

			for _, validRole := range validRoles {
				for _, role := range claims.Roles {

					// if we find a permission, stop and move on
					if strings.Contains(role, validRole) {
						return next(c)
					}
				}
			}

			return routes.AuthErrorReq("roles not found")
		}
	}
}
