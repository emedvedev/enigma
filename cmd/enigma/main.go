package main

import (
	"fmt"
	"os"
	"strings"
	"text/template"

	"github.com/emedvedev/enigma"
	"github.com/mkideal/cli"
)

// DescriptionTemplate is a simple template for help and usage.
const DescriptionTemplate = `
usage: enigma <text> [--rotors=I II III] [--rings=3 4 3] [--reflector=C]
                     [--plugboard=AB CD] [--position=A A A]

Enigma cipher machine emulator

Encrypt all the things with the power of this handy Enigma emulator:
it fully supports the most popular M3 and M4 models, and quite a few
others, too. Choose a rotor set and a reflector, configure rings and
starting positions of the rotors, select plugboard pairs--and you're
all set! Don't forget: cryptography is only real when shared. Make a
friend.

Enjoy!
`

// OutputTemplate is a template for the encryption result
// that will be used if the Concise flag isn't set.
const OutputTemplate = `
{{ (.Ctx.Color).Bold "Result:" }}
  {{ .Encrypted }}

{{ (.Ctx.Color).Bold "Enigma configuration:" }}
  Rotors: {{ .Args.Rotors }}
  Rotor positions: {{ .Args.Position }}
  Rings: {{ .Args.Rings }}
  Plugboard: {{ or (.Args.Plugboard) ("empty") }}
  Reflector: {{ .Args.Reflector }}

{{ (.Ctx.Color).Bold "Original text:" }}
  {{ .Original }}
{{ if ne (.Plain) (.Original) }}
{{ (.Ctx.Color).Bold "Processed original text:" }}
  {{ .Plain }}
{{ end }}
`

func main() {

	cli.SetUsageStyle(cli.DenseManualStyle)
	cli.Run(new(CLIOpts), func(ctx *cli.Context) error {
		argv := ctx.Argv().(*CLIOpts)
		originalPlaintext := strings.Join(ctx.Args(), " ")
		plaintext := enigma.SanitizePlaintext(originalPlaintext)

		if argv.Help || len(plaintext) == 0 {
			com := ctx.Command()
			com.Text = DescriptionTemplate
			ctx.String(com.Usage(ctx))
			return nil
		}

		config := make([]enigma.RotorConfig, len(argv.Rotors))
		for index, rotor := range argv.Rotors {
			ring := argv.Rings[index]
			value := argv.Position[index][0]
			config[index] = enigma.RotorConfig{ID: rotor, Start: byte(value), Ring: ring}
		}

		e := enigma.NewEnigma(config, argv.Reflector, argv.Plugboard)
		encrypted := e.EncryptString(plaintext)

		if argv.Condensed {
			fmt.Print(encrypted)
			return nil
		}

		tmpl, _ := template.New("cli").Parse(OutputTemplate)
		err := tmpl.Execute(os.Stdout, struct {
			Original, Plain, Encrypted string
			Args                       *CLIOpts
			Ctx                        *cli.Context
		}{originalPlaintext, plaintext, encrypted, argv, ctx})
		return err

	})

}

// CLIOpts sets the parameter format for Enigma CLI. It also includes a "help"
// flag and a "condensed" flag telling the program to output plain result.
// Also, this CLI module abuses tags so much it hurts. Oh well. ¯\_(ツ)_/¯
type CLIOpts struct {
	Help      bool `cli:"!h,help" usage:"Show help."`
	Condensed bool `cli:"c,condensed" name:"false" usage:"Output the result without additional information."`

	Rotors    []string `cli:"rotors" name:"I II III" usage:"Rotor configuration. Supported: I, II, III, IV, V, VI, VII, VIII, Beta, Gamma."`
	Rings     []int    `cli:"rings" name:"1 1 1" usage:"Rotor rings offset: from 1 (default) to 26 for each rotor."`
	Position  []string `cli:"position" name:"A A A" usage:"Starting position of the rotors: from A (default) to Z for each."`
	Plugboard []string `cli:"plugboard" name:"[]" usage:"Optional plugboard pairs to scramble the message further."`

	Reflector string `cli:"reflector" name:"C" usage:"Reflector. Supported: A, B, C, B-Thin, C-Thin."`
}

// Validate runs checks on all available Enigma parameters; the checks
// themselves are separate. The defaults are loaded before validation:
// some of the parameter combinations that a user might supply won't work
// with the defaults, so we'll combine first, then check the final form.
func (argv *CLIOpts) Validate(ctx *cli.Context) error {
	SetDefaults(argv)
	validators := [](func(argv *CLIOpts, ctx *cli.Context) error){
		ValidatePlugboard,
		ValidateRotors,
		ValidateReflector,
		ValidatePosition,
		ValidateRings,
		ValidateUniformity,
	}
	for _, validator := range validators {
		if err := validator(argv, ctx); err != nil {
			return err
		}
	}
	return nil
}

// SetDefaults sets values for all Enigma parameters that
// were not set explicitly. Plugboard is the only parameter
// that does not require a default, since it can be empty,
// and some Enigma machines did not have a plugboard at all.
func SetDefaults(argv *CLIOpts) {
	if argv.Reflector == "" {
		argv.Reflector = EnigmaDefaults.Reflector
	}
	if len(argv.Rotors) == 0 {
		argv.Rotors = EnigmaDefaults.Rotors
	}
	loadRings := (len(argv.Rings) == 0)
	loadPosition := (len(argv.Position) == 0)
	if loadRings || loadPosition {
		for range argv.Rotors {
			if loadRings {
				argv.Rings = append(argv.Rings, EnigmaDefaults.Ring)
			}
			if loadPosition {
				argv.Position = append(argv.Position, EnigmaDefaults.Position)
			}
		}
	}
}

// EnigmaDefaults is used to populate default values in case
// one or more of the parameters aren't set. It is assumed
// that rotor rings and positions will be the same for all
// rotors if not set explicitly, so only one value is stored.
var EnigmaDefaults = struct {
	Reflector string
	Ring      int
	Position  string
	Rotors    []string
}{
	Reflector: "C",
	Ring:      1,
	Position:  "A",
	Rotors:    []string{"I", "II", "III"},
}
