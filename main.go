package main

import (
	"fmt"
	"github.com/emedvedev/enigma/enigma"
	"github.com/urfave/cli"
	"os"
	"strconv"
	"strings"
)

// TODO: tests
// TODO: readme and docs
// TODO: CLI
// http://people.physik.hu-berlin.de/~palloks/js/enigma/enigma-u_v20_en.html
func main() {

	app := cli.NewApp()
	app.Name = "enigma"
	app.Usage = "Encrypt text using a given Enigma configuration."
	app.Version = "0.1.0"
	app.Authors = []cli.Author{
		cli.Author{
			Name:  "Edward Medvedev",
			Email: "edward.medvedev@gmail.com",
		},
	}
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "rotors",
			Usage: "Rotor configuration: 3 for M3 emulation, 4 for M4 emulation. Supported: I, II, III, IV, V, VI, VII, VIII, Beta, Gamma.",
		},
		cli.StringFlag{
			Name:  "rings",
			Usage: "Ring configuration, 1-26 for each rotor.",
		},
		cli.StringFlag{
			Name:  "values",
			Usage: "Ring configuration, 1-26 for each rotor.",
		},
		cli.StringFlag{
			Name:  "reflector",
			Value: "C",
			Usage: "Reflector. Supported: A, B, `C`, B-Thin, C-Thin",
		},
		cli.StringFlag{
			Name:  "plugboard",
			Usage: "Optional plugboard pairs: `\"AB,CD,EF\"`, etc.",
		},
	}
	app.Action = func(c *cli.Context) error {

		rotors := strings.Split(c.String("rotors"), ",")
		rings := strings.Split(c.String("rings"), ",")
		values := strings.Split(c.String("values"), ",")
		plugboard := strings.Split(c.String("plugboard"), ",")
		reflector := c.String("reflector")

		fmt.Println("Rotors: ", rotors)
		fmt.Println("Rings: ", rings)
		fmt.Println("Starting values: ", values)
		fmt.Println("Plugboard: ", plugboard)
		fmt.Println("Reflector: ", reflector)

		config := make([]enigma.RotorConfig, len(rotors))
		for index, rotor := range rotors {
			ring, _ := strconv.Atoi(rings[index])
			value := rune(values[index][0])
			config[index] = enigma.RotorConfig{rotor, value, ring}
		}

		mappedPlugboard := make([][2]rune, len(plugboard))

		for index, pair := range plugboard {
			if len(pair) > 0 {
				mappedPlugboard[index] = [2]rune{rune(pair[0]), rune(pair[1])}
			}
		}

		plaintext := c.Args().Get(0)
		fmt.Println("Your text: ", plaintext)

		enigma := enigma.NewEnigma(
			config,
			c.String("reflector"),
			mappedPlugboard,
		)
		fmt.Print("Encrypted: ")
		for index := range plaintext {
			char := rune(plaintext[index])
			enigma.EncryptChar(&char)
			fmt.Print(string(char))
		}
		return nil
	}

	app.Run(os.Args)
}
