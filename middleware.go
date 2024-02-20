package main

import (
	"github.com/antonybholmes/go-auth"
	"github.com/antonybholmes/go-edb-api/routes"
	"github.com/golang-jwt/jwt/v5"
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
// 			return routes.BadReq("ip address invalid")
// 		}

// 		return next(c)
// 	}
// }

func JwtCheckMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {

		user := c.Get("user").(*jwt.Token)
		claims := user.Claims.(*auth.JwtCustomClaims)

		IpAddr := c.RealIP()

		log.Debug().Msgf("ip: %s, %s", IpAddr, claims.IpAddr)

		//t := claims.ExpiresAt.Unix()
		//expired := t != 0 && t < time.Now().Unix()

		if IpAddr != claims.IpAddr {
			return routes.BadReq("ip address invalid")
		}

		return next(c)
	}
}

func JwtIsAccessMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {

		user := c.Get("user").(*jwt.Token)
		claims := user.Claims.(*auth.JwtCustomClaims)

		log.Debug().Msgf("type: %s, %s", claims.Type, auth.TOKEN_TYPE_ACCESS)

		//t := claims.ExpiresAt.Unix()
		//expired := t != 0 && t < time.Now().Unix()

		if claims.Type != auth.TOKEN_TYPE_ACCESS {
			return routes.BadReq("wrong token type")
		}

		return next(c)
	}
}
