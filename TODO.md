# TODO

- [ ] Fix SMPTE encoding/decoding (cf. https://github.com/transcriptaze/midiasm/issues/6)
    - [x] Fix SMPTE offset metaevent binary encoding
    - [x] Fix SMPTE offset metaevent decoding
    - [x] MThd binary unmarshal
    - [x] Rework decoder to use MThd binary unmarshal
    - [ ] Encode MThd SMPTE time divisions correctly
    - [ ] Update JSON encoding/decoding
    - [ ] Update text encoding/decoding
    - [ ] Update TSV encoding/decoding
    - (?) Use enums for SMPTE offset frame rate
    - [ ] Rework MThd struct to explicitly differentiate PPQN and SMPTE time divisions

- [ ] _humanise_ command
      - https://en.wikipedia.org/wiki/Stochastic_computing

- [ ] MIDI-2.0
      - https://github.com/midi2-dev/MIDI2.0Workbench


## NTS

- !!!! https://github.com/SuperDisk/tar.pl
- https://www.masteringemacs.org/article/combobulate-structured-movement-editing-treesitter
- https://github.com/wader/fq
- [Interval Parsing Grammars for File Format Parsing](https://dl.acm.org/doi/pdf/10.1145/3591264)
- https://www.gnu.org/software/gettext/libtextstyle/manual/libtextstyle.html#Introduction

- [ ] Fix bug in tick
```
      00 90 39 4C                           tick:23760      delta:0          90 NoteOn                 channel:0  note:A3, velocity:76
   9C 10 80 2D 40                           tick:27360      delta:3600       80 NoteOff                channel:0  note:A2, velocity:64
      00 80 34 40                           tick:27360      delta:0          80 NoteOff                channel:0  note:E3, velocity:64
      00 80 39 40                           tick:27360      delta:0          80 NoteOff                channel:0  note:A3, velocity:64
      00 FF 2F 00                           tick:27360      delta:0          2F EndOfTrack

```
- [ ] NoteOn with 0 velocity -> NoteOff
      - [x] Rework notes.go
      - [x] Fix unit tests
      - [x] Note.Duration
      - [ ] Unit tests
            - [x] Basic notes
            - [x] NoteOn with velocity 0
            - [x] Tempo changes
            - [ ] EndOfTrack
      - [ ] --use-note-0
      - [ ] --sort
      - [x] Always printing debug info
      - [x] debug output is weird
      - [ ] Calculate bar/beat
      - [x] Figure the dt error in `buildTrackEvents`
      - [ ] Rework algorithm using tick so that delta time error doesn't accumulate

- Rework everything so that SMF & MThd & MTrk are just containers i.e. all the complication around
  marshalling and unmarshalling happens in decoders.
  - [x] UnmarshalBinary
  - [ ] Rework byte reader hack

- https://music.stackexchange.com/questions/39446/where-am-i-going-wrong-in-interpreting-this-midi-string?rq=1
```
[controller] Messages 123 through 127 also function as All Notes Off messages. They will turn off all voices controlled by the assigned Basic Channel. These messages should not be sent periodically, but only for a specific purpose. In no case should they be used in lieu of Note Off commands to turn off notes which have been previously turned on. Any All Notes Off command (123-127) may by ignored by a receiver with no possibility of notes staying on, since any Note On command must have a corresponding specific Note Off command.
```

- https://music.stackexchange.com/questions/127381/duration-of-a-midi-file-by-parsing-it-and-making-a-stream-of-parsed-notes-and-ch
- https://music.stackexchange.com/questions/86241/how-can-i-split-a-midi-file-programatically?rq=1
- https://github.com/WerWolv/ImHex-Patterns

- Check tick to time conversion
  - https://music.stackexchange.com/questions/39446/where-am-i-going-wrong-in-interpreting-this-midi-string?rq=1

### Assembler

- [ ] Assemble: text + templates

### M-IDE
- https://github.com/alecthomas/chroma
- BubbleTea
- templates
- views (e.g. notes)
- macros
- snippets
- apply(...)
- (?) MQL

- [ ] Optimise parsing
      - [ ] context ???
      - [ ] Validate
            - missing/wrong EndOfTrack
      - [ ] Only use ctx when parsing i.e. it shouldn't be a field of MTrk
            - https://dmitrykandalov.com/coroutines-as-threads
      - [ ] Post process tick
      - [ ] NoteOn - unmarshal note and alias and throw an error if they aren't 
            blank and also don't more or less match note value
      - [ ] Move FormatNote to NoteOn (pass ctx as a parameter)
      - [ ] Rework decoder as described in https://go.dev/blog/defer-panic-and-recover
      - [ ] VLQ: TestMarshalBinary
      - [ ] VLQ: UnmarshalBinary
      - [ ] VLF: TestMarshalBinary
      - [ ] VLF: UnmarshalBinary

      - [ ] (maybe) Remove superfluous Event struct
            - (?) or just move bytes and tick back into it so it doesn't clutter event struct
            - (?) bytes are only really used for disassemble
            - (?) tick is used for disassemble and notes
            - (?) so .. either initialise tick and bytes from decoder or move tick and byte into Event
            - https://stackoverflow.com/questions/66118867/go-generics-is-it-possible-to-embed-generic-structs
            - https://stackoverflow.com/questions/71444847/go-with-generics-type-t-is-pointer-to-type-parameter-not-type-parameter


- [ ] MIDI file grammar
      - (?) Treesitter
            - https://tree-sitter.github.io/tree-sitter/
      - (?) ASN.1
      - (?) EBNF
      - (?) Grammar stack
      - (?) https://github.com/codemechanic/midi-sysex-grammar
      - (?) https://www.synalysis.net/grammars/
      - (?) MidiQL

- [ ] Fuzz parser/assembler
      - https://www.pypy.org/posts/2022/12/jit-bug-finding-smt-fuzzing.html

- (?) https://stackoverflow.com/questions/27242652/colorizing-golang-test-run-output
- (?) https://openziti.io/golang-aha-moments-generics
- [ ] TSV
- (?) https://hackaday.com/2023/01/10/imhex-an-open-hex-editor-for-the-modern-hacker/
      - [x] Decode
      - [x] Fixed width output for stdout
      - [x] TSV for file output
      - [x] --delimiter
      - [x] --tabular option
      - [x] Split note and velocity out as separate fields
      - [ ] --piano-roll
      - [ ] --columns
      - [ ] Unit tests
      - [ ] Encode
      - [ ] main_windows.go (Windows doesn't do plugins)


### Transpose
- [ ] Transpose while decoding - otherwise lose track of stuff like note format
      - (!) or - only keep actual MIDI stuff and generate interpretation on the fly for e.g. disassemble
- https://github.com/dolmen-go/goproc

### Click Track

- [ ] --track option (default to track 1)
- [ ] --format consolidated/events
- [ ] Export summary as JSON
- [ ] Export clicks as JSON
- [ ] Export as MIDI
      - [ ] Map beats to clicks
      - [ ] Configurable kick/snare/etc


### Export

### Disassembler

- [ ] Document custom templates
- [ ] Format 2
- [ ] Coloured output (Charm/BubbleTea ?)
- [ ] Add outstanding events to TestDecode
- (?) Rework decoder using tags/reflection/grammar+packrat-parser/kaitai/binpac/somesuch
- [ ] Reference files
      - Format 0
      - [x] Format 1
      - Format 2
- [ ] Decode RPNs
      - https://en.wikipedia.org/wiki/General_MIDI
- [ ] Channel 0-15 or 1-16

### Notes 

- (?) [sync/errgroup](https://pkg.go.dev/golang.org/x/sync/errgroup) for validation 
- [ ] --compact format (JSON): puts notes on one line
- [ ] Configurable formats
- [ ] Check loss of precision
- [ ] Unit tests for tempo map to time conversion
- [ ] Pretty print
- [ ] Format 0
- [ ] Format 2
- [ ] NoteOn with 0 velocity -> NoteOff

## SuperCollider

- [ ] midi2pbind

## TODO

- https://sound.stackexchange.com/questions/51637/velocity-randomizer-any-free-midi-editor-with-this-feature
  - https://stackoverflow.com/questions/46182767/how-to-add-random-number-using-awk-command
- CSS a la Werkmeister
- M-IDE
  - https://arcan-fe.com/2022/10/15/whipping-up-a-new-shell-lashcat9/
  - https://arcan-fe.com/2021/04/12/introducing-pipeworld/
- https://music.stackexchange.com/questions/125617/midi-controlled-metronome

### Other

1.  Tremolo/vibrato
2.  TSV
3.  Export to S-expressions
4.  VSCode plugin
    -  [Language Server Protocol Tutorial: From VSCode to Vim](https://www.toptal.com/javascript/language-server-protocol-tutorial)
5.  Convert between formats 0, 1 and 2
6.  [Manufacturer ID's](https://www.midi.org/specifications-old/category/reference-tables) (?)
7.  Check against reference files from [github:nfroidure/midifile](https://github.com/nfroidure/midifile)
8.  [How to use a field of struct or variable value as template name?](https://stackoverflow.com/questions/28830543/how-to-use-a-field-of-struct-or-variable-value-as-template-name)
9. Online/Javascript version
10. https://github.com/go-interpreter/chezgo
12. SDK (?)
13. mmap
14. REST/GraphQL interface
15. https://sound.stackexchange.com/questions/39457/how-to-open-midi-file-in-text-editor
16. https://github.com/mido/mido
17. [Janet](https://janet-lang.org)
18. [Charm/BubbleTea](https://dlvhdr.me/posts/the-renaissance-of-the-command-line)
19. [GNU Poke](https://youtu.be/Nwb_8VJ5ZeQ)
20. [go-jq](https://github.com/itchyny/gojq)
21. [Katai struct](https://kaitai.io/)
22. [PotLuck](https://www.inkandswitch.com/potluck)
    - [LAPIS](http://groups.csail.mit.edu/graphics/lapis/doc/papers.html)
23. https://futureofcoding.org/catalog/
24. https://github.com/codemechanic/midi-sysex-grammar
25. [WUFFS](https://github.com/google/wuffs)
26. [Apex](https://apexlang.io)
