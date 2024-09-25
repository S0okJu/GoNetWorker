package core

import (
	"fmt"
	"math/rand"
)

func generateRandomValue(dataType string) string {
	switch dataType {
	case "string":
		return randomString(7)
	case "int":
		return fmt.Sprintf("%d", rand.Intn(100))
	case "float":
		return fmt.Sprintf("%.2f", rand.Float64()*100)
	case "bool":
		return fmt.Sprintf("%t", rand.Intn(2) == 1)
	default:
		return "unknown_type"
	}
}

func randomString(length int) string {
	if length == 0 {
		length = 1
	}
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}
