package auth

import (
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"

	"github.com/gofiber/fiber/v2"

	"github.com/golang-jwt/jwt/v5"
)

const tokenExp = time.Hour * 2

var ErrNoJWT = errors.New("no jwt token")

type JWT struct {
	key   string
	token *jwt.Token
}

func NewJWT(userID, key string) (*JWT, error) {
	claims := jwtClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(tokenExp)),
		},
		ID: userID,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return tokenString, nil
}

type jwtClaims struct {
	jwt.RegisteredClaims
	ID string `json:"id"`
}

func GetJWT(c *fiber.Ctx, key string) (string, error) {
	rawToken := string(c.Request().Header.Peek("Authorization"))
	if rawToken == "" {
		return "", "", ErrNoToken
	}
	rawToken = strings.TrimPrefix(rawToken, "Bearer")
	rawToken = strings.TrimSpace(rawToken)
	claims := &jwtClaims{}
	_, err := jwt.ParseWithClaims(rawToken, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(key), nil
	})
	if err != nil {
		return "", "", err
	}
	return rawToken, claims.ID, nil
}

func makeUserID() string {
	return uuid.NewString()
}

func makeTokenString(uid, key string) (string, error) {
	claims := jwtClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(tokenExp)),
		},
		ID: uid,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(key))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

tokenString, err := token.SignedString([]byte(key))
if err != nil {
	return "", err
}