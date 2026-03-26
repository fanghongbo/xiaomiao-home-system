package jwt

import (
	"errors"
	"fmt"
	"strings"
	"time"

	jwtv5 "github.com/golang-jwt/jwt/v5"
)

var (
	ErrSecretNotSet   = errors.New("jwt secret key not set")
	ErrInvalidToken   = errors.New("invalid jwt token")
	ErrExpiredToken   = errors.New("jwt token expired")
	ErrNotValidYet    = errors.New("jwt token not valid yet")
	ErrUnexpectedSign = errors.New("unexpected jwt signing method")
)

type Config struct {
	SecretKey []byte
	Issuer    string
	TTL       time.Duration
}

// Claims 内置用户ID，兼容扩展字段。
type Claims struct {
	UserId   int64  `json:"uid"`
	Nickname string `json:"nickname"`
	jwtv5.RegisteredClaims
}

// GenerateToken 生成用户登录 token（HS256）。
func GenerateToken(config *Config, userId int64, nickname string) (string, error) {
	secret := config.SecretKey
	issuer := config.Issuer
	ttl := config.TTL

	if len(secret) == 0 {
		return "", ErrSecretNotSet
	}
	if userId <= 0 {
		return "", fmt.Errorf("invalid userID: %d", userId)
	}

	now := time.Now()
	claims := Claims{
		UserId:   userId,
		Nickname: nickname,
		RegisteredClaims: jwtv5.RegisteredClaims{
			Issuer:    issuer,
			Subject:   fmt.Sprintf("%d", userId),
			IssuedAt:  jwtv5.NewNumericDate(now),
			NotBefore: jwtv5.NewNumericDate(now.Add(-30 * time.Second)),
			ExpiresAt: jwtv5.NewNumericDate(now.Add(ttl)),
		},
	}

	token := jwtv5.NewWithClaims(jwtv5.SigningMethodHS256, claims)
	return token.SignedString(secret)
}

// ParseToken 解析并校验 token，返回 claims。
func ParseToken(config *Config, tokenString string) (*Claims, error) {
	tokenString = strings.TrimSpace(tokenString)
	if tokenString == "" {
		return nil, ErrInvalidToken
	}

	secret := config.SecretKey
	if len(secret) == 0 {
		return nil, ErrSecretNotSet
	}

	var claims Claims
	_, err := jwtv5.ParseWithClaims(tokenString, &claims, func(t *jwtv5.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwtv5.SigningMethodHMAC); !ok {
			return nil, ErrUnexpectedSign
		}
		if t.Method.Alg() != jwtv5.SigningMethodHS256.Alg() {
			return nil, ErrUnexpectedSign
		}
		return secret, nil
	}, jwtv5.WithLeeway(30*time.Second))
	if err != nil {
		switch {
		case errors.Is(err, jwtv5.ErrTokenExpired):
			return nil, ErrExpiredToken
		case errors.Is(err, jwtv5.ErrTokenNotValidYet):
			return nil, ErrNotValidYet
		default:
			return nil, ErrInvalidToken
		}
	}

	if claims.UserId <= 0 {
		return nil, ErrInvalidToken
	}
	return &claims, nil
}

// ValidateToken 校验 token 并返回 userID。
func ValidateToken(config *Config, tokenString string) (int64, error) {
	claims, err := ParseToken(config, tokenString)
	if err != nil {
		return 0, err
	}
	return claims.UserId, nil
}

// BearerTokenFromAuthorization 从 "Authorization: Bearer <token>" 中提取 token。
func BearerTokenFromAuthorization(authorization string) string {
	authorization = strings.TrimSpace(authorization)
	if authorization == "" {
		return ""
	}
	const prefix = "Bearer "
	if len(authorization) < len(prefix) {
		return ""
	}
	if !strings.EqualFold(authorization[:len(prefix)], prefix) {
		return ""
	}
	return strings.TrimSpace(authorization[len(prefix):])
}
