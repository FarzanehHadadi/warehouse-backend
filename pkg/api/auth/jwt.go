package auth

import (
	"fmt"
	"sync"
	"time"
	"warehouse/pkg/logger"

	"github.com/golang-jwt/jwt/v5"
	"go.uber.org/zap"
)

const MinJWTSecretKeyBytes = 32

type Claims struct {
	UserId uint `json:"user_id"`
	jwt.RegisteredClaims
}

var (
	AccessSecretKey  []byte
	RefreshSecretKey []byte
	mu               sync.RWMutex
)

func LoadSecrets(accessKey, refreshKey []byte) {
	if len(accessKey) < MinJWTSecretKeyBytes || len(refreshKey) < MinJWTSecretKeyBytes {
		panic("JWT secret keys must be at least 32 bytes")
	}

	mu.Lock()
	defer mu.Unlock()

	AccessSecretKey = append([]byte{}, accessKey...)
	RefreshSecretKey = append([]byte{}, refreshKey...)
}
func JWtSigningKey(tokenType TokenType) []byte {
	if tokenType == TokenTypeAccess {
		if len(AccessSecretKey) == 0 {
			return nil
		}
		mu.RLock()
		out := make([]byte, len(AccessSecretKey))
		copy(out, AccessSecretKey)
		mu.RUnlock()
		return out

	}
	mu.RLock()
	out := make([]byte, len(RefreshSecretKey))
	copy(out, RefreshSecretKey)
	mu.RUnlock()
	return out

}
func JWTSignFunc(key []byte) jwt.Keyfunc {
	return func(token *jwt.Token) (interface{}, error) {
		if token.Method != jwt.SigningMethodHS256 {
			return nil, fmt.Errorf("unexpected signing method,%v", token.Header["alg"])

		}
		return key, nil
	}

}

type TokenType string

const (
	TokenTypeAccess  TokenType = "access"
	TokenTypeRefresh TokenType = "refresh"
)
const (
	AccessTokenDuration  = 1 * time.Hour
	RefreshTokenDuration = 7 * 24 * time.Hour
)

func GenerateToken(userId uint, tokenType TokenType) (string, error) {
	if userId == 0 {
		return "", fmt.Errorf("Invalid user id")
	}
	key := JWtSigningKey(tokenType)
	if len(key) < MinJWTSecretKeyBytes {
		return "", fmt.Errorf("Generated token is not valid")
	}
	tokenDuration := AccessTokenDuration
	if tokenType == TokenTypeRefresh {
		tokenDuration = RefreshTokenDuration
	}

	claims := &Claims{
		UserId: userId,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(tokenDuration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(key)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}
func ValidateToken(tokenString string, tokenType TokenType) (*Claims, error) {
	key := JWtSigningKey(tokenType)
	if len(key) < MinJWTSecretKeyBytes {
		return nil, fmt.Errorf("JWT secret key is not configured")
	}

	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, JWTSignFunc(key))
	if err != nil {
		logger.Log.Error("ValidateToken", zap.Error(err), zap.String("tokenType", string(tokenType)))
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}
	return nil, jwt.ErrTokenInvalidId
}
