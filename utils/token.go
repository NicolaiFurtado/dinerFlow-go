package utils

import (
	"errors"

	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	UserID uint `json:"userId"`
	jwt.RegisteredClaims
}

var JwtKey = []byte("suaChaveSuperSecreta")

func ValidateToken(tokenString string) (*Claims, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return JwtKey, nil
	})
	if err != nil || !token.Valid {
		return nil, errors.New("token inválido")
	}
	return claims, nil
}
