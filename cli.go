package main

import (
	"github.com/urfave/cli"
	"os"
	"regexp"
	"strconv"
	"strings"
	"text/template"
)

var outputTemplate, _ = template.New("enigma-cli").Parse(`
Enigma configuration
-------------------------
Rotors: {{ .Context.GlobalString "rotors" }}
Rings: {{ .Context.GlobalString "rings" }}
Plugboard: {{ or (.Context.GlobalString "plugboard") ("empty") }}
Rotor position: {{ .Context.GlobalString "position" }}
Reflector: {{ .Context.GlobalString "reflector" }}
=========================
{{- if ne (.Plain) (.Original) }}

Original plaintext
-------------------------
{{ .Original }}
=========================
{{- end }}

Plaintext
-------------------------
{{ .Plain }}
=========================

Result
-------------------------
{{ .Encrypted }}
=========================

`)

func EnigmaApp() *cli.App {

	var rotors, rings, plugboard, position, reflector string

	var App = cli.NewApp()
	App.Name = "enigma"
	App.Usage = "Encrypt or decrypt text using a given Enigma configuration."
	App.Version = "0.1.0"
	App.UsageText = `enigma <text to encrypt/decrypt> [--rotors "I II III"] [--rings "1 1 1"] [--plugboard "AB CD"] [--position "A A A"] [--reflector C]`
	App.HideVersion = true
	App.Authors = []cli.Author{
		cli.Author{
			Name:  "Edward Medvedev",
			Email: "edward.medvedev@gmail.com",
		},
	}
	App.Flags = []cli.Flag{
		cli.StringFlag{
			Name:        "rotors",
			Value:       "I II III",
			Usage:       "Rotor configuration. Pre-defined rotors: I, II, III, IV, V, VI, VII, VIII, Beta, Gamma",
			Destination: &rotors,
		},
		cli.StringFlag{
			Name:        "rings",
			Value:       "1 1 1",
			Usage:       "Ring configuration, 1-26 for each rotor",
			Destination: &rings,
		},
		cli.StringFlag{
			Name:        "plugboard",
			Value:       "",
			Usage:       "Optional plugboard pairs",
			Destination: &plugboard,
		},
		cli.StringFlag{
			Name:        "position",
			Value:       "A A A",
			Usage:       "Starting position, A-Z for each rotor",
			Destination: &position,
		},
		cli.StringFlag{
			Name:        "reflector",
			Value:       "C",
			Usage:       "Reflector. Supported: A, B, `C`, B-Thin, C-Thin",
			Destination: &reflector,
		},
	}

	App.Action = cli.ActionFunc(
		func(c *cli.Context) error {

			splitRotors := strings.Split(rotors, " ")
			splitRings := strings.Split(rings, " ")
			splitPosition := strings.Split(position, " ")

			allowedRegexp := regexp.MustCompile(`[^A-Z]`)
			originalPlaintext := strings.Join(c.Args(), " ")
			plaintext := strings.TrimSpace(originalPlaintext)
			plaintext = strings.Replace(plaintext, " ", "X", -1)
			plaintext = strings.ToUpper(plaintext)
			plaintext = allowedRegexp.ReplaceAllString(plaintext, "")

			config := make([]RotorConfig, len(splitRotors))
			for index, rotor := range splitRotors {
				ring, _ := strconv.Atoi(splitRings[index])
				value := rune(splitPosition[index][0])
				config[index] = RotorConfig{rotor, value, ring}
			}
			e := NewEnigma(
				config,
				c.String("reflector"),
				strings.Split(plugboard, " "),
			)
			encrypted := e.EncryptString(plaintext)
			outputTemplate.Execute(os.Stdout, struct {
				Context   *cli.Context
				Original  string
				Plain     string
				Encrypted string
			}{c, originalPlaintext, plaintext, encrypted})

			return nil
		})

	return App
}
