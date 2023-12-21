package utils

import "math/rand"

func GenerateTraceID() string {
	const alphabet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890"
	const transIDLen = 10

	transID := make([]byte, transIDLen)

	for i := range transID {
		transID[i] = alphabet[rand.Intn(len(alphabet))]
	}

	return string(transID)
}
