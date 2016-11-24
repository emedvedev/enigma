package alphabet

// ToInt returns the alphabet index of a given letter.
func ToInt(char rune) int {
	return int(char - 'A')
}

// ToChar returns the letter with a given alphabet index.
func ToChar(index int) rune {
	return rune('A' + index)
}
