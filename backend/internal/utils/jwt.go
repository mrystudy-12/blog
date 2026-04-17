package utils

import (
	"GoWork_9/backend/internal/model"
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"time"
)

var jwtKey = []byte("my_secret_key")

// GenerateToken 生成 JWT Token，按照用户指定格式
func GenerateToken(id uint64, username, role string) (string, error) {
	expirationTime := time.Now().Add(30 * time.Minute)
	claims := model.Claims{
		ID:       id,
		Username: username,
		Role:     role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "myapp",
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtKey)
}

// ParseToken 解析 JWT Token
func ParseToken(tokenString string) (*model.Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &model.Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return jwtKey, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*model.Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}
