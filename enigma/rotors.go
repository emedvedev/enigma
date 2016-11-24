package enigma

import (
	"fmt"
	"github.com/emedvedev/enigma/alphabet"
	"strings"
)

var Rotors = map[string]Rotor{
	"I":     Rotor{"EKMFLGDQVZNTOWYHXUSPAIBRCJ", 0, map[rune]bool{'Q': true}},
	"II":    Rotor{"AJDKSIRUXBLHWTMCQGZNPYFVOE", 0, map[rune]bool{'E': true}},
	"III":   Rotor{"BDFHJLCPRTXVZNYEIWGAKMUSQO", 0, map[rune]bool{'V': true}},
	"IV":    Rotor{"ESOVPZJAYQUIRHXLNFTGKDCMWB", 0, map[rune]bool{'J': true}},
	"V":     Rotor{"VZBRGITYUPSDNHLXAWMJQOFECK", 0, map[rune]bool{'Z': true}},
	"VI":    Rotor{"JPGVOUMFYQBENHZRDKASXLICTW", 0, map[rune]bool{'Z': true, 'M': true}},
	"VII":   Rotor{"NZJHGRCXMYSWBOUFAIVLPEKQDT", 0, map[rune]bool{'Z': true, 'M': true}},
	"VIII":  Rotor{"FKQHTLXOCBJSPDZRAMEWNIUYGV", 0, map[rune]bool{'Z': true, 'M': true}},
	"Beta":  Rotor{"LEYJVCNIXWPBQMDRTAKZGFUHOS", 0, map[rune]bool{}},
	"Gamma": Rotor{"FSOKANUERHMBTIYCWLQPZXVGJD", 0, map[rune]bool{}},
}

var Reflectors = map[string]Reflector{
	"A":      Reflector{"EJMZALYXVBWFCRQUONTSPIKHGD"},
	"B":      Reflector{"YRUHQSLDPXNGOKMIEBFZCWVJAT"},
	"C":      Reflector{"FVPJIAOYEDRZXWGCTKUQSBNMHL"},
	"B-Thin": Reflector{"ENKQAUYWJICOPBLMDXZVFTHRGS"},
	"C-Thin": Reflector{"RDOBJNTKVEHMLFCWZAXGYIPSUQ"},
}

type Rotor struct {
	sequence string
	offset   int
	notch    map[rune]bool
}

func (r *Rotor) Step(letter *rune, invert bool) {
	if invert {
		index := alphabet.IntToChar((alphabet.CharToInt(*letter) + r.offset + 26) % 26)
		*letter = alphabet.IntToChar((strings.IndexRune(r.sequence, index) - r.offset + 26) % 26)
	} else {
		index := (alphabet.CharToInt(*letter) + r.offset + 26) % 26
		*letter = alphabet.IntToChar((strings.IndexByte(alphabet.Alphabet, r.sequence[index]) - r.offset + 26) % 26)
	}
	fmt.Print(string(*letter), ">")
}

type RotorConfig struct {
	ID    string
	Start rune
	Ring  int
}

type Reflector struct {
	sequence string
}

func (r *Reflector) Reflect(letter *rune) {
	*letter = alphabet.IntToChar(strings.IndexRune(r.sequence, *letter))
	fmt.Print(string(*letter), ">")
}
