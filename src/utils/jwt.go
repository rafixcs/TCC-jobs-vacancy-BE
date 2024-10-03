package utils

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type UserAuthClaims struct {
	UserId  string `json:"user_id"`
	LoginId string `json:"login_id"`
	jwt.RegisteredClaims
}

const JWT_SECRET = "secretpass"

func ParseToken(tokenStr string) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("invalid token")
		}

		return []byte(JWT_SECRET), nil
	})

	return token, err
}

func CreateUserJwtToken(userId, loginId string) (string, error) {
	claims := UserAuthClaims{
		UserId:  userId,
		LoginId: loginId,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "test",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ss, err := token.SignedString([]byte(JWT_SECRET))
	return ss, err
}
