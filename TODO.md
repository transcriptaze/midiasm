## v0.0

*Disassembler*

- [x] Rework MIDI event parser
- [x] Rework META event parser
- [x] Extract notes
- [x] Use microseconds as integer time base
- [x] Ellipsize too long hex
- [x] Log errors/warning to stderr
- [x] Write to file
- [x] Split tracks to separate files
- [x] Validate (missing end of track, tempo events)
- [x] --debug
- [x] --verbose
- [x] Print note name + octave
- [x] Decode note in context of current scale
- [ ] Configurable formats
  - Remove old 'Render' implementation
  - Load format from file
  - Identify manufacturer for SysEx and SequencerSpecificEvent (http://www.somascape.org/midi/tech/spec.html#sysexnotes)
  - Align event bytes on delta time
  - Keep NoteOff name to be NoteOn name if KeySignature changes during duration of note (?)
- [ ] Format 0
- [ ] Format 2
- [ ] Reference files
- [ ] Check SMTPEOffset only in track 0 for format 1
- [ ] Running status (?)

### Notes 

- [x] Print note name + octave
- [x] Rework as SMF processor
- [ ] Check loss of precision
- [ ] Unit tests for tempo map to time conversion
- [ ] Pretty print
- [ ] Format 0
- [ ] Format 2

### MIDI events

- [x] 8n/Note Off
- [x] 9n/Note On
- [x] An/Polyphonic Pressure
- [x] Bn/Controller
- [x] Cn/Program Change
- [x] Dn/Channel Pressure
- [x] En/Pitch Bend

### META events

- [x] 00/Sequence Number
- [x] 01/Text
- [x] 02/Copyright
- [x] 03/Track Name
- [x] 04/Instrument Name
- [x] 05/Lyric
- [x] 06/Marker
- [x] 07/Cue Point
- [x] 08/Program Name
- [x] 09/Device Name
- [x] 20/MIDI Channel Prefix
- [x] 21/MIDI Port
- [x] 2F/End of Track
- [x] 51/Tempo
- [x] 54/SMPTE Offset
- [x] 58/Time Signature
- [x] 59/Key Signature
- [x] 7F/Sequencer Specific Event

- [x] TimeSignature: [Unicode fractions](http://unicodefractions.com)
- [x] KeySignature:  [Unicode symbols](https://unicode-table.com/en/blocks/musical-symbols/)

### SysEx events

- [x] Single messages
- [x] Continuation events
- [x] Escape sequences
- [x] Casio sequences

## TODO

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
11. Rework decoder using tags/reflection
12. https://github.com/go-interpreter/chezgo
13. SDK (?)

