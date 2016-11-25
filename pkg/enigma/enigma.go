// Package enigma is an Enigma cipher machine emulator.
//
// This package contains a library to use an Enigma machine
// in your own Go code, but there is also a companion CLI
// tool:
//
// go install github.com/emedvedev/enigma/cmd/enigma
//
// While the CLI syntax is a bit verbose, it's actually possible
// to use the tool without any source code modifications, config
// files, or Enigma knowledge.
//
// Support:
//
// 	- Enigma M3 and M4 are supported with the pre-loaded settings.
//  - Other Enigma devices might be supporting depending on rotor
//		turnover behavior.
// 	- Additional devices such as Uhr are not supported.
package enigma

// Enigma represents an Enigma machine with configured rotors, plugs,
// and a reflector. Most states are stored in the rotors themselves.
type Enigma struct {
	reflector Reflector
	plugboard Plugboard
	rotors    []Rotor
}

// NewEnigma is the Enigma constructor, accepting an array of RotorConfig objects
// for rotors, a reflector ID/name and an array of plugs.
func NewEnigma(rotorConfiguration []RotorConfig, reflectorID string, plugs []string) *Enigma {
	rotors := make([]Rotor, len(rotorConfiguration))
	for i, configuration := range rotorConfiguration {
		rotors[i] = Rotors[configuration.ID]
		rotors[i].offset = ToInt(configuration.Start)
		rotors[i].ring = configuration.Ring - 1
	}
	return &Enigma{Reflectors[reflectorID], *NewPlugboard(plugs), rotors}
}

// Plugboard is a two-way mapping between characters modifying the
// encryption/decryption procedure of the Enigma machine.
type Plugboard map[rune]rune

// NewPlugboard is the plugboard constructor accepting an array
// of two-symbol strings representing plug pairs.
func NewPlugboard(pairs []string) *Plugboard {
	p := Plugboard{}
	for _, pair := range pairs {
		if len(pair) > 0 {
			p[rune(pair[0])] = rune(pair[1])
			p[rune(pair[1])] = rune(pair[0])
		}
	}
	return &p
}

func (e *Enigma) moveRotors() {
	var rotate = make(map[int]int, len(e.rotors))
	rotate[len(e.rotors)-1] = 1
	if e.rotors[len(e.rotors)-1].notch[ToChar((e.rotors[len(e.rotors)-1].offset+26)%26)] {
		rotate[len(e.rotors)-2] = 1
	}
	if e.rotors[len(e.rotors)-2].notch[ToChar((e.rotors[len(e.rotors)-2].offset+26)%26)] {
		rotate[len(e.rotors)-2] = 1
		rotate[len(e.rotors)-3] = 1
	}
	for rotor, offset := range rotate {
		var newOffset = (e.rotors[rotor].offset + offset + 26) % 26
		e.rotors[rotor].offset = newOffset
	}
}

// EncryptChar inputs a single character into the machine.
func (e *Enigma) EncryptChar(letter *rune) {
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

// EncryptString inputs a string into the machine.
func (e *Enigma) EncryptString(text string) string {
	var encrypted string
	for _, char := range text {
		e.EncryptChar(&char)
		encrypted += string(char)
	}
	return encrypted
}
