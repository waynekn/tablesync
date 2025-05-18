package utils

import (
	"math/big"
	"strings"

	"github.com/google/uuid"
)

const base62Alphabet = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
const idLength = 22

// GenerateID generates a unique ID using UUID and encodes it in base62 format.
func GenerateID() string {
	uuid := uuid.New()

	var intVal big.Int
	intVal.SetBytes(uuid[:])

	// Encode to base62
	var base62 strings.Builder
	base := big.NewInt(62)
	zero := big.NewInt(0)
	for intVal.Cmp(zero) > 0 {
		mod := new(big.Int)
		intVal.DivMod(&intVal, base, mod)
		base62.WriteByte(base62Alphabet[mod.Int64()])
	}

	// Reverse the string to get the correct order
	id := reverseString(base62.String())

	if len(id) > idLength {
		return id[:idLength]
	}

	return id
}

// reverseString reverses a string.
func reverseString(s string) string {
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}
