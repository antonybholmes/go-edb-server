package authentication

import (
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/antonybholmes/go-auth"
	"github.com/antonybholmes/go-auth/tokengen"
	"github.com/antonybholmes/go-auth/userdbcache"
	"github.com/antonybholmes/go-edb-server/consts"
	"github.com/antonybholmes/go-edb-server/routes"
	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
)

const (
	SESSION_PUBLICID string = "publicId"
	SESSION_ROLES    string = "roles"
)

var SESSION_OPT_ZERO *sessions.Options

//var SESSION_OPT_24H *sessions.Options
//var SESSION_OPT_30_DAYS *sessions.Options
//var SESSION_OPT_7_DAYS *sessions.Options

func init() {

	// HttpOnly and Secure are disabled so we can use them
	// cross domain for testing
	// http only false to allow js to delete etc on the client side

	// For sessions that should end when browser closes
	SESSION_OPT_ZERO = &sessions.Options{
		Path:     "/",
		MaxAge:   0,
		HttpOnly: false,
		Secure:   true,
		SameSite: http.SameSiteNoneMode,
	}

	// SESSION_OPT_24H = &sessions.Options{
	// 	Path:     "/",
	// 	MaxAge:   auth.MAX_AGE_DAY_SECS,
	// 	HttpOnly: false,
	// 	Secure:   true,
	// 	SameSite: http.SameSiteNoneMode,
	// }

	// SESSION_OPT_30_DAYS = &sessions.Options{
	// 	Path:     "/",
	// 	MaxAge:   auth.MAX_AGE_30_DAYS_SECS,
	// 	HttpOnly: false,
	// 	Secure:   true,
	// 	SameSite: http.SameSiteNoneMode,
	// }

	// SESSION_OPT_7_DAYS = &sessions.Options{
	// 	Path:     "/",
	// 	MaxAge:   auth.MAX_AGE_7_DAYS_SECS,
	// 	HttpOnly: false,
	// 	Secure:   true,
	// 	SameSite: http.SameSiteNoneMode,
	// }
}

type SessionRoutes struct {
	options *sessions.Options
}

func NewSessionRoutes() *SessionRoutes {
	maxAge := auth.MAX_AGE_7_DAYS_SECS

	t := os.Getenv("SESSION_TTL_HOURS")

	if t != "" {
		v, err := strconv.ParseUint(t, 10, 32)

		if err == nil {
			maxAge = int((time.Duration(v) * time.Hour).Seconds())
		}
	}

	options := sessions.Options{
		Path:     "/",
		MaxAge:   maxAge,
		HttpOnly: false,
		Secure:   true,
		SameSite: http.SameSiteNoneMode,
	}

	return &SessionRoutes{options: &options}
}

func (sr *SessionRoutes) SessionUsernamePasswordSignInRoute(c echo.Context) error {
	validator, err := NewValidator(c).ParseLoginRequestBody().Ok()

	if err != nil {
		return err
	}

	if validator.Req.Password == "" {
		return PasswordlessSigninEmailRoute(c, validator)
	}

	user := validator.Req.Username

	authUser, err := userdbcache.FindUserByUsername(user)

	// if err != nil {
	// 	//
	// 	email, err := mail.ParseAddress(user)

	// 	if err != nil {
	// 		return routes.InvalidEmailReq()
	// 	}

	// 	// also check if username is valid email and try to login
	// 	// with that
	// 	authUser, err = userdbcache.FindUserByEmail(email)

	// 	if err != nil {
	// 		return routes.UserDoesNotExistReq()
	// 	}
	// }

	if err != nil {
		return routes.UserDoesNotExistReq()
	}

	if !authUser.EmailIsVerified {
		return routes.EmailNotVerifiedReq()
	}

	roles, err := userdbcache.UserRoleList(authUser)

	if err != nil {
		return routes.AuthErrorReq("could not get user roles")
	}

	roleClaim := auth.MakeClaim(roles)

	if !auth.CanLogin(roleClaim) {
		return routes.UserNotAllowedToSignIn()
	}

	err = authUser.CheckPasswordsMatch(validator.Req.Password)

	if err != nil {
		return routes.ErrorReq(err)
	}

	sess, err := session.Get(consts.SESSION_NAME, c)
	if err != nil {
		return routes.ErrorReq("error creating session")
	}

	if validator.Req.StaySignedIn {
		sess.Options = sr.options
	} else {
		sess.Options = SESSION_OPT_ZERO
	}

	sess.Values[SESSION_PUBLICID] = authUser.PublicId
	sess.Values[SESSION_ROLES] = roleClaim //auth.MakeClaim(authUser.Roles)

	sess.Save(c.Request(), c.Response())

	return UserSignedInResp(c)
	//return c.NoContent(http.StatusOK)
}

// Validate the passwordless token we generated and create
// a user session. The session acts as a refresh token and
// can be used to generate access tokens to use resources
func (sr *SessionRoutes) SessionPasswordlessValidateSignInRoute(c echo.Context) error {

	return NewValidator(c).LoadAuthUserFromToken().CheckUserHasVerifiedEmailAddress().Success(func(validator *Validator) error {

		if validator.Claims.Type != auth.PASSWORDLESS_TOKEN {
			return routes.WrongTokentTypeReq()
		}

		authUser := validator.AuthUser

		roles, err := userdbcache.UserRoleList(authUser)

		if err != nil {
			return routes.ErrorReq("error getting user roles")
		}

		roleClaim := auth.MakeClaim(roles)

		//log.Debug().Msgf("user %v", authUser)

		if !auth.CanLogin(roleClaim) {
			return routes.UserNotAllowedToSignIn()
		}

		sess, err := session.Get(consts.SESSION_NAME, c)

		if err != nil {
			return routes.ErrorReq("error creating session")
		}

		// set session options such as if cookie secure and how long it
		// persists
		sess.Options = sr.options //SESSION_OPT_30_DAYS
		sess.Values[SESSION_PUBLICID] = authUser.PublicId
		sess.Values[SESSION_ROLES] = roleClaim //auth.MakeClaim(authUser.Roles)

		sess.Save(c.Request(), c.Response())

		return UserSignedInResp(c)
	})
}

func SessionSignOutRoute(c echo.Context) error {
	sess, err := session.Get(consts.SESSION_NAME, c)
	if err != nil {
		return routes.ErrorReq("error creating session")
	}

	// invalidate by time
	sess.Values[SESSION_PUBLICID] = ""
	sess.Values[SESSION_ROLES] = ""
	sess.Options.MaxAge = 0

	sess.Save(c.Request(), c.Response())

	return routes.MakeOkPrettyResp(c, "user was signed out")
}

func SessionNewAccessTokenRoute(c echo.Context) error {
	sess, _ := session.Get(consts.SESSION_NAME, c)
	publicId, _ := sess.Values[SESSION_PUBLICID].(string)
	roles, _ := sess.Values[SESSION_ROLES].(string)

	t, err := tokengen.AccessToken(c, publicId, roles)

	if err != nil {
		return routes.TokenErrorReq()
	}

	return routes.MakeDataPrettyResp(c, "", &routes.AccessTokenResp{AccessToken: t})
}

func SessionUserRoute(c echo.Context) error {
	sess, _ := session.Get(consts.SESSION_NAME, c)
	publicId, _ := sess.Values[SESSION_PUBLICID].(string)

	authUser, err := userdbcache.FindUserByPublicId(publicId)

	if err != nil {
		return routes.UserDoesNotExistReq()
	}

	return routes.MakeDataPrettyResp(c, "", *authUser)
}

// // Start passwordless login by sending an email
// func SessionSendResetPasswordEmailRoute(c echo.Context) error {

// 	return NewValidator(c).LoadAuthUserFromSession().CheckUserHasVerifiedEmailAddress().Success(func(validator *Validator) error {

// 		authUser := validator.AuthUser
// 		req := validator.Req

// 		otpJwt, err := tokengen.ResetPasswordToken(c, authUser)

// 		if err != nil {
// 			return routes.ErrorReq(err)
// 		}

// 		var file string

// 		if req.CallbackUrl != "" {
// 			file = "templates/email/password/reset/web.html"
// 		} else {
// 			file = "templates/email/password/reset/api.html"
// 		}

// 		go SendEmailWithToken("Password Reset",
// 			authUser,
// 			file,
// 			otpJwt,
// 			req.CallbackUrl,
// 			req.VisitUrl)

// 		//if err != nil {
// 		//	return routes.ErrorReq(err)
// 		//}

// 		return routes.MakeOkPrettyResp(c, "check your email for a password reset link")
// 	})
// }

// func SessionUpdatePasswordRoute(c echo.Context) error {

// 	return NewValidator(c).LoadAuthUserFromSession().ParseLoginRequestBody().Success(func(validator *Validator) error {

// 		err := userdbcache.SetPassword(validator.AuthUser.PublicId, validator.Req.Password)

// 		if err != nil {
// 			return routes.ErrorReq(err)
// 		}

// 		return SendEmailChangedEmail(c, validator.AuthUser, validator.Req.Password)
// 	})
// }

// func SessionSendResetEmailEmailRoute(c echo.Context) error {

// 	return NewValidator(c).LoadAuthUserFromSession().ParseLoginRequestBody().Success(func(validator *Validator) error {

// 		req := validator.Req

// 		newEmail, err := mail.ParseAddress(req.Email)

// 		if err != nil {
// 			return routes.ErrorReq(err)
// 		}

// 		otpJwt, err := tokengen.ResetEmailToken(c, validator.AuthUser, newEmail)

// 		if err != nil {
// 			return routes.ErrorReq(err)
// 		}

// 		file := "templates/email/email/reset/web.html"

// 		go BaseSendEmailWithToken("Change Your Email Address",
// 			validator.AuthUser,
// 			newEmail,
// 			file,
// 			otpJwt,
// 			req.CallbackUrl,
// 			req.VisitUrl)

// 		//if err != nil {
// 		//	return routes.ErrorReq(err)
// 		//}

// 		return routes.MakeOkPrettyResp(c, "check your email for a change email link")
// 	})
// }
