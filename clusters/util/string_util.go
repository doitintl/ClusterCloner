package util

import (
	"math/rand"
	"strings"
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

// CommaSeparatedKeyValPairsToMap ...
func CommaSeparatedKeyValPairsToMap(s string) map[string]string {
	s = strings.Trim(s, " ")

	entries := strings.Split(s, ",")

	m := make(map[string]string)
	if s == "" {
		return m
	}
	for _, e := range entries {
		parts := strings.Split(e, "=")
		m[parts[0]] = parts[1]
	}
	return m
}
