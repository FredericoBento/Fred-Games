package utils

import (
	"math/rand"
	"time"
)

func RandomString(length int) string {
	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	rand.Seed(time.Now().UnixNano()) // Seed the random number generator
	result := make([]byte, length)

	for i := range result {
		result[i] = letters[rand.Intn(len(letters))] // Select a random character
	}

	return string(result)
}
