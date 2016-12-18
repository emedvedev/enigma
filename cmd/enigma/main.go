package main

import (
	"fmt"
	"os"
	"strings"
	"text/template"

	"github.com/emedvedev/enigma"
	"github.com/mkideal/cli"
)

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

// CLIDefaults is used to populate default values in case
// one or more of the parameters aren't set. It is assumed
// that rotor rings and positions will be the same for all
// rotors if not set explicitly, so only one value is stored.
var CLIDefaults = struct {
	Reflector string
	Ring      int
	Position  string
	Rotors    []string
}{
	Reflector: "B",
	Ring:      1,
	Position:  "A",
	Rotors:    []string{"I", "II", "III"},
}

// SetDefaults sets values for all Enigma parameters that
// were not set explicitly.
// Plugboard is the only parameter that does not require a
// default, since it may not be set, and in some Enigma versions
// there was no plugboard at all.
func SetDefaults(argv *CLIOpts) {
	if argv.Reflector == "" {
		argv.Reflector = CLIDefaults.Reflector
	}
	if len(argv.Rotors) == 0 {
		argv.Rotors = CLIDefaults.Rotors
	}
	loadRings := (len(argv.Rings) == 0)
	loadPosition := (len(argv.Position) == 0)
	if loadRings || loadPosition {
		for range argv.Rotors {
			if loadRings {
				argv.Rings = append(argv.Rings, CLIDefaults.Ring)
			}
			if loadPosition {
				argv.Position = append(argv.Position, CLIDefaults.Position)
			}
		}
	}
}

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
			config[index] = enigma.RotorConfig{ID: rotor, Start: value, Ring: ring}
		}

		e := enigma.NewEnigma(config, argv.Reflector, argv.Plugboard)
		encoded := e.EncodeString(plaintext)

		if argv.Condensed {
			fmt.Print(encoded)
			return nil
		}

		tmpl, _ := template.New("cli").Parse(OutputTemplate)
		err := tmpl.Execute(os.Stdout, struct {
			Original, Plain, Encoded string
			Args                     *CLIOpts
			Ctx                      *cli.Context
		}{originalPlaintext, plaintext, encoded, argv, ctx})
		return err

	})

}
