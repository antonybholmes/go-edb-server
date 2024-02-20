package userdb

import (
	"github.com/rs/zerolog/log"

	"github.com/antonybholmes/go-auth"
)

var USERS *auth.UserDb = nil

func init() {
	var err error

	USERS, err = auth.NewUserDb("data/users.db")

	if err != nil {
		log.Fatal().Msgf("Error loading user db: %s", err)
	}
}

func CreateUser(user *auth.SignupReq) (*auth.AuthUser, error) {
	return USERS.CreateUser(user)
}

func FindUserByEmail(email string) (*auth.AuthUser, error) {
	return USERS.FindUserByEmail(email)
}

func FindUserById(user string) (*auth.AuthUser, error) {
	return USERS.FindUserById(user)
}

func SetIsVerified(user string) error {
	return USERS.SetIsVerified(user)
}

func SetPassword(user string, password string) error {
	return USERS.SetPassword(user, password)
}
