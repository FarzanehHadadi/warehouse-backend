package utils

import (
	"crypto/rand"
	"encoding/base64"
)

func GenerateRandomKey() string {
	key := make([]byte, 32)
	_, err := rand.Read(key)
	if err != nil {
		panic("Failed to generate random key:" + err.Error())
	}
	return base64.StdEncoding.EncodeToString(key)
}
