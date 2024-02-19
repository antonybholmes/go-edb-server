package main

import (
	"github.com/antonybholmes/go-edb-api/routes"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
)

func JWTCheckMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {

		user := c.Get("user").(*jwt.Token)
		claims := user.Claims.(*routes.JwtCustomClaims)

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
