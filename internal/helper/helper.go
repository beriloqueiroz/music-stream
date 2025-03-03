package helper

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"log"

	"golang.org/x/crypto/bcrypt"
)

func GenerateHash(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Hash da senha: %s\n", string(hash))

	return string(hash), nil
}

func RemoveFromSlice[T any](slice []T, condition func(T) bool) []T {
	for i, v := range slice {
		if condition(v) {
			return append(slice[:i], slice[i+1:]...)
		}
	}
	return slice
}

func GenerateRandomCode() string {
	b := make([]byte, 6)
	rand.Read(b)
	return base64.URLEncoding.EncodeToString(b)[:8]
}
