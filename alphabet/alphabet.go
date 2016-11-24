package alphabet

import (
	"strings"
)

var Alphabet = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"

func CharToInt(char rune) int {
	return strings.IndexRune(Alphabet, char)
}
func IntToChar(index int) rune {
	return rune(Alphabet[index])
}
