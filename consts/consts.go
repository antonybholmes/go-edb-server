package consts

import (
	"crypto/rsa"
	"os"

	"github.com/joho/godotenv"
	"github.com/rs/zerolog/log"

	"github.com/antonybholmes/go-sys/env"

	"github.com/golang-jwt/jwt/v5"
)

var NAME string
var APP_NAME string
var APP_URL string
var VERSION string
var COPYRIGHT string
var JWT_RSA_PRIVATE_KEY *rsa.PrivateKey //[]byte
var JWT_RSA_PUBLIC_KEY *rsa.PublicKey   //[]byte
var SESSION_SECRET string
var SESSION_NAME string
var UPDATED string

const DO_NOT_REPLY = "Please do not reply to this message. It was sent from a notification-only email address that we don't monitor."

func init() {
	env.Load()

	godotenv.Load("consts.env")
	godotenv.Load("version.env")

	NAME = os.Getenv("NAME")
	APP_NAME = os.Getenv("APP_NAME")
	APP_URL = os.Getenv("APP_URL")
	VERSION = os.Getenv("VERSION")
	UPDATED = os.Getenv("UPDATED")
	COPYRIGHT = os.Getenv("COPYRIGHT")
	//JWT_PRIVATE_KEY = []byte(os.Getenv("JWT_SECRET"))
	//JWT_PUBLIC_KEY = []byte(os.Getenv("JWT_SECRET"))
	SESSION_SECRET = os.Getenv("SESSION_SECRET")
	SESSION_NAME = os.Getenv("SESSION_NAME")

	bytes, err := os.ReadFile("jwtRS256.key")
	if err != nil {
		log.Fatal().Msgf("%s", err)
	}

	JWT_RSA_PRIVATE_KEY, err = jwt.ParseRSAPrivateKeyFromPEM(bytes)
	if err != nil {
		log.Fatal().Msgf("%s", err)
	}

	bytes, err = os.ReadFile("jwtRS256.key.pub")
	if err != nil {
		log.Fatal().Msgf("%s", err)
	}

	JWT_RSA_PUBLIC_KEY, err = jwt.ParseRSAPublicKeyFromPEM(bytes)
	if err != nil {
		log.Fatal().Msgf("%s", err)
	}
}
