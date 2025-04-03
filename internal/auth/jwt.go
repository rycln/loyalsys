package auth

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/rycln/loyalsys/internal/models"
)

var ErrInvalidJWT = errors.New("invalid jwt")

const tokenExp = time.Hour * 2

type JwtClaims struct {
	jwt.RegisteredClaims
	UserID models.UserID `json:"id"`
}

func (claims *JwtClaims) GetUserID() models.UserID {
	return claims.UserID
}

func NewJWTString(userID models.UserID, key string) (string, error) {
	claims := JwtClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(tokenExp)),
		},
		UserID: userID,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(key))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func ParseJWT(tokenString, key string) (*JwtClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JwtClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(key), nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(*JwtClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, ErrInvalidJWT
}
