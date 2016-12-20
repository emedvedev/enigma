package enigma

// HistoricRotors match the original Enigma configurations, including the
// notches. "Beta" and "Gamma" are additional rotors used in M4
// at the leftmost position.
var HistoricRotors = Rotors{
	*NewRotor("EKMFLGDQVZNTOWYHXUSPAIBRCJ", "I", "Q"),
	*NewRotor("AJDKSIRUXBLHWTMCQGZNPYFVOE", "II", "E"),
	*NewRotor("BDFHJLCPRTXVZNYEIWGAKMUSQO", "III", "V"),
	*NewRotor("ESOVPZJAYQUIRHXLNFTGKDCMWB", "IV", "J"),
	*NewRotor("VZBRGITYUPSDNHLXAWMJQOFECK", "V", "Z"),
	*NewRotor("JPGVOUMFYQBENHZRDKASXLICTW", "VI", "ZM"),
	*NewRotor("NZJHGRCXMYSWBOUFAIVLPEKQDT", "VII", "ZM"),
	*NewRotor("FKQHTLXOCBJSPDZRAMEWNIUYGV", "VIII", "ZM"),
	*NewRotor("LEYJVCNIXWPBQMDRTAKZGFUHOS", "Beta", ""),
	*NewRotor("FSOKANUERHMBTIYCWLQPZXVGJD", "Gamma", ""),
}

// HistoricReflectors in the list are pre-loaded with historically accurate data
// from Enigma machines. Use "B-Thin" and "C-Thin" with M4 (4 rotors).
var HistoricReflectors = Reflectors{
	*NewReflector("EJMZALYXVBWFCRQUONTSPIKHGD", "A"),
	*NewReflector("YRUHQSLDPXNGOKMIEBFZCWVJAT", "B"),
	*NewReflector("FVPJIAOYEDRZXWGCTKUQSBNMHL", "C"),
	*NewReflector("ENKQAUYWJICOPBLMDXZVFTHRGS", "B-thin"),
	*NewReflector("RDOBJNTKVEHMLFCWZAXGYIPSUQ", "C-thin"),
}
