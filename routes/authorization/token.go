package authorization

import (
	"github.com/antonybholmes/go-auth"
	"github.com/antonybholmes/go-auth/tokengen"
	"github.com/antonybholmes/go-edb-server/consts"
	"github.com/antonybholmes/go-edb-server/routes"
	authenticationroutes "github.com/antonybholmes/go-edb-server/routes/authentication"
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
// 		return routes.ErrorReq("error signing token")
// 	}

// 	return MakeDataResp(c, "", &JwtResp{t})
// }

func TokenInfoRoute(c echo.Context) error {
	t, err := routes.HeaderAuthToken(c)

	if err != nil {
		return routes.ErrorReq(err)
	}

	claims := auth.TokenClaims{}

	_, err = jwt.ParseWithClaims(t, &claims, func(token *jwt.Token) (interface{}, error) {
		return consts.JWT_RSA_PUBLIC_KEY, nil
	})

	if err != nil {
		return routes.ErrorReq(err)
	}

	return routes.MakeDataPrettyResp(c, "", &routes.JwtInfo{
		Uuid: claims.Uuid,
		Type: claims.Type, //.TokenTypeString(claims.Type),
		//IpAddr:  claims.IpAddr,
		Expires: claims.ExpiresAt.UTC().String()})

}

func NewAccessTokenRoute(c echo.Context) error {
	return authenticationroutes.NewValidator(c).CheckIsValidRefreshToken().Success(func(validator *authenticationroutes.Validator) error {

		// Generate encoded token and send it as response.
		accessToken, err := tokengen.AccessToken(c, validator.Claims.Uuid, validator.Claims.Roles)

		if err != nil {
			return routes.ErrorReq("error creating access token")
		}

		return routes.MakeDataPrettyResp(c, "", &routes.AccessTokenResp{AccessToken: accessToken})
	})

}
