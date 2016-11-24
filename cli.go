package main

import (
	"fmt"
	"github.com/mkideal/cli"
	"os"
	"regexp"
	"strings"
	"text/template"
)

const DESCRIPTION = `
usage: enigma <text> [--rotors=I II III] [--rings=3 4 3] [--reflector=C]
                     [--plugboard=AB CD] [--position=A A A]

Enigma cipher machine simulator

Encrypt or decrypt text using this configurable Enigma simulation:
it fully supports both M3 and M4, as well as some older models.
You are free to pick a rotor set, configure rings and starting
positions of the rotors, choose from a variety of reflectors and
mess with the plugboard.

Enjoy!
`

const OUTPUT = `
Enigma configuration
--------------------------------------------------
Rotors: {{ .Args.Rotors }}
Rotor positions: {{ .Args.Position }}
Rings: {{ .Args.Rings }}
Plugboard: {{ or (.Args.Plugboard) ("empty") }}
Reflector: {{ .Args.Reflector }}
==================================================

Plaintext
--------------------------------------------------
{{ .Original }}
==================================================
{{ if ne (.Plain) (.Original) }}
Modified plaintext
--------------------------------------------------
{{ .Plain }}
==================================================
{{ end }}
Result
--------------------------------------------------
{{ .Encrypted }}
==================================================

`

type CLIOpts struct {
	Help      bool     `cli:"!h,help" usage:"Show help."`
	Rotors    []string `cli:"rotors" name:"I II III" usage:"Rotor configuration. Supported: I, II, III, IV, V, VI, VII, VIII, Beta, Gamma."`
	Rings     []int    `cli:"rings" name:"3 4 3" usage:"Each rotor ring can be shifted: 1 is the default location, 26 is the maximum."`
	Plugboard []string `cli:"plugboard" name:"AB CD" usage:"Optional plugboard pairs. Letters must be unique across the plugboard."`
	Position  []string `cli:"position" name:"A A A" usage:"Starting position, A-Z for each rotor."`
	Reflector string   `cli:"reflector" name:"C" usage:"Reflector. Supported: A, B, C, B-Thin, C-Thin."`
}

func (argv *CLIOpts) Validate(ctx *cli.Context) error {

	for _, char := range argv.Position {
		if matched, _ := regexp.MatchString(`^[A-Z]$`, char); !matched {
			return fmt.Errorf("Rotor positions should be single letters in the A-Z range.")
		}
	}

	if !(len(argv.Rotors) == len(argv.Position) && len(argv.Position) == len(argv.Rings)) {
		return fmt.Errorf("You should configure equal number of rotors, rings, and position settings.")
	}

	for _, ring := range argv.Rings {
		if ring < 1 || ring > 26 {
			return fmt.Errorf("Ring out of range! Must be in the range of 1-26.")
		}
	}

	return nil
}

func EnigmaCLI() {

	cli.SetUsageStyle(cli.DenseManualStyle)
	cli.Run(new(CLIOpts), func(ctx *cli.Context) error {
		argv := ctx.Argv().(*CLIOpts)
		if argv.Help {
			com := ctx.Command()
			com.Text = DESCRIPTION
			ctx.String(com.Usage(ctx))
			return nil
		}

		originalPlaintext := strings.Join(ctx.Args(), " ")
		plaintext := SanitizePlaintext(originalPlaintext)

		config := make([]RotorConfig, len(argv.Rotors))
		for index, rotor := range argv.Rotors {
			ring := argv.Rings[index]
			value := rune(argv.Position[index][0])
			config[index] = RotorConfig{rotor, value, ring}
		}

		e := NewEnigma(config, argv.Reflector, argv.Plugboard)
		encrypted := e.EncryptString(plaintext)

		tmpl, _ := template.New("cli").Parse(OUTPUT)
		tmpl.Execute(os.Stdout, struct {
			Original, Plain, Encrypted string
			Args                       *CLIOpts
		}{originalPlaintext, plaintext, encrypted, argv})

		return nil
	})

}
