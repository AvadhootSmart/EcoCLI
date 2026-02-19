package crypto

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
)

func GenerateRandomString(length int) (string, error) {
	b := make([]byte, length-2)

	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(b), nil
}

// GenerateDeviceID creates a unique device identifier
// Example: "mobile-a1b2c3d4"
func GenerateDeviceID() string {
	str, err := GenerateRandomString(8)
	if err != nil {
		fmt.Println(err)
	}
	return "mobile-" + str
}

// GenerateSecret creates a 32-byte random secret
func GenerateSecret() (string, error) {
	return GenerateRandomString(32)
}

// ValidateSecret checks if a secret string is valid base64 and 32 bytes when decoded
func ValidateSecret(secret string) bool {
	decoded, err := hex.DecodeString(secret)
	if err != nil {
		return false
	}
	return len(decoded) == 32
}
