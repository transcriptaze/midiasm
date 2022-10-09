# TODO


- [ ] dissassemble // FIXME MThd suddenly got one space shorter

- [ ] Optimise parsing
- (?) SMF: TextMarshal/Unmarshal
- (?) SMF: JsonMarshal/Unmarshal

### Assembler
- [x] Command skeleton
- [x] Assemble MThd
- [ ] Assemble MTrk
- [ ] Parse MThd
- [ ] Parse MTrk
- [ ] VLQ: TestMarshalBinary
- [ ] VLQ: UnmarshalBinary
- [ ] VLF: TestMarshalBinary
- [ ] VLF: UnmarshalBinary
- [ ] Event needs become embedded in Meta/Midi

### Transpose
      - [ ] Transpose while decoding - otherwise lose track of stuff 

### Click Track

- [ ] --track option (default to track 1)
- [ ] --format consolidated/events
- [ ] Export summary as JSON
- [ ] Export clicks as JSON
- [ ] Export as MIDI
      - [ ] Map beats to clicks
      - [ ] Configurable kick/snare/etc


### Disassembler

- [ ] Document custom templates
- [ ] Move TimeSignature in reference-01.mid to track 0
- [ ] Format 2
- [ ] Coloured output (Charm/BubbleTea ?)
- [ ] Add outstanding events to TestDecode
- (?) Rework decoder using tags/reflection/grammar+packrat-parser/kaitai/binpac/somesuch
- [ ] Reference files
      - Format 0
      - Format 1
      - Format 2
- [ ] Decode RPNs
      - https://en.wikipedia.org/wiki/General_MIDI
- [ ] Channel 0-15 or 1-16

### Notes 

- [ ] --compact format (JSON): puts notes on one line
- [ ] Configurable formats
- [ ] Check loss of precision
- [ ] Unit tests for tempo map to time conversion
- [ ] Pretty print
- [ ] Format 0
- [ ] Format 2
- [ ] NoteOn with 0 velocity -> NoteOff

## TODO

- https://sound.stackexchange.com/questions/51637/velocity-randomizer-any-free-midi-editor-with-this-feature

### Other

1.  Tremolo/vibrato
2.  TSV
3.  Export to JSON
4.  Export to S-expressions
5.  VSCode plugin
    -  [Language Server Protocol Tutorial: From VSCode to Vim](https://www.toptal.com/javascript/language-server-protocol-tutorial)
6.  Convert between formats 0, 1 and 2
7.  [Manufacturer ID's](https://www.midi.org/specifications-old/category/reference-tables) (?)
8.  Check against reference files from [github:nfroidure/midifile](https://github.com/nfroidure/midifile)
9.  [How to use a field of struct or variable value as template name?](https://stackoverflow.com/questions/28830543/how-to-use-a-field-of-struct-or-variable-value-as-template-name)
10. Online/Javascript version
12. https://github.com/go-interpreter/chezgo
13. SDK (?)
14. mmap
15. REST/GraphQL interface
16. https://sound.stackexchange.com/questions/39457/how-to-open-midi-file-in-text-editor
17. https://github.com/mido/mido
18. [Janet](https://janet-lang.org)
19. [Charm/BubbleTea](https://dlvhdr.me/posts/the-renaissance-of-the-command-line)
20. [GNU Poke](https://youtu.be/Nwb_8VJ5ZeQ)
21. [go-jq](https://github.com/itchyny/gojq)
