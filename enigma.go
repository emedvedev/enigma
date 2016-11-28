// Package enigma is an Enigma cipher machine emulator.
//
// This is a neat little emulator of various Enigma machines
// with a lot of confugurable parameters and a CLI tool.
// Somebody hurt your feelings by saying "my grandmother
// encrypts better than you"? I've got you covered! With
// this port of the amazing 1940's technology you'll be
// just as good at encrypting things as anyone's grandmother.
//
// This package contains a library to use an Enigma machine
// in your own Go code. For the CLI tool use cmd/enigma:
//
// 	go install github.com/emedvedev/enigma/cmd/enigma
//
// The full CLI syntax is a bit verbose, but it's possible to
// use the tool without any source code modifications, config
// files, or Enigma knowledge:
//
//  enigma Never gonna give you up, never gonna let you down!
//
// Using an Enigma machine with default settings is somewhat similar to
// setting your password to "0000". Let's up our security game:
//
//  enigma youtu.be/dQw4w9WgXcQ --rotors Beta VI I III --reflector C-Thin --plugboard AD SF ET RY HK JL QZ WX UM OP --rings 10 5 16 10
//
// Much better. And of course, `enigma -h` will give you the complete
// description of parameters and usage.
//
// Importantly, since Enigma machines only have 26 keys, spaces
// are replaced with X, and everything outside of the English alphabet
// is discarded. It's up to you to come up with a suitable encoding.
//
// Enjoy!
//
// Enigma models and features
//
// Almost everything from the German Enigma machines can be configured in this
// emulator:
//
// — Rotor set: rotors from M3 and M4, the most famous Enigma machines,
// are pre-loaded.
//
// — Reflector: reflectors A, B, and C — as well as the thin B and C
// versions used in M4 — are supported.
//
// — Plugboard: any number of letter pairs is accepted. Plugboard
// configuration is optional.
//
// — Ring offsets and starting position of the rotors.
//
// M3 and M4 can be fully emulated with the right parameters, and if it's
// not enough, new rotors and reflectors can be added quite easily: just
// add a new entry to the list in `rotors.go`, and that's it. Notches for
// rotor turnover are optional.
//
// Some exotic Enigma variants and implementations, as well
// as devices such as Uhr, are not supported due to my chronic lack of
// spare time. Your pull requests would be most welcome!
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
type Plugboard map[byte]byte

// NewPlugboard is the plugboard constructor accepting an array
// of two-symbol strings representing plug pairs.
func NewPlugboard(pairs []string) *Plugboard {
	p := Plugboard{}
	for _, pair := range pairs {
		if len(pair) > 0 {
			p[pair[0]] = pair[1]
			p[pair[1]] = pair[0]
		}
	}
	return &p
}

func (r *Rotor) move(offset int) {
	r.Offset = (r.Offset + offset + 26) % 26
}

func (e *Enigma) moveRotors() {
	var (
		rotorLen         = len(e.Rotors)
		farRight         = &e.Rotors[rotorLen-1]
		farRightNotch    = farRight.Notch[ToChar((farRight.Offset+26)%26)]
		secondRight      = &e.Rotors[rotorLen-2]
		secondRightNotch = secondRight.Notch[ToChar((secondRight.Offset+26)%26)]
		thirdRight       = &e.Rotors[rotorLen-3]
	)
	farRight.move(1)
	if farRightNotch {
		secondRight.move(1)
	}
	if secondRightNotch {
		if !farRightNotch {
			secondRight.move(1)
		}
		thirdRight.move(1)
	}
}

// EncryptChar inputs a single character into the machine.
func (e *Enigma) EncryptChar(letter byte) byte {
	e.moveRotors()
	if value, ok := e.Plugboard[letter]; ok {
		letter = value
	}
	for i := len(e.Rotors) - 1; i >= 0; i-- {
		e.Rotors[i].Step(&letter, false)
	}
	e.Reflector.Reflect(&letter)
	for i := 0; i < len(e.Rotors); i++ {
		e.Rotors[i].Step(&letter, true)
	}
	if value, ok := e.Plugboard[letter]; ok {
		letter = value
	}
	return letter
}

// EncryptString inputs a string into the machine.
func (e *Enigma) EncryptString(text string) string {
	encrypted := make([]byte, len(text))
	for i := range text {
		encrypted[i] = e.EncryptChar(text[i])
	}
	return string(encrypted)
}
