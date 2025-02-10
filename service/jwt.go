package service

import (
	"context"
	"errors"
	"fmt"
	"nexus-ai/constant"
	"nexus-ai/redis"
	"time"

	"github.com/dgrijalva/jwt-go"
)

var jwtSecret = []byte(constant.JwtSecret)

var (
	ErrInvalidToken     = errors.New("invalid token")
	ErrExpiredToken     = errors.New("token has expired")
	ErrTokenBlacklisted = errors.New("token has been blacklisted")
)

const (
	AccessTokenExpiry  = constant.JwtAccessTokenExpiry
	RefreshTokenExpiry = constant.JwtRefreshTokenExpiry
	// token黑名单的Redis key前缀
	tokenBlacklistPrefix = "token:blacklist:"
	// 用户黑名单的Redis key前缀
	userBlacklistPrefix = "user:blacklist:"
)

type TokenType string

const (
	AccessToken  TokenType = "access"
	RefreshToken TokenType = "refresh"
)

// GenerateToken 生成指定类型的token
func GenerateToken(userID string, tokenType TokenType) (string, error) {
	expiry := AccessTokenExpiry
	if tokenType == RefreshToken {
		expiry = RefreshTokenExpiry
	}

	claims := jwt.MapClaims{
		"user_id":    userID,
		"token_type": string(tokenType),
		"exp":        time.Now().Add(expiry).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

// ValidateToken 验证token并返回解析后的token对象，同时验证userID是否匹配
func ValidateToken(tokenString string, expectedUserID string) (*jwt.Token, error) {
	// 检查token是否在黑名单中
	ctx := context.Background()
	isBlacklisted, err := redis.Exists(ctx, tokenBlacklistPrefix+tokenString)
	if err != nil {
		return nil, fmt.Errorf("failed to check token blacklist: %v", err)
	}
	if isBlacklisted {
		return nil, ErrTokenBlacklisted
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrInvalidToken
		}
		return jwtSecret, nil
	})

	if err != nil {
		if ve, ok := err.(*jwt.ValidationError); ok {
			if ve.Errors&jwt.ValidationErrorExpired != 0 {
				return nil, ErrExpiredToken
			}
		}
		return nil, err
	}

	// 检查用户级别的黑名单
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, ErrInvalidToken
	}

	// 验证用户ID是否匹配
	tokenUserID, ok := claims["user_id"].(string)
	if !ok || tokenUserID != expectedUserID {
		return nil, errors.New("user id not match")
	}

	if userID, ok := claims["user_id"].(string); ok {
		isUserBlacklisted, err := redis.Exists(ctx, userBlacklistPrefix+userID)
		if err != nil {
			return nil, fmt.Errorf("failed to check user blacklist: %v", err)
		}
		if isUserBlacklisted {
			return nil, ErrTokenBlacklisted
		}
	}

	return token, nil
}

// RefreshTokenPair 包含新生成的访问令牌和刷新令牌
type TokenPair struct {
	AccessToken  string
	RefreshToken string
}

// RefreshAccessToken 使用refresh token置换新的access token和refresh token对
func RefreshAccessToken(refreshToken string, expectedUserID string) (*TokenPair, error) {
	// 验证refresh token
	token, err := ValidateToken(refreshToken, expectedUserID)
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, ErrInvalidToken
	}

	// 验证token类型
	tokenType, ok := claims["token_type"].(string)
	if !ok || TokenType(tokenType) != RefreshToken {
		return nil, errors.New("invalid token type: not a refresh token")
	}

	// 获取用户ID
	userID, ok := claims["user_id"].(string)
	if !ok {
		return nil, errors.New("invalid token: missing user_id")
	}

	// 生成新的token对
	newAccessToken, err := GenerateToken(userID, AccessToken)
	if err != nil {
		return nil, err
	}

	newRefreshToken, err := GenerateToken(userID, RefreshToken)
	if err != nil {
		return nil, err
	}

	// 将旧的refresh token加入黑名单
	if err := InvalidateToken(refreshToken); err != nil {
		return nil, err
	}

	return &TokenPair{
		AccessToken:  newAccessToken,
		RefreshToken: newRefreshToken,
	}, nil
}

// InvalidateAllUserTokens 使指定用户的所有token失效
func InvalidateAllUserTokens(userID string) error {
	ctx := context.Background()
	// 使用最长的token过期时间（refresh token的过期时间）
	return redis.Set(ctx, userBlacklistPrefix+userID, time.Now().Unix(), RefreshTokenExpiry)
}

// IsTokenBlacklisted 检查token是否在黑名单中
func IsTokenBlacklisted(tokenString string) (bool, error) {
	ctx := context.Background()
	return redis.Exists(ctx, tokenBlacklistPrefix+tokenString)
}

// InvalidateToken 使指定的token立即失效
func InvalidateToken(tokenString string) error {
	// 首先验证token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrInvalidToken
		}
		return jwtSecret, nil
	})
	if err != nil {
		return err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return ErrInvalidToken
	}

	// 获取token的过期时间
	exp, ok := claims["exp"].(float64)
	if !ok {
		return errors.New("invalid token: missing expiration")
	}

	// 计算剩余过期时间
	expTime := time.Unix(int64(exp), 0)
	ttl := time.Until(expTime)
	if ttl <= 0 {
		return nil // token已经过期，无需加入黑名单
	}

	// 将token加入Redis黑名单，使用相同的过期时间
	ctx := context.Background()
	return redis.Set(ctx, tokenBlacklistPrefix+tokenString, "1", ttl)
}

// IsUserBlacklisted 检查用户是否在黑名单中
func IsUserBlacklisted(userID string) (bool, error) {
	ctx := context.Background()
	return redis.Exists(ctx, userBlacklistPrefix+userID)
}

// CleanupBlacklist 已废弃 - Redis会自动清理过期的key
// Deprecated: 此函数不再需要，因为Redis会自动清理过期的key
func CleanupBlacklist() error {
	// 为了向后兼容保留此函数，但实际上不需要执行任何操作
	return nil
}
