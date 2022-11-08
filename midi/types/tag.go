package types

import (
	"fmt"
	"strings"
)

type Tag int

const (
	TagUnknown Tag = iota
	TagSequenceNumber
	TagText
	TagCopyright
	TagTrackName
	TagInstrumentName
	TagLyric
	TagMarker
	TagCuePoint
	TagProgramName
	TagDeviceName
	TagMIDIChannelPrefix
	TagMIDIPort
	TagTempo
	TagSMPTEOffset
	TagTimeSignature
	TagKeySignature
	TagEndOfTrack
	TagSequencerSpecificEvent

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
	TagSequenceNumber:         "SequenceNumber",
	TagText:                   "Text",
	TagCopyright:              "Copyright",
	TagTrackName:              "TrackName",
	TagInstrumentName:         "InstrumentName",
	TagLyric:                  "Lyric",
	TagMarker:                 "Marker",
	TagCuePoint:               "CuePoint",
	TagProgramName:            "ProgramName",
	TagDeviceName:             "DeviceName",
	TagMIDIChannelPrefix:      "MIDIChannelPrefix",
	TagMIDIPort:               "MIDIPort",
	TagTempo:                  "Tempo",
	TagSMPTEOffset:            "SMPTEOffset",
	TagTimeSignature:          "TimeSignature",
	TagKeySignature:           "KeySignature",
	TagEndOfTrack:             "EndOfTrack",
	TagSequencerSpecificEvent: "SequencerSpecificEvent",

	TagChannelPressure:    "ChannelPressure",
	TagController:         "Controller",
	TagNoteOff:            "NoteOff",
	TagNoteOn:             "NoteOn",
	TagPitchBend:          "PitchBend",
	TagPolyphonicPressure: "PolyphonicPressure",
	TagProgramChange:      "ProgramChange",

	TagSysExMessage:      "SysExMessage",
	TagSysExContinuation: "SysExContinuation",
	TagSysExEscape:       "SysExEscape",
}

func (t Tag) String() string {
	return []string{
		"",
		"SequenceNumber",
		"Text",
		"Copyright",
		"TrackName",
		"InstrumentName",
		"Lyric",
		"Marker",
		"CuePoint",
		"ProgramName",
		"DeviceName",
		"MIDIChannelPrefix",
		"MIDIPort",
		"Tempo",
		"SMPTEOffset",
		"TimeSignature",
		"KeySignature",
		"EndOfTrack",
		"SequencerSpecificEvent",
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
