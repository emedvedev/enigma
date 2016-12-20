package enigma

// Reflector is used to reverse a signal inside the Enigma: the current
// goes from the keys through the rotors to the reflector, then it is
// reversed and goes through the rotors again in the opposite direction.
type Reflector struct {
	ID       string
	Sequence [26]int
}

// NewReflector is a constuctor, taking a reflector mapping and
// its ID (name).
func NewReflector(mapping string, id string) *Reflector {
	var seq [26]int
	for i, value := range mapping {
		seq[i] = CharToIndex(byte(value))
	}
	return &Reflector{id, seq}
}

// Reflectors is a simple list of reflector pointers.
type Reflectors []Reflector

// GetByID takes a "name" of the reflector (e.g. "B") and returns the
// Reflector pointer.
func (refs *Reflectors) GetByID(id string) *Reflector {
	for _, ref := range *refs {
		if ref.ID == id {
			return &ref
		}
	}
	return nil
}
