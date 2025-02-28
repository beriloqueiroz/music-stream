package main

import (
	"fmt"
	"log"

	"golang.org/x/crypto/bcrypt"
)

func main2() {
	password := []byte("12365478")
	hash, err := bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Hash da senha: %s\n", string(hash))

	// Teste de verificação
	err = bcrypt.CompareHashAndPassword(hash, password)
	if err != nil {
		fmt.Println("Falha na verificação:", err)
	} else {
		fmt.Println("Verificação bem sucedida!")
	}
}
