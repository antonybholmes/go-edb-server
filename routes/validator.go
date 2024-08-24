package routes

import (
	"net/mail"

	"github.com/antonybholmes/go-auth"
	"github.com/antonybholmes/go-auth/userdbcache"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
)

const SESSION_NAME string = "session"
const SESSION_UUID string = "uuid"

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
	Err      *echo.HTTPError
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

	if validator.Req == nil {
		var req auth.LoginReq

		err := validator.c.Bind(&req)

		if err != nil {
			validator.Err = ErrorReq(err)
		} else {
			validator.Req = &req
		}
	}

	return validator
}

func (validator *Validator) CheckEmailIsWellFormed() *Validator {
	validator.ParseLoginRequestBody()

	if validator.Err != nil {
		return validator
	}

	address, err := auth.CheckEmailIsWellFormed(validator.Req.Username)

	if err != nil {
		validator.Err = ErrorReq(err)
	} else {
		validator.Address = address
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
		validator.Err = UserDoesNotExistReq()
	} else {
		validator.AuthUser = authUser
	}

	return validator

}

func (validator *Validator) LoadAuthUserFromId() *Validator {
	validator.ParseLoginRequestBody()

	if validator.Err != nil {
		return validator
	}

	authUser, err := userdbcache.FindUserById(validator.Req.Username)

	if err != nil {
		validator.Err = UserDoesNotExistReq()
	} else {
		validator.AuthUser = authUser
	}

	return validator

}

func (validator *Validator) LoadAuthUserFromSession() *Validator {
	sess, _ := session.Get(SESSION_NAME, validator.c)
	uuid, _ := sess.Values[SESSION_UUID].(string)

	if validator.Err != nil {
		return validator
	}

	authUser, err := userdbcache.FindUserByUuid(uuid)

	if err != nil {
		validator.Err = UserDoesNotExistReq()
	} else {
		validator.AuthUser = authUser
	}

	return validator
}

func (validator *Validator) CheckAuthUserIsLoaded() *Validator {
	if validator.Err != nil {
		return validator
	}

	if validator.AuthUser == nil {
		validator.Err = ErrorReq("no auth user")
	}

	return validator
}

func (validator *Validator) CheckUserHasVerifiedEmailAddress() *Validator {
	validator.CheckAuthUserIsLoaded()

	if validator.Err != nil {
		return validator
	}

	if !validator.AuthUser.EmailVerified {
		validator.Err = ErrorReq("email address not verified")
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
		validator.Claims = user.Claims.(*auth.JwtCustomClaims)
	}

	return validator
}

// Extracts uuid from token, checks user exists and calls success function.
// If claims argument is nil, function will search for claims automatically.
// If claims are supplied, this step is skipped. This is so this function can
// be nested in other call backs that may have already extracted the claims
// without having to repeat this part.
func (validator *Validator) LoadAuthUserFromToken() *Validator {
	validator.LoadTokenClaims()

	if validator.Err != nil {
		return validator
	}

	//log.Debug().Msgf("from uuid %s", validator.Claims.Uuid)

	authUser, err := userdbcache.FindUserByUuid(validator.Claims.Uuid)

	if err != nil {
		validator.Err = UserDoesNotExistReq()
	} else {
		validator.AuthUser = authUser
	}

	return validator
}

func (validator *Validator) CheckIsValidRefreshToken() *Validator {
	validator.LoadAuthUserFromToken()

	if validator.Err != nil {
		return validator
	}

	if validator.Claims.Type != auth.TOKEN_TYPE_REFRESH {
		validator.Err = ErrorReq("no refresh token")
	}

	return validator

}

func (validator *Validator) CheckIsValidAccessToken() *Validator {
	validator.LoadAuthUserFromToken()

	if validator.Err != nil {
		return validator
	}

	if validator.Claims.Type != auth.TOKEN_TYPE_ACCESS {
		validator.Err = ErrorReq("no access token")
	}

	return validator
}
