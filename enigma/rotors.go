package enigma

import (
	"github.com/emedvedev/enigma/alphabet"
	"strings"
)

var Rotors = map[string]Rotor{
	"I":     Rotor{"EKMFLGDQVZNTOWYHXUSPAIBRCJ", 0, 0, map[rune]bool{'Q': true}},
	"II":    Rotor{"AJDKSIRUXBLHWTMCQGZNPYFVOE", 0, 0, map[rune]bool{'E': true}},
	"III":   Rotor{"BDFHJLCPRTXVZNYEIWGAKMUSQO", 0, 0, map[rune]bool{'V': true}},
	"IV":    Rotor{"ESOVPZJAYQUIRHXLNFTGKDCMWB", 0, 0, map[rune]bool{'J': true}},
	"V":     Rotor{"VZBRGITYUPSDNHLXAWMJQOFECK", 0, 0, map[rune]bool{'Z': true}},
	"VI":    Rotor{"JPGVOUMFYQBENHZRDKASXLICTW", 0, 0, map[rune]bool{'Z': true, 'M': true}},
	"VII":   Rotor{"NZJHGRCXMYSWBOUFAIVLPEKQDT", 0, 0, map[rune]bool{'Z': true, 'M': true}},
	"VIII":  Rotor{"FKQHTLXOCBJSPDZRAMEWNIUYGV", 0, 0, map[rune]bool{'Z': true, 'M': true}},
	"Beta":  Rotor{"LEYJVCNIXWPBQMDRTAKZGFUHOS", 0, 0, map[rune]bool{}},
	"Gamma": Rotor{"FSOKANUERHMBTIYCWLQPZXVGJD", 0, 0, map[rune]bool{}},
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
	ring     int
	notch    map[rune]bool
}

func (r *Rotor) Step(letter *rune, invert bool) {

	number := strings.IndexRune(alphabet.Alphabet, *letter)
	number = (number - r.ring + r.offset + 26) % 26

	if invert {
		number = strings.IndexRune(r.sequence, rune(alphabet.Alphabet[number]))
	} else {
		number = strings.IndexRune(alphabet.Alphabet, rune(r.sequence[number]))
	}

	number = (number + r.ring - r.offset + 26) % 26
	*letter = rune(alphabet.Alphabet[number])

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
}
