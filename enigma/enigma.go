package enigma

import (
	"github.com/emedvedev/enigma/alphabet"
	"strings"
)

type Enigma struct {
	reflector Reflector
	rotors    []Rotor
	plugboard map[rune]rune
}

func (e *Enigma) MoveRotors() {
	rotate := make(map[int]int, len(e.rotors))
	rotate[len(e.rotors)-1] = 1

	if e.rotors[len(e.rotors)-1].notch[alphabet.IntToChar((e.rotors[len(e.rotors)-1].offset+26)%26)] {
		rotate[len(e.rotors)-2] = 1
	}
	if e.rotors[len(e.rotors)-2].notch[alphabet.IntToChar((e.rotors[len(e.rotors)-2].offset+26)%26)] {
		rotate[len(e.rotors)-2] = 1
		rotate[len(e.rotors)-3] = 1
	}

	for rotor, offset := range rotate {
		e.rotors[rotor].offset += offset
	}
}

func (e *Enigma) EncryptChar(letter *rune) {
	e.MoveRotors()

	if value, ok := e.plugboard[*letter]; ok {
		*letter = value
	}
	for i := len(e.rotors) - 1; i >= 0; i-- {
		e.rotors[i].Step(letter, false)
	}
	e.reflector.Reflect(letter)
	for i := 0; i < len(e.rotors); i++ {
		e.rotors[i].Step(letter, true)
	}
	if value, ok := e.plugboard[*letter]; ok {
		*letter = value
	}

}

func NewEnigma(rotorConfiguration []RotorConfig, reflectorID string, plugboardConfiguration [][2]rune) *Enigma {
	rotors := make([]Rotor, len(rotorConfiguration))
	for i, configuration := range rotorConfiguration {
		rotors[i] = Rotors[configuration.ID]
		rotors[i].offset = strings.IndexRune(alphabet.Alphabet, configuration.Start)
		rotors[i].ring = configuration.Ring - 1
	}
	plugboard := map[rune]rune{}
	for _, pair := range plugboardConfiguration {
		plugboard[pair[0]] = pair[1]
		plugboard[pair[1]] = pair[0]
	}
	return &Enigma{Reflectors[reflectorID], rotors, plugboard}
}
