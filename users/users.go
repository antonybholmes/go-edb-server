package users

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

func CreateUser(user *auth.LoginUser, otp string) (*auth.AuthUser, error) {
	return USERS.CreateUser(user, otp)
}

func FindUserByEmail(user *auth.LoginUser) (*auth.AuthUser, error) {
	return USERS.FindUserByEmail(user)
}

func FindUserById(user string) (*auth.AuthUser, error) {
	return USERS.FindUserById(user)
}

func SetIsVerified(user string) error {
	return USERS.SetIsVerified(user)
}
