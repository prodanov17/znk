package utils

import (
	"math/rand"
	"time"
)

func StringPtr(s string) *string {
	return &s
}

const letters = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"

// GenerateRandomString generates a random string of a specified length
func GenerateRandomString(length int) string {
	rand.Seed(time.Now().UnixNano()) // Seed the random number generator
	result := make([]byte, length)

	for i := range result {
		result[i] = letters[rand.Intn(len(letters))]
	}
	return string(result)
}
