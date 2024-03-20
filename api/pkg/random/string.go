package random

import (
	"math/rand"
)

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890")

func NewStrings(digit int) string {
	b := make([]rune, digit)
	for i := range b {
		//nolint:gosec
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}
