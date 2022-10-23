package types

type Tag int

const (
	TagUnknown Tag = iota
	TagCopyright
	TagCuePoint
	TagDeviceName
	TagEndOfTrack
	TagInstrumentName
	TagKeySignature
	TagLyric
	TagMarker
	TagMIDIChannelPrefix
	TagMIDIPort
	TagProgramName
	TagSequenceNumber
	TagSMPTEOffset
	TagTempo
	TagText
	TagTimeSignature
	TagTrackName
	TagChannelPressure
	TagController
	TagNoteOff
	TagNoteOn
	TagPitchBend
	TagPolyphonicPressure
	TagProgramChange
	TagSysExContinuation
	TagSysExEscape
	TagSysExMessage
)

func (t Tag) String() string {
	return []string{
		"",
		"Copyright",
		"CuePoint",
		"DeviceName",
		"EndOfTrack",
		"InstrumentName",
		"KeySignature",
		"Lyric",
		"Marker",
		"MIDIChannelPrefix",
		"MIDIPort",
		"ProgramName",
		"SequenceNumber",
		"SMPTEOffset",
		"Tempo",
		"Text",
		"TimeSignature",
		"TrackName",
		"ChannelPressure",
		"Controller",
		"NoteOff",
		"NoteOn",
		"PitchBend",
		"PolyphonicPressure",
		"ProgramChange",
		"SysExContinuation",
		"SysExEscape",
		"SysExMessage",
	}[t]
}
