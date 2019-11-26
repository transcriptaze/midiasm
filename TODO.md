## v0.0

*Disassembler*

- [x] Rework MIDI event parser
- [x] Rework META event parser
- [x] Extract notes
- [x] Use microseconds as integer time base
- [x] Ellipsize too long hex
- [ ] Write to file, log errors/warning to stderr
- [ ] Split tracks to separate files
- [ ] --debug
- [ ] --verbose
- [ ] Configurable formats
- [ ] EventList: pretty print
- [ ] NoteList: pretty print
- [ ] Check loss of precision

### MIDI events

- [x] 8n/Note Off
- [x] 9n/Note On
- [x] An/Polyphonic Pressure
- [x] Bn/Controller
- [x] Cn/Program Change
- [x] Dn/Channel Pressure
- [x] En/Pitch Bend

### META events

- [ ] 00/Sequence Number
- [ ] 01/Text
- [ ] 02/Copyright
- [x] 03/Track Name
- [ ] 04/Instrument Name
- [ ] 05/Lyric
- [ ] 06/Marker
- [ ] 07/Cue Point
- [x] 08/Program Name
- [ ] 09/Device Name
- [ ] 20/MIDI Channel Prefix
- [ ] 21/MIDI Port
- [x] 2F/End of Track
- [x] 51/Tempo
- [ ] 54/SMPTE Offset
- [x] 58/Time Signature
- [x] 59/Key Signature
- [ ] 7F/Sequencer Specific Event

- [x] TimeSignature: [Unicode fractions](http://unicodefractions.com)
- [x] KeySignature:  [Unicode symbols](https://unicode-table.com/en/blocks/musical-symbols/)

### SysEx events

- [ ] Single messages
- [ ] Continuation events
- [ ] Escape sequences


## TODO

1.  Assembler
2.  Export to JSON
3.  Export to S-expressions
4.  VSCode plugin

