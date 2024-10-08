![build](https://github.com/transcriptaze/midiasm/workflows/build/badge.svg)

# midiasm

MIDI assembler/disassembler to convert between standard MIDI files and a text/JSON equivalent.

## Raison d'être

A Go reimplementation of Jeff Glatt's long defunct Windows-only MIDIASM (last seen archived at [MIDI Technical Fanatic's Brainwashing Center](http://midi.teragonaudio.com)). Because sometimes it's easier to just programmatically deal with
text or JSON.

## Installation

Executables for all the supported operating systems are packaged in the [releases](https://github.com/transcriptaze/midiasm/releases). Installation is straightforward - download the archive and extract it to a directory of your choice. 

`midiasm help` will list the available commands and associated options (documented below).

### Building from source

Required tools:
- [Go 1.23+](https://go.dev)
- _make_ (optional but recommended)

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

_None_

## midiasm

Usage: ```midiasm <command> <options>```

Supported commands:

- `help`
- `version`
- [`disassemble`](#disassemble)
- [`assemble`](#assemble)
- [`export`](#export)
- [`notes`](#notes)
- [`click`](#click)
- [`transpose`](#transpose)
- [`tsv`](#tsv)

Defaults to `disassemble` if the command is not provided.

### `disassemble`

Disassembles a MIDI file and displays the tracks in a human readable format.

Command line:

` midiasm [--debug] [--verbose] [--C4] [--split] [--out <file>] <MIDI file>`

```
  --out <file>  Writes the disassembly to a file. Default is to write to stdout.
  --split       Writes each track to a separate file. Default is `false`.

  Options:

  --C4       Uses C4 as middle C (Yamaha convention). Defaults to C3.
  --debug    Displays internal information while processing a MIDI file. Defaults to false
  --verbose  Enables 'verbose' logging. Defaults to false

  Example:

  midiasm --debug --verbose --out one-time.txt one-time.mid
```

### `assemble`

Assembles a MIDI file from a text or JSON source.

Command line:

` midiasm assemble [--debug] [--verbose] [--C4] [--out <MIDI file>] <file>`

```
  --out <file>  Output MIDI file. Defaults to the input file with a .midi extension.

  Options:

  --C4       Uses C4 as middle C (Yamaha convention). Defaults to C3.
  --debug    Displays internal information while processing a MIDI file. Defaults to false
  --verbose  Enables 'verbose' logging. Defaults to false

  Example:

  midiasm assemble --debug --verbose --out one-time.mid one-time.json
```

### `export`

Extracts the MIDI information as JSON for use with other tools (e.g. _jq_).

Command line:

` midiasm export [--debug] [--verbose] [--C4] [--out <file>] <MIDI file>`

```
  --out <file>     Writes the JSON to a file. Default is to write to stdout.
  --json           Formats the output as JSON - the default is human readable text.
  --transpose <N>  Transposes the notes up or down by N semitones.

  Options:

  --C4       Uses C4 as middle C (Yamaha convention). Defaults to C3.
  --debug    Displays internal information while processing a MIDI file. Defaults to false
  --verbose  Enables 'verbose' logging. Defaults to false

  Example:

  midiasm notes --debug --verbose --out one-time.json one-time.mid
```


### `notes`

Extracts the _NoteOn_ and _NoteOff_ events to generate a list of notes with start times and durations.

Command line:

` midiasm notes [--debug] [--verbose] [--C4] [--out <file>] <MIDI file>`

```
  --out <file>  Writes the notes to a file. Default is to write to stdout.

  Options:

  --C4       Uses C4 as middle C (Yamaha convention). Defaults to C3.
  --debug    Displays internal information while processing a MIDI file. Defaults to false
  --verbose  Enables 'verbose' logging. Defaults to false

  Example:

  midiasm notes --debug --verbose --out one-time.notes one-time.mid
```

### `click`

Extracts the _beats_ from the MIDI file in a format that can be used to create a click track.

Command line:

` midiasm click [--debug] [--verbose] [--C4] [--out <file>] <MIDI file>`

```
  --out <file>  Writes the click track to a file. Default is to write to stdout.

  Options:

  --C4       Uses C4 as middle C (Yamaha convention). Defaults to C3.
  --debug    Displays internal information while processing a MIDI file. Defaults to false
  --verbose  Enables 'verbose' logging. Defaults to false

  Example:
  
  midiasm click --debug --verbose --out one-time.click one-time.mid
```

### `transpose`

Transposes the key of the notes (and key signature) and writes it back as MIDI file.

Command line:

` midiasm transpose [--debug] [--verbose] [--C4] --semitones <steps> --out <file> <MIDI file>`

```
  --semitones <N>  Number of semitones to transpose up or down. Defaults to 0.
  --out <file>     (required) Destination file for the transposed MIDI. 

  Options:

  --C4       Uses C4 as middle C (Yamaha convention). Defaults to C3.
  --debug    Displays internal information while processing a MIDI file. Defaults to false
  --verbose  Enables 'verbose' logging. Defaults to false

  Example:
  
  midiasm transpose --debug --verbose --semitones +5 --out one-time+5.mid one-time.mid
```

### `tsv`

Extracts the MIDI information as a TSV or fixed width file for use with other tools (e.g. [miller](https://github.com/johnkerl/miller))

Command line:

` midiasm tsv [--debug] [--verbose] [--C4] [--out <file>] <MIDI file>`

```
  --out <file>   Output filepath. Default is to write to stdout.
  --delimiter    Column delimiter for TSV files. Defaults to TAB.
  --tabular      Formats the outputs as fixed width columns

  Options:

  --C4       Uses C4 as middle C (Yamaha convention). Defaults to C3.
  --debug    Displays internal information while processing a MIDI file. Defaults to false
  --verbose  Enables 'verbose' logging. Defaults to false

  Example:

  midiasm tsv --debug --verbose --out one-time.tsv one-time.mid
```

## Tools

1. [jq](https://jqlang.github.io/jq)
2. [miller](https://github.com/johnkerl/miller)


## References

1. [The Complete MIDI 1.0 Detailed Specification](https://www.midi.org/specifications/item/the-midi-1-0-specification)
2. [Somascape - MIDI Files Specification](http://www.somascape.org/midi/tech/mfile.html)
3. [(archive) MIDI Technical Fanatic's Brainwashing Center](http://midi.teragonaudio.com)
4. [StackExchange::Music Transposing key signatures - how to do so quickly?](https://music.stackexchange.com/questions/110078/transposing-key-signatures-how-to-do-so-quickly)
5. [An Introduction to MIDI](https://www.midi.org/articles/an-intro-to-midi)
6. [Roland: Midi Guidebook](https://ia902502.us.archive.org/3/items/manuallib-id-2725326/2725326_text.pdf)
7. [Timing in MIDI files](https://sites.uci.edu/camp2014/2014/05/19/timing-in-midi-files)


## Alternatives
1.  [mido](https://github.com/mido/mido)
2.  [midicsv](https://www.fourmilab.ch/webtools/midicsv)
3.  [Werkmeister](https://werckme.github.io/werckmeister)
4.  [pymidi](https://pypi.org/project/py-midi)
5.  [DryWetMIDI](https://github.com/melanchall/drywetmidi)
6.  [Online MIDI Editor](http://midi.mathewvp.com)
7.  [go-midi](https://github.com/moutend/go-midi)
8.  [midiutil](https://pypi.org/project/MIDIUtil)
9.  [music21](https://web.mit.edu/music21)
10. [fq](https://github.com/wader/fq)

