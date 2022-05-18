## v0.0

- Take a look at [GNU Poke](https://youtu.be/Nwb_8VJ5ZeQ)
- [ ] Redo logging

### Disassembler

- [ ] Format 2
- [ ] Add outstanding events to TestDecode
- (?) Rework decoder using tags/reflection/grammar+packrat-parser/kaitai/binpac/somesuch
- [ ] Reference files
      - Format 0
      - Format 1
      - Format 2

### Click Track

- [x] Extract time signature changes
- [x] Extract bars
- [ ] Export as JSON
- [ ] Export as MIDI
      - [ ] Map beats to clicks
      - [ ] Configurable kick/snare/etc
- [ ] Include tempo changes

### Notes 

- [ ] Check loss of precision
- [ ] Unit tests for tempo map to time conversion
- [ ] Configurable formats
- [ ] Pretty print
- [ ] Format 0
- [ ] Format 2
- [ ] NoteOn with 0 velocity -> NoteOff

### MIDI events

### META events

- [x] TimeSignature: [Unicode fractions](http://unicodefractions.com)
- [x] KeySignature:  [Unicode symbols](https://unicode-table.com/en/blocks/musical-symbols/)

### SysEx events

## TODO

### Other

1.  Assembler
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

