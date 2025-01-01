package main

import (
	"strings"

	"github.com/antonybholmes/go-auth"
	"github.com/antonybholmes/go-edb-server/routes"
	authenticationroutes "github.com/antonybholmes/go-edb-server/routes/authentication"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
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

// func JwtLoadTokenClaimsMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
// 	return func(c echo.Context) error {
// 		_, err := NewValidator(c).CheckIsValidAccessToken().Ok()

// 		if err != nil {
// 			return err
// 		}

// 		return next(c)
// 	}
// }

func JwtIsRefreshTokenMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		user := c.Get("user").(*jwt.Token)
		claims := user.Claims.(*auth.TokenClaims)

		if claims.Type != auth.REFRESH_TOKEN {
			routes.AuthErrorReq("not a refresh token")
		}

		return next(c)
	}
}

func JwtIsAccessTokenMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		user := c.Get("user").(*jwt.Token)
		claims := user.Claims.(*auth.TokenClaims)

		if claims.Type != auth.ACCESS_TOKEN {
			routes.AuthErrorReq("not an access token")
		}

		return next(c)
	}
}

func JwtHasAdminPermissionMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		user := c.Get("user").(*jwt.Token)
		claims := user.Claims.(*auth.TokenClaims)

		if !auth.IsAdmin((claims.Roles)) {
			return routes.AuthErrorReq("user is not an admin")
		}

		return next(c)
	}
}

func JwtHasLoginPermissionMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		user := c.Get("user").(*jwt.Token)
		claims := user.Claims.(*auth.TokenClaims)

		if !auth.CanSignin((claims.Roles)) {
			return routes.AuthErrorReq("user is not allowed to login")
		}

		return next(c)
	}
}

// basic check that session exists and seems to be populated with the user
func SessionIsValidMiddleware(next echo.HandlerFunc) echo.HandlerFunc {

	return func(c echo.Context) error {
		sessData, err := authenticationroutes.ReadSessionInfo(c)

		if err != nil {
			return routes.AuthErrorReq("cannot get user id from session")
		}

		c.Set("authUser", sessData.AuthUser)

		return next(c)
	}
}

// func ValidateJwtMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
// 	return func(c echo.Context) error {
// 		authorizationHeader := c.Request().Header.Get("authorization")

// 		if len(authorizationHeader) == 0 {
// 			return routes.AuthErrorReq("missing Authentication header")

// 		}

// 		log.Debug().Msgf("parsing authentication header")

// 		authPair := strings.SplitN(authorizationHeader, " ", 2)

// 		if len(authPair) != 2 {
// 			return routes.AuthErrorReq("wrong Authentication header definiton")
// 		}

// 		headerAuthScheme := authPair[0]
// 		headerAuthToken := authPair[1]

// 		if headerAuthScheme != "Bearer" {
// 			return routes.AuthErrorReq("wrong Authentication header definiton")
// 		}

// 		log.Debug().Msgf("validating JWT token")

// 		token, err := validateJwtToken(headerAuthToken)

// 		if err != nil {
// 			return routes.AuthErrorReq(err)
// 		}

// 		log.Debug().Msgf("JWT token is valid")
// 		c.Set("user", token)
// 		return next(c)

// 	}
// }

// Create a permissions middleware to verify jwt permissions on a token
func JwtRoleMiddleware(validRoles ...string) echo.MiddlewareFunc {

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {

			user := c.Get("user").(*jwt.Token)

			// if user == nil {
			// 	return routes.AuthErrorReq("no jwt available")
			// }

			claims := user.Claims.(*auth.TokenClaims)

			// shortcut for admin, as we allow this for everything
			if auth.IsAdmin(claims.Roles) {
				//log.Debug().Msgf("is admin")
				return next(c)
			}

			for _, validRole := range validRoles {

				// if we find a permission, stop and move on
				if strings.Contains(claims.Roles, validRole) {
					return next(c)
				}

			}

			return routes.AuthErrorReq("roles not found")
		}
	}
}

func RDFMiddleware() echo.MiddlewareFunc {
	return JwtRoleMiddleware("RDF")
}
