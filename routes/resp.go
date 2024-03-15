package routes

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type StatusResp struct {
	Status int `json:"status"`
}

type StatusMessageResp struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

type DataResp struct {
	StatusMessageResp
	Data interface{} `json:"data"`
}

type SuccessResp struct {
	Success bool `json:"success"`
}

type ValidResp struct {
	Valid bool `json:"valid"`
}

type JwtResp struct {
	Jwt string `json:"jwt"`
}

type RefreshTokenResp struct {
	RefreshToken string `json:"refreshToken"`
}

type AccessTokenResp struct {
	AccessToken string `json:"accessToken"`
}

type LoginResp struct {
	RefreshToken string `json:"refreshToken"`
	AccessToken  string `json:"accessToken"`
}

func JsonResp[V any](c echo.Context, status int, data V) error {
	return c.JSONPretty(status, data, " ")
}

// func MakeBadResp(c echo.Context, err error) error {
// 	return JsonRep(c, http.StatusBadRequest, StatusResp{StatusResp: StatusResp{Status: http.StatusBadRequest}, Message: err.Error()})
// }

func MakeDataResp[V any](c echo.Context, message string, data V) error {
	return JsonResp(c, http.StatusOK, DataResp{StatusMessageResp: StatusMessageResp{Status: http.StatusOK, Message: message}, Data: data})
}

// func MakeValidResp(c echo.Context, message string, valid bool) error {
// 	return MakeDataResp(c, message, &ValidResp{Valid: valid})
// }

func MakeOkResp(c echo.Context, message string) error {
	return MakeSuccessResp(c, message, true)
}

func MakeSuccessResp(c echo.Context, message string, success bool) error {
	return MakeDataResp(c, message, &SuccessResp{Success: success})
}

func PasswordUpdatedResp(c echo.Context) error {
	return MakeOkResp(c, "password updated")
}
