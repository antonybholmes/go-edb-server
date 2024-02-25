package routes

import (
	"net/mail"

	"github.com/antonybholmes/go-auth"
	"github.com/antonybholmes/go-auth/userdb"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
)

//
// Standardized data checkers for checking header and body contain
// the correct data for a route
//

type Validator struct {
	c        echo.Context
	Address  *mail.Address
	Req      *auth.LoginReq
	AuthUser *auth.AuthUser
	Claims   *auth.JwtCustomClaims
	Err      error
}

func NewValidator(c echo.Context) *Validator {
	return &Validator{c: c, Address: nil, Req: nil, AuthUser: nil, Claims: nil, Err: nil}

}

func (validator *Validator) Ok() (*Validator, error) {
	if validator.Err != nil {
		return nil, validator.Err
	} else {
		return validator, nil
	}
}

func (validator *Validator) Success(success func(validator *Validator) error) error {
	if validator.Err != nil {
		return validator.Err
	}

	return success(validator)
}

func (validator *Validator) ReqBind() *Validator {
	if validator.Err != nil {
		return validator
	}

	if validator.Req == nil {
		req := new(auth.LoginReq)

		err := validator.c.Bind(req)

		if err != nil {
			validator.Err = err
		} else {
			validator.Req = req
		}
	}

	return validator
}

func (validator *Validator) ValidEmail() *Validator {
	validator.ReqBind()

	if validator.Err != nil {
		return validator
	}

	address, err := mail.ParseAddress(validator.Req.Username)

	if err != nil {
		validator.Err = InvalidEmailReq()
	} else {
		validator.Address = address
	}

	return validator
}

func (validator *Validator) AuthUserFromEmail() *Validator {
	validator.ValidEmail()

	if validator.Err != nil {
		return validator
	}

	authUser, err := userdb.FindUserByEmail(validator.Address)

	if err != nil {
		validator.Err = UserDoesNotExistReq()
	} else {
		validator.AuthUser = authUser
	}

	return validator

}

func (validator *Validator) AuthUserFromUsername() *Validator {
	validator.ReqBind()

	if validator.Err != nil {
		return validator
	}

	authUser, err := userdb.FindUserByUsername(validator.Req.Username)

	if err != nil {
		validator.Err = UserDoesNotExistReq()
	} else {
		validator.AuthUser = authUser
	}

	return validator

}

func (validator *Validator) IsAuthUser() *Validator {
	if validator.Err != nil {
		return validator
	}

	if validator.AuthUser == nil {
		validator.Err = BadReq("no auth user")
	}

	return validator
}

func (validator *Validator) VerifiedEmail() *Validator {
	validator.IsAuthUser()

	if validator.Err != nil {
		return validator
	}

	if !validator.AuthUser.EmailVerified {
		validator.Err = BadReq("email address not verified")
	}

	return validator
}

func (validator *Validator) JwtClaims() *Validator {
	if validator.Err != nil {
		return validator
	}

	if validator.Claims == nil {
		user := validator.c.Get("user").(*jwt.Token)
		validator.Claims = user.Claims.(*auth.JwtCustomClaims)
	}

	return validator
}

// Extracts uuid from token, checks user exists and calls success function.
// If claims argument is nil, function will search for claims automatically.
// If claims are supplied, this step is skipped. This is so this function can
// be nested in other call backs that may have already extracted the claims
// without having to repeat this part.
func (validator *Validator) AuthUserFromUuid() *Validator {
	validator.JwtClaims()

	if validator.Err != nil {
		return validator
	}

	log.Debug().Msgf("from uuiid %s", validator.Claims.Uuid)

	authUser, err := userdb.FindUserByUuid(validator.Claims.Uuid)

	if err != nil {
		validator.Err = UserDoesNotExistReq()
	} else {
		validator.AuthUser = authUser
	}

	return validator
}

func (validator *Validator) IsValidRefreshToken() *Validator {
	validator.JwtClaims()

	if validator.Err != nil {
		return validator
	}

	if validator.Claims.Type != auth.TOKEN_TYPE_REFRESH {
		validator.Err = BadReq("wrong token type")
	}

	return validator

}

func (validator *Validator) IsValidAccessToken() *Validator {
	validator.JwtClaims()

	if validator.Err != nil {
		return validator
	}

	if validator.Claims.Type != auth.TOKEN_TYPE_ACCESS {
		validator.Err = BadReq("wrong token type")
	}

	return validator
}
