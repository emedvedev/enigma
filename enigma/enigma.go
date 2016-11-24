package enigma

import (
	abc "github.com/emedvedev/enigma/alphabet"
)

// Enigma represents an Enigma machine with configured rotors, plugs,
// and a reflector. Most states are stored in the rotors themselves.
type Enigma struct {
	reflector Reflector
	plugboard Plugboard
	rotors    []Rotor
}

// NewEnigma is a simplified constructor for the Enigma machine, accepting
// an array of RotorConfig objects for rotors, and strings for reflectorID
// and plugs.
func NewEnigma(rotorConfiguration []RotorConfig, reflectorID string, plugs []string) *Enigma {
	rotors := make([]Rotor, len(rotorConfiguration))
	for i, configuration := range rotorConfiguration {
		rotors[i] = Rotors[configuration.ID]
		rotors[i].offset = abc.ToInt(configuration.Start)
		rotors[i].ring = configuration.Ring - 1
	}
	return &Enigma{Reflectors[reflectorID], *NewPlugboard(plugs...), rotors}
}

// Plugboard is a two-way mapping between characters modifying the
// encryption/decryption procedure of the Enigma machine.
type Plugboard map[rune]rune

// NewPlugboard is the plugboard constructor accepting pairs as a
// series of two-symbol strings: NewPlugboard("AB", "CD", "EF")
func NewPlugboard(pairs ...string) *Plugboard {
	p := Plugboard{}
	for _, pair := range pairs {
		p[rune(pair[0])] = rune(pair[1])
		p[rune(pair[1])] = rune(pair[0])
	}
	return &p
}

func (e *Enigma) moveRotors() {
	rotate := make(map[int]int, len(e.rotors))
	rotate[len(e.rotors)-1] = 1
	if e.rotors[len(e.rotors)-1].notch[abc.ToChar((e.rotors[len(e.rotors)-1].offset+26)%26)] {
		rotate[len(e.rotors)-2] = 1
	}
	if e.rotors[len(e.rotors)-2].notch[abc.ToChar((e.rotors[len(e.rotors)-2].offset+26)%26)] {
		rotate[len(e.rotors)-2] = 1
		rotate[len(e.rotors)-3] = 1
	}
	for rotor, offset := range rotate {
		e.rotors[rotor].offset += offset
	}
}

func (e *Enigma) encryptChar(letter *rune) {
	e.moveRotors()
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

// EncryptString is the only exposed function of an Enigma machine
// after it's been created. Does exactly what it says it does.
func (e *Enigma) EncryptString(text string) string {
	var encrypted string
	for _, char := range text {
		e.encryptChar(&char)
		encrypted += string(char)
	}
	return encrypted
}
