package utils

import (
	"math/rand"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func RandomUrl() string {
	str := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	return RandomString(str, 6)

}

func RandomString(str string, length int) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = str[RandomInt(0, len(str))]
	}
	return string(b)
}

func RandomInt(min, max int) int {
	return min + rand.Intn(max-min)
}
