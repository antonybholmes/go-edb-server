package main

import (
	"github.com/antonybholmes/go-auth"
	"github.com/antonybholmes/go-edb-server/consts"
	"github.com/antonybholmes/go-edb-server/routes"
	"github.com/golang-jwt/jwt/v5"
)

// you can add your implementation here.
func validateJwtToken(tokenString string) (*jwt.Token, error) {

	token, err := jwt.ParseWithClaims(tokenString, &auth.JwtCustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return consts.JWT_RSA_PUBLIC_KEY, nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, routes.AuthErrorReq("invalid token")
	}

	return token, nil
}
