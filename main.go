package main

import (
	"fmt"
	"strings"
)

var Rotors = map[string]Rotor{
	"I":     *NewRotor("EKMFLGDQVZNTOWYHXUSPAIBRCJ", []rune{'R'}, []rune{'Q'}, []rune{'Y'}),
	"II":    *NewRotor("AJDKSIRUXBLHWTMCQGZNPYFVOE", []rune{'F'}, []rune{'E'}, []rune{'M'}),
	"III":   *NewRotor("BDFHJLCPRTXVZNYEIWGAKMUSQO", []rune{'W'}, []rune{'V'}, []rune{'D'}),
	"IV":    *NewRotor("ESOVPZJAYQUIRHXLNFTGKDCMWB", []rune{'K'}, []rune{'J'}, []rune{'R'}),
	"V":     *NewRotor("VZBRGITYUPSDNHLXAWMJQOFECK", []rune{'A'}, []rune{'Z'}, []rune{'H'}),
	"VI":    *NewRotor("JPGVOUMFYQBENHZRDKASXLICTW", []rune{'A', 'N'}, []rune{'Z', 'M'}, []rune{'H', 'U'}),
	"VII":   *NewRotor("NZJHGRCXMYSWBOUFAIVLPEKQDT", []rune{'A', 'N'}, []rune{'Z', 'M'}, []rune{'H', 'U'}),
	"VIII":  *NewRotor("FKQHTLXOCBJSPDZRAMEWNIUYGV", []rune{'A', 'N'}, []rune{'Z', 'M'}, []rune{'H', 'U'}),
	"Beta":  *NewRotor("LEYJVCNIXWPBQMDRTAKZGFUHOS", []rune{}, []rune{}, []rune{}),
	"Gamma": *NewRotor("FSOKANUERHMBTIYCWLQPZXVGJD", []rune{}, []rune{}, []rune{}),
}

var Reflectors = map[string]Reflector{
	"A":      *NewReflector("EJMZALYXVBWFCRQUONTSPIKHGD"),
	"B":      *NewReflector("YRUHQSLDPXNGOKMIEBFZCWVJAT"),
	"C":      *NewReflector("FVPJIAOYEDRZXWGCTKUQSBNMHL"),
	"B Thin": *NewReflector("ENKQAUYWJICOPBLMDXZVFTHRGS"),
	"C Thin": *NewReflector("RDOBJNTKVEHMLFCWZAXGYIPSUQ"),
}

type Rotor struct {
	Sequence                string
	StraightPairs           [26][2]rune
	ReversePairs            [26][2]rune
	Offset                  int
	Turnover, Window, Notch []rune
}

var baseSequence = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"

func NewRotor(sequence string, turnover, window, notch []rune) *Rotor {
	rotor := &Rotor{Turnover: turnover, Window: window, Notch: notch, Sequence: sequence}
	for i, letter := range sequence {
		rotor.StraightPairs[i] = [2]rune{rune(baseSequence[i]), letter}
		rotor.ReversePairs[i] = [2]rune{letter, rune(baseSequence[i])}
	}
	return rotor
}

func (r *Rotor) Step(letter *rune, invert bool) rune {
	if invert {
		*letter = rune(baseSequence[(strings.IndexRune(baseSequence, *letter)+r.Offset+26)%26])
		index := (strings.IndexRune(r.Sequence, *letter) - r.Offset + 26) % 26
		*letter = rune(baseSequence[index])
	} else {
		index := (strings.IndexRune(baseSequence, *letter) + r.Offset + 26) % 26
		*letter = rune(baseSequence[(strings.IndexByte(baseSequence, r.Sequence[index])-r.Offset+26)%26])
	}
	//fmt.Print(string(*letter))
	return *letter
}

type Reflector struct {
	Pairs map[rune]rune
}

func NewReflector(sequence string) *Reflector {
	reflector := &Reflector{}
	reflector.Pairs = make(map[rune]rune, len(sequence))
	for i, letter := range sequence {
		reflector.Pairs[letter] = rune(baseSequence[i])
	}
	return reflector
}

func (r *Reflector) Reflect(letter *rune) rune {
	*letter = r.Pairs[*letter]
	return *letter
}

type Enigma struct {
	rotors    []Rotor
	reflector Reflector
}

func (e *Enigma) MoveRotors() {
	e.rotors[2].Offset++
}

func (e *Enigma) EncryptChar(letter *rune) rune {
	e.MoveRotors()
	for i := len(e.rotors) - 1; i >= 0; i-- {
		e.rotors[i].Step(letter, false)
	}
	e.reflector.Reflect(letter)
	for i := 0; i < len(e.rotors); i++ {
		e.rotors[i].Step(letter, true)
	}
	return *letter
}

func NewEnigma(rotorConfigurations [][]string, reflectorID string) *Enigma {
	rotors := make([]Rotor, len(rotorConfigurations))
	for i, configuration := range rotorConfigurations {
		rotors[i] = Rotors[configuration[0]]
		rotors[i].Offset = strings.IndexRune(baseSequence, rune(configuration[1][0]))
	}
	return &Enigma{rotors, Reflectors[reflectorID]}
}

func main() {
	enigma := NewEnigma([][]string{{"III", "A"}, {"II", "A"}, {"I", "A"}}, "B")

	plaintext := "HELLOWORLD"
	fmt.Println(plaintext)

	for index := range plaintext {
		char := rune(plaintext[index])
		enigma.EncryptChar(&char)
		fmt.Print(string(char))
	}
}
