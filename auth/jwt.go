package auth

import (
	"github.com/dgrijalva/jwt-go"
	"time"
)

func GenerateToken() (string, error) {
	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		ExpiresAt: time.Now().Add(time.Hour * 24 * 3).Unix(),
	})

	return tokenClaims.SignedString([]byte("gopanel-secret"))
}
