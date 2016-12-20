package enigma

// Rotor is the device performing letter substitutions inside
// the Enigma machine. Rotors can be put in different positions,
// swapped, and replaced; they are also rotated during the encoding
// process, following the machine configuration. As a result, there
// are millions of possible combinations, making brute-forcing attacks
// on Enigma unfeasible (and even more so when the plugboard is used).
type Rotor struct {
	ID          string
	StraightSeq [26]int
	ReverseSeq  [26]int
	Turnover    []int

	Offset int
	Ring   int
}

// NewRotor is a constructor for rotors, taking a mapping string
// and a turnover position.
func NewRotor(mapping string, id string, turnovers string) *Rotor {
	r := &Rotor{ID: id, Offset: 0, Ring: 0}
	r.Turnover = make([]int, len(turnovers))
	for i := range turnovers {
		r.Turnover[i] = CharToIndex(turnovers[i])
	}
	for i, letter := range mapping {
		index := CharToIndex(byte(letter))
		r.StraightSeq[i] = index
		r.ReverseSeq[index] = i
	}
	return r
}

// Move the rotor, shifting the offset by a given number.
func (r *Rotor) move(offset int) {
	r.Offset = (r.Offset + offset) % 26
}

// ShouldTurnOver checks if the current rotor position corresponds
// to a notch that is supposed to move the next rotor.
func (r *Rotor) ShouldTurnOver() bool {
	for _, turnover := range r.Turnover {
		if r.Offset == turnover {
			return true
		}
	}
	return false
}

// Step through the rotor, performing the letter substitution depending
// on the offset and direction.
func (r *Rotor) Step(letter int, invert bool) int {
	letter = (letter - r.Ring + r.Offset + 26) % 26
	if invert {
		letter = r.ReverseSeq[letter]
	} else {
		letter = r.StraightSeq[letter]
	}
	letter = (letter + r.Ring - r.Offset + 26) % 26
	return letter
}

// Rotors is a simple list of rotor pointers.
type Rotors []Rotor

// GetByID takes a "name" of the rotor (e.g. "III") and returns the
// Rotor pointer.
func (rs *Rotors) GetByID(id string) *Rotor {
	for _, rotor := range *rs {
		if rotor.ID == id {
			return &rotor
		}
	}
	return nil
}
