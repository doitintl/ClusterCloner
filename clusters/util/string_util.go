package util

import (
	"github.com/tjarratt/babble"
	"math/rand"
	"regexp"
	"strings"
	"time"
)

func init() {
	rand.Seed(time.Now().UTC().UnixNano())

}

// RandomWord , lowercase, ASCII-only, max 7 letters. Meant for human-memorable identifers.
func RandomWord() string {
	babbler := babble.NewBabbler()
	babbler.Count = 1
	word := babbler.Babble()
	for len(word) > 7 {
		word = babbler.Babble()
	}
	wordLC := strings.ToLower(word)
	re := regexp.MustCompile("[[:^ascii:]]")
	wordASCII := re.ReplaceAllLiteralString(wordLC, "")
	return wordASCII
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
		if e == "" {
			continue
		}
		parts := strings.Split(e, "=")
		if len(parts) == 1 {
			parts = []string{parts[0], ""}
		}
		m[parts[0]] = parts[1]
	}
	return m
}

// ToCommaSeparateKeyValuePairs ...
func ToCommaSeparateKeyValuePairs(m map[string]string) (ret string) {
	for k, v := range m {
		ret += k + "=" + v + ","
	}
	ret = ret[:len(ret)-1]
	return
}
