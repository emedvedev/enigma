package enigma

import (
	"strings"
)

// Rotor is the device performing letter substitutions inside
// the Enigma machine. Rotors can be put in different positions,
// swapped, and replaced; they are also rotated during the encoding
// process, following the machine configuration. As a result, there
// are billions of possible combinations, making brute-forcing attacks
// on Enigma unfeasible.
type Rotor struct {
	ID          string
	Sequence    string
	IntSequence []int
	Notch       []int

	Offset int
	Ring   int
}

// NewRotor is a constructor taking a mapping string and a notch position
// that will trigger the rotation of the next rotor.
func NewRotor(mapping string, id string, notches string) *Rotor {
	notchArr := make([]int, len(notches))
	for i := range notches {
		notchArr[i] = ToInt(notches[i])
	}
	intSeq := make([]int, 26)
	for i, value := range mapping {
		intSeq[i] = ToInt(byte(value))
	}
	return &Rotor{id, mapping, intSeq, notchArr, 0, 0}
}

// IsAtNotch ...
func (r *Rotor) IsAtNotch() bool {
	for _, notch := range r.Notch {
		if r.Offset == notch {
			return true
		}
	}
	return false
}

// Step through the rotor, performing the letter substitution depending
// on the offset and direction.
func (r *Rotor) Step(letter *int, invert bool) {
	*letter = (*letter - r.Ring + r.Offset + 26) % 26
	if invert {
		*letter = strings.IndexByte(r.Sequence, ToChar(*letter))
	} else {
		*letter = r.IntSequence[*letter]
	}
	*letter = (*letter + r.Ring - r.Offset + 26) % 26
}

// RotorConfig sets the initial configuration for a rotor: ID from
// the pre-defined list, a starting position (A-Z), and a ring setting (1-26).
type RotorConfig struct {
	ID    string
	Start byte
	Ring  int
}

// Rotors ...
type Rotors []*Rotor

// Reflectors ...
type Reflectors []*Reflector

// HistoricRotors match the original Enigma configurations, including the
// notches. "Beta" and "Gamma" are additional rotors used in M4
// at the 4th position.
var HistoricRotors = Rotors{
	NewRotor("EKMFLGDQVZNTOWYHXUSPAIBRCJ", "I", "Q"),
	NewRotor("AJDKSIRUXBLHWTMCQGZNPYFVOE", "II", "E"),
	NewRotor("BDFHJLCPRTXVZNYEIWGAKMUSQO", "III", "V"),
	NewRotor("ESOVPZJAYQUIRHXLNFTGKDCMWB", "IV", "J"),
	NewRotor("VZBRGITYUPSDNHLXAWMJQOFECK", "V", "Z"),
	NewRotor("JPGVOUMFYQBENHZRDKASXLICTW", "VI", "ZM"),
	NewRotor("NZJHGRCXMYSWBOUFAIVLPEKQDT", "VII", "ZM"),
	NewRotor("FKQHTLXOCBJSPDZRAMEWNIUYGV", "VIII", "ZM"),
	NewRotor("LEYJVCNIXWPBQMDRTAKZGFUHOS", "Beta", ""),
	NewRotor("FSOKANUERHMBTIYCWLQPZXVGJD", "Gamma", ""),
}

// GetByID ...
func (rs *Rotors) GetByID(id string) *Rotor {
	for _, rotor := range *rs {
		if rotor.ID == id {
			return rotor
		}
	}
	return nil
}

// Reflector is used to reverse a signal inside the Enigma: the current
// goes from the keys through the rotors to the reflector, then it is
// reversed and goes through the rotors again in the opposite direction.
type Reflector struct {
	ID          string
	Sequence    string
	IntSequence []int
}

// Reflect is a method for reversing the Enigma signal in a reflector:
// it is just a simple substitution, essentially.
func (ref *Reflector) Reflect(letter int) int {
	return ref.IntSequence[letter]
}

// HistoricReflectors in the list are pre-loaded with historically accurate data
// from Enigma machines. Use "B-Thin" and "C-Thin" with M4 (4 rotors).
var HistoricReflectors = Reflectors{
	NewReflector("EJMZALYXVBWFCRQUONTSPIKHGD", "A"),
	NewReflector("YRUHQSLDPXNGOKMIEBFZCWVJAT", "B"),
	NewReflector("FVPJIAOYEDRZXWGCTKUQSBNMHL", "C"),
	NewReflector("ENKQAUYWJICOPBLMDXZVFTHRGS", "B-thin"),
	NewReflector("RDOBJNTKVEHMLFCWZAXGYIPSUQ", "C-thin"),
}

// NewReflector ...
func NewReflector(mapping string, id string) *Reflector {
	intSeq := make([]int, 26)
	for i, value := range mapping {
		intSeq[i] = ToInt(byte(value))
	}
	return &Reflector{id, mapping, intSeq}
}

// GetByID ...
func (refs *Reflectors) GetByID(id string) *Reflector {
	for _, ref := range *refs {
		if ref.ID == id {
			return ref
		}
	}
	return nil
}
