package helper

import (
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
