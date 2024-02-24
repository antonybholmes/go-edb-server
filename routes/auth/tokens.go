package auth

import (
	"github.com/antonybholmes/go-auth"
	"github.com/antonybholmes/go-edb-api/consts"
	"github.com/antonybholmes/go-edb-api/routes"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

type AuthReq struct {
	Authorization string `header:"Authorization"`
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
// 		return routes.BadReq("error signing token")
// 	}

// 	return MakeDataResp(c, "", &JwtResp{t})
// }

func TokenInfoRoute(c echo.Context) error {
	t, err := routes.HeaderAuthToken(c)

	if err != nil {
		return routes.BadReq(err)
	}

	claims := auth.JwtCustomClaims{}

	_, err = jwt.ParseWithClaims(t, &claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(consts.JWT_SECRET), nil
	})

	if err != nil {
		return routes.BadReq(err)
	}

	return routes.MakeDataResp(c, "", &routes.JwtInfo{
		Uuid: claims.Uuid,
		Type: auth.TokenTypeString(claims.Type),
		//IpAddr:  claims.IpAddr,
		Expires: claims.ExpiresAt.UTC().String()})

}

func NewAccessTokenRoute(c echo.Context) error {
	return routes.IsValidRefreshTokenCB(c, func(c echo.Context, claims *auth.JwtCustomClaims) error {
		// Generate encoded token and send it as response.
		t, err := auth.AccessToken(c, claims.Uuid, consts.JWT_SECRET)

		if err != nil {
			return routes.BadReq("error creating access token")
		}

		return routes.MakeDataResp(c, "", &routes.AccessTokenResp{AccessToken: t})
	})

}
