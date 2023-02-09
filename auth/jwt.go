package auth

import (
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"gitlab.com/xiayesuifeng/gopanel/core/config"
	"math/rand"
	"time"
)

func GenerateToken() (string, error) {
	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24 * 3)),
	})

	if config.Conf.Secret == "" {
		rand.Seed(time.Now().Unix())
		config.Conf.Secret = fmt.Sprintf("gopanel-secret-%.6d", rand.Intn(999999))
	}
	return tokenClaims.SignedString([]byte(config.Conf.Secret))
}

func ParseToken(token string) (claims *jwt.RegisteredClaims, err error) {
	tokenClaims, err := jwt.ParseWithClaims(token, &jwt.RegisteredClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.Conf.Secret), nil
	})

	if tokenClaims != nil {
		if c, ok := tokenClaims.Claims.(*jwt.RegisteredClaims); ok {
			claims = c
		}
	}

	return claims, err
}
