package utils

import (
	"crypto/rand"
	"encoding/hex"
	"log"
)

func GenerateCSRFToken() string {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		log.Printf("Error generating CSRF token: %v", err)
		return ""
	}
	return hex.EncodeToString(bytes)
}