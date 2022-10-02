# midiasm

MIDI assembler/disassembler to convert between standard MIDI files and a text/JSON equivalent.

## Raison d'être

A Go reimplementation of Jeff Glatt's long defunct Windows-only MIDIASM (last seen archived at [MIDI Technical Fanatic's Brainwashing Center](http://midi.teragonaudio.com)). Because sometimes it's easier to just programmatically deal with
text or JSON.

## Releases

*In development*


| *Version* | *Description*                                                                             |
| --------- | ----------------------------------------------------------------------------------------- |
|           |                                                                                           |

## Installation

Executables for all the supported operating systems are packaged in the [releases](https://github.com/transcriptaze/midiasm/releases). Installation is straightforward - download the archive and extract it to a directory of your choice. 

`midiasm help` will list the available commands and associated options (documented below).

### Building from source

Required tools:
- [Go 1.19+](https://go.dev)
- make (optional but recommended)

To build using the included Makefile:

```
git clone https://github.com/transcriptaze/midiasm.git
cd midiasm
make build
```

Without using `make`:
```
git clone https://github.com/transcriptaze/midiasm.git
cd midiasm
go build -trimpath -o bin/ ./...
```

The above commands build the `midiasm` executable to the `bin` directory.


#### Dependencies

| *Dependency*                                                            | *Description*                        |
| ----------------------------------------------------------------------- | -------------------------------------|
|                                                                         |                                      |
|                                                                         |                                      |


## midiasm

Usage: ```midiasm <command> <options>```

Supported commands:

- `help`
- `version`
- `disassemble`
- `notes`
- `export`
- `click`
- `transpose`

Defaults to `disassemble` if the command is not provided.

### `disassemble`

Disassemble a MIDI file and displays the tracks in a human readable format.

Command line:

` midiasm [--debug] [--verbose] [--split] [--out <file>] <MIDI file>`

```
  --out <file>  Writes the disassembly to a file. Default is to write to _stdout_.
  --split       Writes each track to a separate file. Default is `false`.

  Options:

  --debug    Displays internal information while processing a MIDI file. Defaults to false
  --verbose  Enables 'verbose' logging. Defaults to false

  Example:

  uhppoted-codegen generate --models bindings/.models --templates bindings/rust --out generated/rust
```

### `export`

Generates a _models.json_ file that represents the internal UHPPOTE models used to generate the functions,
requests and responses.

Command line:

` uhppoted-codegen export [--out <file>]`

```
  --out <file> File for models JSON. Defaults to models.json.

  Example:
  
  uhppoted-codegen export --out my-models.json
```

















## Installation

### Building from source

#### Dependencies

## midiasm

## References

1. [The Complete MIDI 1.0 Detailed Specification](https://www.midi.org/specifications/item/the-midi-1-0-specification)
2. [Somascape - MIDI Files Specification](http://www.somascape.org/midi/tech/mfile.html)
3. [(archive) MIDI Technical Fanatic's Brainwashing Center](http://midi.teragonaudio.com)
4. [(github) mido](https://github.com/mido/mido)
5. [midicsv](https://www.fourmilab.ch/webtools/midicsv)
6. [StackExchange::Music Transposing key signatures - how to do so quickly?](https://music.stackexchange.com/questions/110078/transposing-key-signatures-how-to-do-so-quickly)



