package routes

import (
	"strings"

	"github.com/antonybholmes/go-auth"
	"github.com/antonybholmes/go-edb-api/consts"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
)

type AuthReq struct {
	Authorization string `header:"authorization"`
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

	return MakeDataResp(c, "", &JwtInfo{
		UserId:  claims.UserId,
		Type:    auth.TokenTypeString(claims.Type),
		IpAddr:  claims.IpAddr,
		Expires: claims.ExpiresAt.UTC().String()})

}

func NewAccessTokenRoute(c echo.Context) error {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*auth.JwtCustomClaims)

	if claims.Type != auth.TOKEN_TYPE_REFRESH {
		return BadReq("wrong token type")
	}

	// Generate encoded token and send it as response.
	t, err := auth.AccessToken(claims.UserId, c.RealIP(), consts.JWT_SECRET)

	if err != nil {
		return MakeDataResp(c, "error signing token", &JwtResp{""})
	}

	return MakeDataResp(c, "", &JwtResp{t})
}
