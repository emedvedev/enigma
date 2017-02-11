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
// 	go get github.com/emedvedev/enigma/cmd/enigma
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

import "bytes"

// Enigma represents an Enigma machine with configured rotors, plugs,
// and a reflector. Most states are stored in the rotors themselves.
type Enigma struct {
	Reflector Reflector
	Plugboard Plugboard
	Rotors    []*Rotor
}

// RotorConfig reprensents a configuration for a rotor as set by the user:
// ID from the pre-defined list, a starting position (A to Z), and a ring
// setting (1 to 26).
type RotorConfig struct {
	ID    string
	Start byte
	Ring  int
}

// NewEnigma is the Enigma constructor, accepting an array of RotorConfig objects
// for rotors, a reflector ID/name, and an array of plugboard pairs.
func NewEnigma(rotorConfiguration []RotorConfig, refID string, plugs []string) *Enigma {
	rotors := make([]*Rotor, len(rotorConfiguration))
	for i, configuration := range rotorConfiguration {
		rotors[i] = HistoricRotors.GetByID(configuration.ID)
		rotors[i].Offset = CharToIndex(configuration.Start)
		rotors[i].Ring = configuration.Ring - 1
	}
	return &Enigma{*HistoricReflectors.GetByID(refID), *NewPlugboard(plugs), rotors}
}

func (e *Enigma) moveRotors() {
	var (
		rotorLen            = len(e.Rotors)
		farRight            = e.Rotors[rotorLen-1]
		farRightTurnover    = farRight.ShouldTurnOver()
		secondRight         = e.Rotors[rotorLen-2]
		secondRightTurnover = secondRight.ShouldTurnOver()
		thirdRight          = e.Rotors[rotorLen-3]
	)
	if secondRightTurnover {
		if !farRightTurnover {
			secondRight.move(1)
		}
		thirdRight.move(1)
	}
	if farRightTurnover {
		secondRight.move(1)
	}
	farRight.move(1)
}

// EncodeChar encodes a single character.
func (e *Enigma) EncodeChar(letter byte) byte {
	e.moveRotors()

	letterIndex := CharToIndex(letter)
	letterIndex = e.Plugboard[letterIndex]

	for i := len(e.Rotors) - 1; i >= 0; i-- {
		letterIndex = e.Rotors[i].Step(letterIndex, false)
	}

	letterIndex = e.Reflector.Sequence[letterIndex]

	for i := 0; i < len(e.Rotors); i++ {
		letterIndex = e.Rotors[i].Step(letterIndex, true)
	}

	letterIndex = e.Plugboard[letterIndex]
	letter = IndexToChar(letterIndex)

	return letter
}

// EncodeString encodes a string.
func (e *Enigma) EncodeString(text string) string {
	var result bytes.Buffer
	for i := range text {
		result.WriteByte(e.EncodeChar(text[i]))
	}
	return result.String()
}
