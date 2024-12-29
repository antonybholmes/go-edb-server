package auth0routes

import (
	"github.com/antonybholmes/go-auth"
	"github.com/antonybholmes/go-edb-server/routes"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
)

func ValidateAuth0TokenRoute(c echo.Context) error {

	// bytes, err := os.ReadFile("auth0.key.pub")
	// if err != nil {
	// 	log.Fatal().Msgf("%s", err)
	// }

	// key, err := jwt.ParseRSAPublicKeyFromPEM(bytes)
	// if err != nil {
	// 	log.Fatal().Msgf("%s", err)
	// }

	// h := c.Request().Header.Get("Authorization")

	// tokens := strings.SplitN(h, " ", 2)
	// token := tokens[1]

	// log.Debug().Msgf("tok: %v", h)

	// hmm, err := jwt.ParseWithClaims(token, &TokenClaims{}, func(token *jwt.Token) (interface{}, error) {
	// 	return key, nil
	// })

	// if err != nil {
	// 	log.Debug().Msgf("%s", err)
	// }

	user := c.Get("user").(*jwt.Token)
	myClaims := user.Claims.(*auth.Auth0TokenClaims)

	//myClaims := user.Claims.(*auth.TokenClaims) //hmm.Claims.(*TokenClaims)

	//user := c.Get("user").(*jwt.Token)
	//claims := user.Claims.(*TokenClaims)

	log.Debug().Msgf("auth0 claims %v", myClaims)

	log.Debug().Msgf("auth0 claims %v", myClaims.Email)

	return routes.MakeOkPrettyResp(c, "user was signed out")
}
