package common

import (
	"crypto/sha512"
	"encoding/hex"
	"math/rand"
	"time"
)

// It generates random token (insecure for prod)
var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func GenerateToken() string {
	rand.Seed(time.Now().UnixNano())
	b := make([]rune, 30)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func HashPassword(text string) string {
	algorithm := sha512.New()
	algorithm.Write([]byte(text))
	return hex.EncodeToString(algorithm.Sum(nil))
}
