package enigma

import (
	"regexp"
	"strings"
)

// CharToIndex returns the alphabet index of a given letter.
func CharToIndex(char byte) int {
	return int(char - 'A')
}

// IndexToChar returns the letter with a given alphabet index.
func IndexToChar(index int) byte {
	return byte('A' + index)
}

// SanitizePlaintext will prepare a string to be encoded
// in the Enigma machine: everything except A-Z will be
// stripped, spaces will be replaced with "X".
func SanitizePlaintext(plaintext string) string {
	plaintext = strings.TrimSpace(plaintext)
	plaintext = strings.ToUpper(plaintext)
	plaintext = strings.Replace(plaintext, " ", "X", -1)
	plaintext = regexp.MustCompile(`[^A-Z]`).ReplaceAllString(plaintext, "")
	return plaintext
}
