package toolsroutes

import (
	"crypto/rand"
	"math/big"
	"strconv"

	"github.com/antonybholmes/go-auth"
	"github.com/antonybholmes/go-edb-server/routes"
	"github.com/labstack/echo/v4"
)

type HashResp struct {
	Password string `json:"password"`
	Hash     string `json:"hash"`
}

func HashedPasswordRoute(c echo.Context) error {

	password := c.QueryParam("password")

	if len(password) == 0 {
		return routes.ErrorReq("password cannot be empty")
	}

	hash := auth.HashPassword(password)

	ret := HashResp{Password: password, Hash: hash}

	return routes.MakeDataPrettyResp(c, "", ret)
}

type KeyResp struct {
	Key    string `json:"key"`
	Length int    `json:"length"`
}

func RandomKeyRoute(c echo.Context) error {

	l, err := strconv.Atoi(c.QueryParam("l"))

	if err != nil || l < 1 {
		return routes.ErrorReq("length cannot be zero")
	}

	key, err := generateRandomString(l)

	if err != nil {
		return routes.ErrorReq(err)
	}

	ret := KeyResp{Key: key, Length: l}

	return routes.MakeDataPrettyResp(c, "", ret)
}

// func generateRandomString(length int) string {
// 	randomBytes := make([]byte, length)
// 	_, err := rand.Read(randomBytes)
// 	if err != nil {
// 		panic(err) // Handle the error appropriately in your application
// 	}

// 	return base64.URLEncoding.EncodeToString(randomBytes)
// }

const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

// generateRandomString generates a random string of specified length from the letters set.
func generateRandomString(length int) (string, error) {
	b := make([]byte, length)
	for i := range b {
		// Generate a random index
		index, err := rand.Int(rand.Reader, big.NewInt(int64(len(letters))))
		if err != nil {
			return "", err
		}
		b[i] = letters[index.Int64()]
	}
	return string(b), nil
}
