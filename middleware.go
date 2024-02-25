package main

import (
	"github.com/antonybholmes/go-edb-api/routes"
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
// 			return routes.BadReq("ip address invalid")
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
// 			return routes.BadReq("ip address invalid")
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
