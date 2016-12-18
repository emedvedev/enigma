package main

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

// OutputTemplate is a template for the encoding result
// that will be used if the Condensed flag isn't set.
const OutputTemplate = `
{{ (.Ctx.Color).Bold "Original text:" }}
  {{ .Original }}
{{ if ne (.Plain) (.Original) }}
{{ (.Ctx.Color).Bold "Processed original text:" }}
  {{ .Plain }}
{{ end }}
{{ (.Ctx.Color).Bold "Enigma configuration:" }}
  Rotors: {{ .Args.Rotors }}
  Rotor positions: {{ .Args.Position }}
  Rings: {{ .Args.Rings }}
  Plugboard: {{ or (.Args.Plugboard) ("empty") }}
  Reflector: {{ .Args.Reflector }}

{{ (.Ctx.Color).Bold "Result:" }}
  {{ .Encoded }}
`
