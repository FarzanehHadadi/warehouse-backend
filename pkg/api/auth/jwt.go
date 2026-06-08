package auth

import (
	"fmt"
	"sync"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const MinJWTSecretKeyBytes = 32

type Claims struct {
	UserId uint `json:"user_id"`
	jwt.RegisteredClaims
}

var (
	SecretKey []byte
	signingMu sync.RWMutex
)

func LoadSecretKey(secret []byte) {

	if len(secret) < MinJWTSecretKeyBytes {
		panic("Invalid jwt secret key")
	}
	k := make([]byte, len(secret))
	copy(k, secret)
	signingMu.Lock()
	SecretKey = k
	signingMu.Unlock()

}
func JWtSigningKey() []byte {
	if len(SecretKey) == 0 {
		return nil
	}
	signingMu.RLock()
	out := make([]byte, len(SecretKey))
	copy(out, SecretKey)
	signingMu.RUnlock()
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
func GenerateToken(userId uint) (string, error) {
	if userId == 0 {
		return "", fmt.Errorf("Invalid user id")
	}
	key := JWtSigningKey()
	if len(key) < MinJWTSecretKeyBytes {
		return "", fmt.Errorf("Generated token is not valid")
	}

	claims := &Claims{
		UserId: userId,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(80 * 24 * time.Hour)),
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
func ValidateToken(tokenString string) (*Claims, error) {
	key := JWtSigningKey()
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, JWTSignFunc(key))

	if err != nil {
		return nil, err

	}
	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}
	return nil, jwt.ErrTokenInvalidId
}
