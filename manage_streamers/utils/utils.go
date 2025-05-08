package utils

import (
	"crypto/rand"
	"encoding/hex"
)

func GenerateStreamKey() string {
	bytes := make([]byte, 16)
	if _, err := rand.Read(bytes); err != nil {
		return "fallback-streamkey"
	}
	return hex.EncodeToString(bytes)
}
