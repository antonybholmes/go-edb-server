package consts

import (
	"crypto/rsa"
	"os"

	"github.com/rs/zerolog/log"

	"github.com/antonybholmes/go-sys/env"

	"github.com/golang-jwt/jwt/v5"
)

var NAME string
var VERSION string
var COPYRIGHT string
var JWT_PRIVATE_KEY *rsa.PrivateKey //[]byte
var JWT_PUBLIC_KEY *rsa.PublicKey   //[]byte
var SESSION_SECRET string

func init() {
	env.Load()

	NAME = os.Getenv("NAME")
	VERSION = os.Getenv("VERSION")
	COPYRIGHT = os.Getenv("COPYRIGHT")
	//JWT_PRIVATE_KEY = []byte(os.Getenv("JWT_SECRET"))
	//JWT_PUBLIC_KEY = []byte(os.Getenv("JWT_SECRET"))
	SESSION_SECRET = os.Getenv("SESSION_SECRET")

	signBytes, err := os.ReadFile("jwtRS256.key")
	if err != nil {
		log.Fatal().Msgf("%s", err)
	}

	JWT_PRIVATE_KEY, err = jwt.ParseRSAPrivateKeyFromPEM(signBytes)
	if err != nil {
		log.Fatal().Msgf("%s", err)
	}

	verifyBytes, err := os.ReadFile("jwtRS256.key.pub")
	if err != nil {
		log.Fatal().Msgf("%s", err)
	}

	JWT_PUBLIC_KEY, err = jwt.ParseRSAPublicKeyFromPEM(verifyBytes)
	if err != nil {
		log.Fatal().Msgf("%s", err)
	}
}
