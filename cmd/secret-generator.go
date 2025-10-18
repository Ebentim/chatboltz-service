package main

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
)

func GenerateSecret() string {
	b := make([]byte, 64)
	rand.Read(b)
	secret := base64.StdEncoding.EncodeToString(b)
	fmt.Println("Generated secret:", secret)
	return secret
}

// func main() {
// 	GenerateSecret()
// }
