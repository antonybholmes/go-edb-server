package authenticationroutes

import (
	"net/mail"

	"github.com/antonybholmes/go-auth"
	"github.com/antonybholmes/go-auth/userdbcache"
	"github.com/antonybholmes/go-edb-server/routes"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

//
// Standardized data checkers for checking header and body contain
// the correct data for a route
//

type Validator struct {
	c            echo.Context
	Address      *mail.Address
	LoginBodyReq *auth.LoginBodyReq

	AuthUser *auth.AuthUser
	Claims   *auth.TokenClaims
	Err      *echo.HTTPError
}

func NewValidator(c echo.Context) *Validator {
	return &Validator{c: c, Address: nil, LoginBodyReq: nil, AuthUser: nil, Claims: nil, Err: nil}

}

func (validator *Validator) Ok() (*Validator, error) {
	if validator.Err != nil {
		return nil, validator.Err
	} else {
		return validator, nil
	}
}

// If the validator does not encounter errors, it will run the success function
// allowing you to extract data from the validator, otherwise it returns an error
// without running the function
func (validator *Validator) Success(success func(validator *Validator) error) error {

	if validator.Err != nil {
		return validator.Err
	}

	return success(validator)
}

func (validator *Validator) ParseLoginRequestBody() *Validator {
	if validator.Err != nil {
		return validator
	}

	if validator.LoginBodyReq == nil {
		var req auth.LoginBodyReq

		err := validator.c.Bind(&req)

		if err != nil {
			validator.Err = routes.ErrorReq(err)
		} else {
			validator.LoginBodyReq = &req
		}
	}

	return validator
}

func (validator *Validator) CheckUsernameIsWellFormed() *Validator {
	validator.ParseLoginRequestBody()

	if validator.Err != nil {
		return validator
	}

	//address, err := auth.CheckEmailIsWellFormed(validator.Req.Email)

	err := auth.CheckUsername(validator.LoginBodyReq.Username)

	if err != nil {
		validator.Err = routes.ErrorReq(err)
	}

	return validator
}

func (validator *Validator) CheckEmailIsWellFormed() *Validator {
	validator.ParseLoginRequestBody()

	if validator.Err != nil {
		return validator
	}

	//address, err := auth.CheckEmailIsWellFormed(validator.Req.Email)

	address, err := mail.ParseAddress(validator.LoginBodyReq.Email)

	if err != nil {
		validator.Err = routes.ErrorReq(err)
	} else {
		validator.Address = address
	}

	return validator
}

func (validator *Validator) LoadAuthUserFromUuid() *Validator {

	if validator.Err != nil {
		return validator
	}

	authUser, err := userdbcache.FindUserByUuid(validator.LoginBodyReq.Uuid)

	if err != nil {
		validator.Err = routes.UserDoesNotExistReq()
	} else {
		validator.AuthUser = authUser
	}

	return validator

}

func (validator *Validator) LoadAuthUserFromEmail() *Validator {
	validator.CheckEmailIsWellFormed()

	if validator.Err != nil {
		return validator
	}

	authUser, err := userdbcache.FindUserByEmail(validator.Address)

	if err != nil {
		validator.Err = routes.UserDoesNotExistReq()
	} else {
		validator.AuthUser = authUser
	}

	return validator

}

func (validator *Validator) LoadAuthUserFromUsername() *Validator {
	validator.ParseLoginRequestBody()

	if validator.Err != nil {
		return validator
	}

	authUser, err := userdbcache.FindUserByUsername(validator.LoginBodyReq.Username)

	//log.Debug().Msgf("beep2 %s", authUser.Username)

	if err != nil {
		validator.Err = routes.UserDoesNotExistReq()
	} else {
		validator.AuthUser = authUser
	}

	return validator

}

func (validator *Validator) LoadAuthUserFromSession() *Validator {
	validator.ParseLoginRequestBody()

	if validator.Err != nil {
		return validator
	}

	sessionData, err := ReadSessionInfo(validator.c)

	if err != nil {
		validator.Err = routes.ErrorReq("user not in session")
		validator.CheckIsValidRefreshToken().CheckUsernameIsWellFormed()
	}

	validator.AuthUser = sessionData.AuthUser

	return validator
}

func (validator *Validator) CheckAuthUserIsLoaded() *Validator {
	if validator.Err != nil {
		return validator
	}

	if validator.AuthUser == nil {
		validator.Err = routes.ErrorReq("no auth user")
	}

	return validator
}

func (validator *Validator) CheckUserHasVerifiedEmailAddress() *Validator {
	validator.CheckAuthUserIsLoaded()

	if validator.Err != nil {
		return validator
	}

	if validator.AuthUser.EmailVerifiedAt == 0 {
		validator.Err = routes.ErrorReq("email address not verified")
	}

	return validator
}

// If using jwt middleware, token is put into user variable
// and we can extract data from the jwt
func (validator *Validator) LoadTokenClaims() *Validator {
	if validator.Err != nil {
		return validator
	}

	if validator.Claims == nil {
		user := validator.c.Get("user").(*jwt.Token)
		validator.Claims = user.Claims.(*auth.TokenClaims)
	}

	return validator
}

// Extracts public id from token, checks user exists and calls success function.
// If claims argument is nil, function will search for claims automatically.
// If claims are supplied, this step is skipped. This is so this function can
// be nested in other call backs that may have already extracted the claims
// without having to repeat this part.
func (validator *Validator) LoadAuthUserFromToken() *Validator {
	validator.LoadTokenClaims()

	if validator.Err != nil {
		return validator
	}

	authUser, err := userdbcache.FindUserByUuid(validator.Claims.Uuid)

	if err != nil {
		validator.Err = routes.UserDoesNotExistReq()
	} else {
		validator.AuthUser = authUser
	}

	return validator
}

func (validator *Validator) CheckIsValidRefreshToken() *Validator {
	validator.LoadTokenClaims()

	if validator.Err != nil {
		return validator
	}

	if validator.Claims.Type != auth.REFRESH_TOKEN {
		validator.Err = routes.ErrorReq("no refresh token")
	}

	return validator

}

func (validator *Validator) CheckIsValidAccessToken() *Validator {
	validator.LoadTokenClaims()

	if validator.Err != nil {
		return validator
	}

	if validator.Claims.Type != auth.ACCESS_TOKEN {
		validator.Err = routes.ErrorReq("no access token")
	}

	return validator
}
