package pkg

import (
	"errors"
	"github.com/fynn404/gin-demo/backend/internal/common/config"
	"github.com/golang-jwt/jwt/v5"
	"time"
)

// 错误定义
var (
	ErrInvalidToken = errors.New("invalid token")
	ErrExpiredToken = errors.New("token has expired")
)

type GenerateTokenReq struct {
	UserID   string `json:"user_id"`
	Role     string `json:"role"`
	UserName string `json:"username"`
}

// Claims 自定义的JWT声明结构
type Claims struct {
	UserID   string `json:"user_id"`
	Role     string `json:"role"`
	UserName string `json:"username"`
	jwt.RegisteredClaims
}

// GenerateToken 生成JWT令牌
func GenerateToken(req *GenerateTokenReq) (string, error) {
	var (
		expHours       = 24 // 默认24小时
		configExpHours = config.GlobalConfig.JWT.ExpireHours
	)
	if configExpHours != 0 {
		expHours = configExpHours
	}

	claims := Claims{
		UserID:   req.UserID,
		Role:     req.Role,
		UserName: req.UserName,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(expHours) * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(config.GlobalConfig.JWT.Secret))
}

// ParseToken 解析JWT令牌
func ParseToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.GlobalConfig.JWT.Secret), nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, ErrExpiredToken
		}
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, ErrInvalidToken
}

// ValidateToken 验证JWT令牌并返回解析后的声明
func ValidateToken(tokenString string) (*Claims, error) {
	claims, err := ParseToken(tokenString)
	if err != nil {
		return nil, err
	}

	return claims, nil
}
