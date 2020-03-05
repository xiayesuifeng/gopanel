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

func ParseToken(token string) (claims *jwt.StandardClaims, err error) {
	tokenClaims, err := jwt.ParseWithClaims(token, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte("gopanel-secret"), nil
	})

	if tokenClaims != nil {
		if c, ok := tokenClaims.Claims.(*jwt.StandardClaims); ok {
			claims = c
		}
	}

	return claims, err
}
