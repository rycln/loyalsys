package auth

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/rycln/loyalsys/internal/models"
)

var ErrInvalidJWT = errors.New("invalid jwt")

const tokenExp = time.Hour * 2

func NewJWTString(userID models.UserID, key string) (string, error) {
	claims := jwt.MapClaims{
		"uid": userID,
		"exp": tokenExp,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(key))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func ParseIDFromJWT(token *jwt.Token) (models.UserID, error) {
	if !token.Valid {
		return 0, ErrInvalidJWT
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return 0, ErrInvalidJWT
	}
	uidFloat, ok := claims["uid"].(float64)
	if !ok {
		return 0, ErrInvalidJWT
	}
	return models.UserID(uidFloat), nil
}
