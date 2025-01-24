package routes

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type StatusResp struct {
	Status int `json:"status"`
}

type StatusMessageResp struct {
	Message string `json:"message"`
	Status  int    `json:"status"`
}

type DataResp struct {
	Data interface{} `json:"data"`
	StatusMessageResp
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

func JsonResp[V any](c echo.Context, status int, data V, pretty bool) error {

	if pretty {
		return c.JSONPretty(status, data, " ")
	} else {
		return c.JSON(status, data)
	}
}

// func MakeBadResp(c echo.Context, err error) error {
// 	return JsonRep(c, http.StatusBadRequest, StatusResp{StatusResp: StatusResp{Status: http.StatusBadRequest}, Message: err.Error()})
// }

func MakeDataPrettyResp[V any](c echo.Context, message string, data V) error {
	return MakeDataResp(c, message, data, true)
}

func MakeDataResp[V any](c echo.Context, message string, data V, pretty bool) error {
	return JsonResp(c,
		http.StatusOK,
		DataResp{
			StatusMessageResp: StatusMessageResp{
				Status:  http.StatusOK,
				Message: message,
			},
			Data: data,
		},
		pretty)
}

// func MakeValidResp(c echo.Context, message string, valid bool) error {
// 	return MakeDataResp(c, message, &ValidResp{Valid: valid})
// }

func MakeOkPrettyResp(c echo.Context, message string) error {
	return MakeOkResp(c, message, true)
}

func MakeOkResp(c echo.Context, message string, pretty bool) error {
	return MakeSuccessResp(c, message, true, pretty)
}

func MakeSuccessPrettyResp(c echo.Context, message string, success bool) error {
	return MakeSuccessResp(c, message, success, true)
}

func MakeSuccessResp(c echo.Context, message string, success bool, pretty bool) error {
	return MakeDataResp(c, message, &SuccessResp{Success: success}, pretty)
}
