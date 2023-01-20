# TODO

- !!!! https://github.com/SuperDisk/tar.pl

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

- https://github.com/WerWolv/ImHex-Patterns

- Check tick to time conversion
  - https://music.stackexchange.com/questions/39446/where-am-i-going-wrong-in-interpreting-this-midi-string?rq=1

### Assembler

- [ ] Assemble: text + templates

### M-IDE
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
      - (?) ASN.1
      - (?) EBNF
      - (?) Grammar stack
      - (?) Treesitter
      - (?) https://github.com/codemechanic/midi-sysex-grammar
      - (?) https://www.synalysis.net/grammars/
      - (?) MidiQL

- [ ] Fuzz parser/assembler
      - https://www.pypy.org/posts/2022/12/jit-bug-finding-smt-fuzzing.html

- (?) https://stackoverflow.com/questions/27242652/colorizing-golang-test-run-output
- (?) https://openziti.io/golang-aha-moments-generics
- (?) https://hackaday.com/2023/01/10/imhex-an-open-hex-editor-for-the-modern-hacker/

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
- [ ] TSV

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
