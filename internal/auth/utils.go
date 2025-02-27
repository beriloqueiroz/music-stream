package auth

import (
	"crypto/rand"
	"encoding/base64"
)

func generateRandomCode() string {
	b := make([]byte, 6)
	rand.Read(b)
	return base64.URLEncoding.EncodeToString(b)[:8]
}
