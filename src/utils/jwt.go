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

func GetUserIdFromToken(tokenHeader string) (string, error) {
	token, err := ParseToken(tokenHeader)
	if err != nil {
		return "", fmt.Errorf("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", fmt.Errorf("invalid token claims")
	}

	userId, ok := claims["user_id"].(string)
	if !ok {
		return "", fmt.Errorf("token missing user field")
	}

	return userId, nil
}

func GetUserAuthIdsFromToken(tokenHeader string) (string, string, error) {
	token, err := ParseToken(tokenHeader)
	if err != nil {
		return "", "", fmt.Errorf("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", "", fmt.Errorf("invalid token claims")
	}

	loginId, ok := claims["login_id"].(string)
	if !ok {
		return "", "", fmt.Errorf("token missing login field")
	}

	userId, ok := claims["user_id"].(string)
	if !ok {
		return "", "", fmt.Errorf("token missing user field")
	}

	return userId, loginId, nil
}
