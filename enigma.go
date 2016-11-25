// Package enigma is an Enigma cipher machine emulator.
//
// This package contains a library to use an Enigma machine
// in your own Go code, but there is also a companion CLI
// tool:
//
// 	go install github.com/emedvedev/enigma/cmd/enigma
//
// While the full CLI syntax is a bit verbose, it's actually possible
// to use the tool without any source code modifications, config
// files, or Enigma knowledge:
//
//  enigma Never gonna give you up, never gonna let you down!
//
// Using an Enigma machine with default settings is somewhat similar to
// setting your password to "0000". Let's up our security game:
//
//  enigma youtu.be/dQw4w9WgXcQ --rotors Beta VI I III  --reflector C-Thin --plugboard AD SF ET RY HK JL QZ WX UM OP --rings 10 5 16 10
//
// Much better.
//
// Support
//
// - Enigma M3 and M4 are supported with the pre-defined settings.
//
// - Other Enigma models might be supported depending on rotor
// turnover behavior.
//
// - Additional features and devices such as Uhr are not supported.
//
// Importantly, everything except English letters is discarded, since
// Enigma machines only had 26 keys. It's up to you to come up with
// a suitable encoding.
package enigma

// Enigma represents an Enigma machine with configured rotors, plugs,
// and a reflector. Most states are stored in the rotors themselves.
type Enigma struct {
	Reflector Reflector
	Plugboard Plugboard
	Rotors    []Rotor
}

// NewEnigma is the Enigma constructor, accepting an array of RotorConfig objects
// for rotors, a reflector ID/name, and an array of plugboard pairs.
func NewEnigma(rotorConfiguration []RotorConfig, reflectorID string, plugs []string) *Enigma {
	rotors := make([]Rotor, len(rotorConfiguration))
	for i, configuration := range rotorConfiguration {
		rotors[i] = Rotors[configuration.ID]
		rotors[i].Offset = ToInt(configuration.Start)
		rotors[i].Ring = configuration.Ring - 1
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
	var rotate = make(map[int]int, len(e.Rotors))
	rotate[len(e.Rotors)-1] = 1
	if e.Rotors[len(e.Rotors)-1].Notch[ToChar((e.Rotors[len(e.Rotors)-1].Offset+26)%26)] {
		rotate[len(e.Rotors)-2] = 1
	}
	if e.Rotors[len(e.Rotors)-2].Notch[ToChar((e.Rotors[len(e.Rotors)-2].Offset+26)%26)] {
		rotate[len(e.Rotors)-2] = 1
		rotate[len(e.Rotors)-3] = 1
	}
	for rotor, offset := range rotate {
		var newOffset = (e.Rotors[rotor].Offset + offset + 26) % 26
		e.Rotors[rotor].Offset = newOffset
	}
}

// EncryptChar inputs a single character into the machine.
func (e *Enigma) EncryptChar(letter *rune) {
	e.moveRotors()
	if value, ok := e.Plugboard[*letter]; ok {
		*letter = value
	}
	for i := len(e.Rotors) - 1; i >= 0; i-- {
		e.Rotors[i].Step(letter, false)
	}
	e.Reflector.Reflect(letter)
	for i := 0; i < len(e.Rotors); i++ {
		e.Rotors[i].Step(letter, true)
	}
	if value, ok := e.Plugboard[*letter]; ok {
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
