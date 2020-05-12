package util

import (
	"math/rand"
	"time"
)

func init() {
	rand.Seed(time.Now().UTC().UnixNano())

}

// RandomAlphaNumSequence ...
func RandomAlphaNumSequence(length int, includeUpperCase, includeLowerCase, includeDigits bool) string {
	s := ""
	if includeUpperCase {
		s += "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	}
	if includeLowerCase {
		s += "abcdefghijklmnopqrstuvwxyz"
	}
	if includeDigits {
		s += "01234356789"
	}
	letters := []rune(s)
	b := make([]rune, length)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}
