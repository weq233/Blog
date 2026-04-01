package utils

import (
	"errors"
	"time"
	"github.com/golang-jwt/jwt/v5"
)

// JWT 密钥（实际项目中应该从配置文件读取）
var jwtSecret = []byte("your-secret-key-change-in-production")

// Claims JWT 声明
type Claims struct {
	UserId   int    `json:"user_id"`
	Username string `json:"username"`
	Role     int    `json:"role"`
	jwt.RegisteredClaims
}

// GenerateToken 生成 JWT 令牌
func GenerateToken(userId int, username string, role int) (string, error) {
	claims := Claims{
		UserId:   userId,
		Username: username,
		Role:     role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)), // 24 小时有效期
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "blog-system",
		},
	}
	
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

// ParseToken 解析 JWT 令牌
func ParseToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})
	
	if err != nil {
		return nil, err
	}
	
	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}
	
	return nil, errors.New("invalid token")
}

// RefreshToken 刷新令牌
func RefreshToken(tokenString string) (string, error) {
	claims, err := ParseToken(tokenString)
	if err != nil {
		return "", err
	}
	
	// 生成新令牌
	return GenerateToken(claims.UserId, claims.Username, claims.Role)
}
