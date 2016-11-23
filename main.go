package main

import (
	"fmt"
)

type Rotor struct {
	StraightPairs           map[rune]rune
	ReversePairs            map[rune]rune
	Turnover, Window, Notch []rune
}

var baseSequence = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"

func NewRotor(sequence string, turnover, window, notch []rune) *Rotor {
	rotor := &Rotor{Turnover: turnover, Window: window, Notch: notch}
	rotor.StraightPairs = make(map[rune]rune, len(sequence))
	rotor.ReversePairs = make(map[rune]rune, len(sequence))
	for i, letter := range sequence {
		rotor.StraightPairs[rune(baseSequence[i])] = letter
		rotor.ReversePairs[letter] = rune(baseSequence[i])
	}
	return rotor
}

func (r *Rotor) Step(letter *rune, invert bool) rune {
	if invert {
		*letter = r.ReversePairs[*letter]
	} else {
		*letter = r.StraightPairs[*letter]
	}
	fmt.Print(string(*letter))
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
	fmt.Print(string(*letter))
	return *letter
}

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

type Enigma struct {
	rotors    []Rotor
	reflector Reflector
}

func (e *Enigma) EncryptChar(letter *rune) rune {
	for i := len(e.rotors) - 1; i >= 0; i-- {
		e.rotors[i].Step(letter, false)
	}
	e.reflector.Reflect(letter)
	for i := 0; i < len(e.rotors); i++ {
		e.rotors[i].Step(letter, true)
	}
	return *letter
}

func NewEnigma(rotorIDs []string, reflectorID string) *Enigma {
	rotors := make([]Rotor, len(rotorIDs))
	for i, name := range rotorIDs {
		rotors[i] = Rotors[name]
	}
	return &Enigma{rotors, Reflectors[reflectorID]}
}

func main() {
	enigma := NewEnigma([]string{"I", "II", "III"}, "B")

	letter := 'G'
	fmt.Print(string(letter))

	enigma.EncryptChar(&letter)
}
