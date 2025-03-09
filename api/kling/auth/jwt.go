package auth

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWTClaims struct {
	jwt.RegisteredClaims
}

func GenerateToken(ak, sk string) (string, error) {
	now := time.Now()
	claims := jwt.RegisteredClaims{
		Issuer:    ak,
		ExpiresAt: jwt.NewNumericDate(now.Add(1800 * time.Second)), // 30分钟有效期
		NotBefore: jwt.NewNumericDate(now.Add(-5 * time.Second)),   // 立即生效
		IssuedAt:  jwt.NewNumericDate(now),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(sk))
}

func ValidateToken(tokenString, sk string) (*jwt.RegisteredClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &jwt.RegisteredClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(sk), nil
	})

	if claims, ok := token.Claims.(*jwt.RegisteredClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, err
}
