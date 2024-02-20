package routes

import (
	"strings"
	"time"

	"github.com/antonybholmes/go-auth"
	"github.com/antonybholmes/go-edb-api/consts"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
)

type AuthReq struct {
	Authorization string `header:"authorization"`
}

func TokenValidRoute(c echo.Context) error {
	// jwtReq := new(ReqJwt)

	// err := c.Bind(jwtReq)

	// if err != nil {
	// 	return err
	// }

	// token, err := jwt.ParseWithClaims(jwtReq.Jwt, &JwtCustomClaims{}, func(token *jwt.Token) (interface{}, error) {
	// 	return []byte(consts.JWT_SECRET), nil
	// })

	// if err != nil {
	// 	return MakeDataResp(c, &JWTValidResp{JwtIsValid: false})
	// }

	// claims := token.Claims.(*JwtCustomClaims)

	// user := c.Get("user").(*jwt.Token)
	// claims := user.Claims.(*JwtCustomClaims)

	// IpAddr := c.RealIP()

	// log.Debug().Msgf("ip: %s, %s", IpAddr, claims.IpAddr)

	// //t := claims.ExpiresAt.Unix()
	// //expired := t != 0 && t < time.Now().Unix()

	// if IpAddr != claims.IpAddr {
	// 	return MakeDataResp(c, &JWTValidResp{JwtIsValid: false})
	// }

	return MakeValidResp(c, "", true)

}

// func RenewTokenRoute(c echo.Context) error {
// 	user := c.Get("user").(*jwt.Token)
// 	claims := user.Claims.(*auth.JwtCustomClaims)

// 	// Throws unauthorized error
// 	//if username != "edb" || password != "tod4EwVHEyCRK8encuLE" {
// 	//	return echo.ErrUnauthorized
// 	//}

// 	// Set custom claims
// 	renewClaims := auth.JwtCustomClaims{
// 		UserId: claims.UserId,
// 		//Email: authUser.Email,
// 		IpAddr: claims.IpAddr,
// 		RegisteredClaims: jwt.RegisteredClaims{
// 			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * auth.JWT_TOKEN_EXPIRES_HOURS)),
// 		},
// 	}

// 	// Create token with claims
// 	token := jwt.NewWithClaims(jwt.SigningMethodHS256, renewClaims)

// 	// Generate encoded token and send it as response.
// 	t, err := token.SignedString([]byte(consts.JWT_SECRET))

// 	if err != nil {
// 		return BadReq("error signing token")
// 	}

// 	return MakeDataResp(c, "", &JwtResp{t})
// }

func NewAccessTokenRoute(c echo.Context) error {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*auth.JwtCustomClaims)

	// Generate encoded token and send it as response.
	t, err := auth.CreateAccessToken(claims.UserId, c.RealIP(), consts.JWT_SECRET)

	if err != nil {
		return MakeDataResp(c, "error signing token", &JwtResp{""})
	}

	return MakeDataResp(c, "", &JwtResp{t})
}

func TokenInfoRoute(c echo.Context) error {

	h := c.Request().Header.Get("Authorization")

	if h == "" {
		return BadReq("authorization header not present")
	}

	if !strings.Contains(h, "Bearer") {
		return BadReq("bearer not present")
	}

	tokens := strings.Split(h, " ")

	if len(tokens) < 2 {
		return BadReq("jwt not present")
	}

	t := tokens[1]

	log.Debug().Msgf("%s", t)

	claims := auth.JwtCustomClaims{}

	_, err := jwt.ParseWithClaims(t, &claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(consts.JWT_SECRET), nil
	})

	if err != nil {
		return BadReq(err)
	}

	log.Debug().Msgf("%s %s", claims.ExpiresAt.UTC(), time.Now().UTC())

	// user := c.Get("user").(*jwt.Token)
	// claims := user.Claims.(*auth.JwtCustomClaims)

	expired := claims.ExpiresAt.UTC().Before(time.Now().UTC())

	return MakeDataResp(c, "", &JwtInfo{
		UserId:  claims.UserId,
		Type:    claims.Type,
		IpAddr:  claims.IpAddr,
		Expires: claims.ExpiresAt.UTC().String(),
		Expired: expired})

}
