package types

import (
	"fmt"
	"strings"
)

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

var tags = map[Tag]string{
	TagCopyright:          "Copyright",
	TagCuePoint:           "CuePoint",
	TagDeviceName:         "DeviceName",
	TagEndOfTrack:         "EndOfTrack",
	TagInstrumentName:     "InstrumentName",
	TagKeySignature:       "KeySignature",
	TagLyric:              "Lyric",
	TagMarker:             "Marker",
	TagMIDIChannelPrefix:  "MIDIChannelPrefix",
	TagMIDIPort:           "MIDIPort",
	TagProgramName:        "ProgramName",
	TagSequenceNumber:     "SequenceNumber",
	TagSMPTEOffset:        "SMPTEOffset",
	TagTempo:              "Tempo",
	TagText:               "Text",
	TagTimeSignature:      "TimeSignature",
	TagTrackName:          "TrackName",
	TagChannelPressure:    "ChannelPressure",
	TagController:         "Controller",
	TagNoteOff:            "NoteOff",
	TagNoteOn:             "NoteOn",
	TagPitchBend:          "PitchBend",
	TagPolyphonicPressure: "PolyphonicPressure",
	TagProgramChange:      "ProgramChange",
	TagSysExContinuation:  "SysExContinuation",
	TagSysExEscape:        "SysExEscape",
	TagSysExMessage:       "SysExMessage",
}

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

func (t *Tag) UnmarshalText(bytes []byte) error {
	text := string(bytes)

	for k, v := range tags {
		if strings.Contains(text, v) {
			*t = k
			return nil
		}
	}

	return fmt.Errorf("No matching tag")

}
