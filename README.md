# Enigma cipher machine emulator

![](https://www.dropbox.com/s/5wb3u29ybxrzphl/Screenshot%202016-11-25%2015.34.47.png?dl=1)

This is a neat little emulator of the Enigma machines with simple CLI and a lot of
confugurable parameters. Somebody hurt your feelings by saying "my grandmother
encrypts better than you"? I've got you covered! With this amazing port of the
40's technology you'll be just as good at encrypting things as anyone's grandmother.

### Configuration

There's a bunch of things that can be configured when encrypting/decrypting your
text in CLI:

* Rotor set: rotors from M3 and M4, the most famous Enigma machines, are
  pre-loaded.

* Reflector: reflectors A, B, C (as well as thin B and C versions for M4) are
  supported.

* Plugboard: any number of letter pairs is accepted. Plugboard configuration
  is optional.

* Ring offsets and starting position of the rotors.

M3 and M4 can be fully emulated with the right parameters, and if it's not enough, new 
rotors and reflectors can be added quite easily: just add a new entry to the list in 
`rotors.go`, and that's it. You can also specify notches for turnover.

There's a bunch of more exotic Enigma variants and implementations, as well as devices
such as Uhr, that are not supported in this version, but pull requests are always 
welcome.

### Further reading

A bunch of material on Enigma machines, in no particular order. Explanations, specs,
other emulators.

- http://users.telenet.be/d.rijmenants/en/enigmatech.htm
- http://www.codesandciphers.org.uk/enigma/index.htm
- http://www.codesandciphers.org.uk/enigma/rotorspec.htm
- http://kerryb.github.io/enigma/
- http://enigma.louisedade.co.uk/enigma.html
- http://people.physik.hu-berlin.de/~palloks/js/enigma/enigma-u_v20_en.html
