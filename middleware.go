package main

import (
	"fmt"

	"github.com/antonybholmes/go-edb-api/routes"
	authroutes "github.com/antonybholmes/go-edb-api/routes/auth"
	"github.com/labstack/echo-contrib/session"
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

func JwtIsAccessTokenMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		_, err := routes.NewValidator(c).IsValidAccessToken().Ok()

		if err != nil {
			return err
		}

		return next(c)
	}
}

func SessionIsValidMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		sess, err := session.Get(authroutes.SESSION_NAME, c)
		if err != nil {
			return err
		}

		//log.Debug().Msgf("validate session %s", sess.ID)

		_, ok := sess.Values[authroutes.SESSION_UUID].(string)

		if !ok {
			return fmt.Errorf("cannot get user id from session")
		}

		return next(c)
	}
}
