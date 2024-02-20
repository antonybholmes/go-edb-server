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

func CreateUser(user *auth.SignupReq, otp string) (*auth.AuthUser, error) {
	return USERS.CreateUser(user, otp)
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

func SetOtp(user string, code string) error {
	return USERS.SetOtp(user, code)
}
