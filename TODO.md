## v0.0

*Disassembler*

- [x] Rework MIDI event parser
- [x] Rework META event parser
- [x] Extract notes
- [x] Use microseconds as integer time base
- [x] Ellipsize too long hex
- [x] Log errors/warning to stderr
- [x] Write to file
- [ ] Split tracks to separate files
- [ ] Validate (missing end of track, tempo events, etc)
- [x] --debug
- [x] --verbose
- [x] Print note name + octave
- [ ] Configurable formats
- [ ] Pretty print
- [ ] Format 0
- [ ] Format 2
- [ ] MThd hex
- [ ] MTrk hex

### Notes 

- [x] Print note name + octave
- [ ] Rework as SMF processor
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
- [xq] 7F/Sequencer Specific Event

- [x] TimeSignature: [Unicode fractions](http://unicodefractions.com)
- [x] KeySignature:  [Unicode symbols](https://unicode-table.com/en/blocks/musical-symbols/)

### SysEx events

- [ ] Single messages
- [ ] Continuation events
- [ ] Escape sequences


## TODO

1.  Assembler
2.  TSV
3.  Export to JSON
4.  Export to S-expressions
5.  VSCode plugin
6.  Convert between formats 0, 1 and 2

