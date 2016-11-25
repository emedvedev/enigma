package main

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/emedvedev/enigma"
	"github.com/mkideal/cli"
)

// ValidatePlugboard checks that all plugboard pairs are formatted correctly,
// and letters in pairs do not repeat.
func ValidatePlugboard(argv *CLIOpts, ctx *cli.Context) error {
	var plugboard string
	for _, pair := range argv.Plugboard {
		if matched, _ := regexp.MatchString(`^[A-Z]{2}$`, pair); !matched {
			return fmt.Errorf(
				`plugboard should be grouped by letter pairs ("AB CD"), got "%s"`,
				ctx.Color().Yellow(pair))
		}
		if strings.ContainsAny(pair, plugboard) || pair[0] == pair[1] {
			return fmt.Errorf(
				`letters cannot repeat across the plugboard, check "%s"`,
				ctx.Color().Yellow(pair))
		}
		plugboard += pair
	}
	return nil
}

// ValidateRotors checks that the requested rotors are present
// in the pre-defined list.
func ValidateRotors(argv *CLIOpts, ctx *cli.Context) error {
	for _, rotor := range argv.Rotors {
		if _, ok := enigma.Rotors[rotor]; !ok {
			return fmt.Errorf(`unknown rotor "%s"`, ctx.Color().Yellow(rotor))
		}
	}
	return nil
}

// ValidateReflector checks that the requested reflector is present
// in the pre-defined list.
func ValidateReflector(argv *CLIOpts, ctx *cli.Context) error {
	if _, ok := enigma.Reflectors[argv.Reflector]; !ok {
		return fmt.Errorf(`unknown reflector "%s"`, ctx.Color().Yellow(argv.Reflector))
	}
	return nil
}

// ValidatePosition checks that the rotor positions are in the right
// range and format.
func ValidatePosition(argv *CLIOpts, ctx *cli.Context) error {
	for _, char := range argv.Position {
		if matched, _ := regexp.MatchString(`^[A-Z]$`, char); !matched {
			return fmt.Errorf(
				`rotor positions should be single letters in the A-Z range, got "%s"`,
				ctx.Color().Yellow(char))
		}
	}
	return nil
}

// ValidateRings checks that the rotor rings are in the right
// range and format.
func ValidateRings(argv *CLIOpts, ctx *cli.Context) error {
	for _, ring := range argv.Rings {
		if ring < 1 || ring > 26 {
			return fmt.Errorf(
				`ring out of range: must be 1-26, got "%s"`,
				ctx.Color().Yellow(ring))
		}
	}
	return nil
}

// ValidateUniformity checks that the number of rotors, positions,
// and rings is equal.
func ValidateUniformity(argv *CLIOpts, ctx *cli.Context) error {
	if !(len(argv.Rotors) == len(argv.Position) && len(argv.Position) == len(argv.Rings)) {
		return fmt.Errorf(
			"number of configured rotors, rings, and position settings should be equal")
	}
	return nil
}
