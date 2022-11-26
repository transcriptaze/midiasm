package lib

import (
	"fmt"
	"testing"
)

func TestTagStringer(t *testing.T) {
	tests := []struct {
		tag      Tag
		expected string
	}{
		{TagUnknown, ""},
		{TagSequenceNumber, "SequenceNumber"},
		{TagCopyright, "Copyright"},
		{TagCuePoint, "CuePoint"},
		{TagDeviceName, "DeviceName"},
		{TagInstrumentName, "InstrumentName"},
		{TagKeySignature, "KeySignature"},
		{TagLyric, "Lyric"},
		{TagMarker, "Marker"},
		{TagMIDIChannelPrefix, "MIDIChannelPrefix"},
		{TagMIDIPort, "MIDIPort"},
		{TagProgramName, "ProgramName"},
		{TagSMPTEOffset, "SMPTEOffset"},
		{TagTempo, "Tempo"},
		{TagText, "Text"},
		{TagTimeSignature, "TimeSignature"},
		{TagTrackName, "TrackName"},
		{TagEndOfTrack, "EndOfTrack"},
		{TagSequencerSpecificEvent, "SequencerSpecificEvent"},
		{TagChannelPressure, "ChannelPressure"},
		{TagController, "Controller"},
		{TagNoteOff, "NoteOff"},
		{TagNoteOn, "NoteOn"},
		{TagPitchBend, "PitchBend"},
		{TagPolyphonicPressure, "PolyphonicPressure"},
		{TagProgramChange, "ProgramChange"},
		{TagSysExContinuation, "SysExContinuation"},
		{TagSysExEscape, "SysExEscape"},
		{TagSysExMessage, "SysExMessage"},
	}

	for _, test := range tests {
		if s := fmt.Sprintf("%v", test.tag); s != test.expected {
			t.Errorf("Incorrect label for tag %v - expected:%v, got:%v", test.tag, test.expected, s)
		}
	}

}

func TestTagUnmarshalText(t *testing.T) {
	tests := []struct {
		text     string
		expected Tag
	}{
		{"Copyright", TagCopyright},
		{"CuePoint", TagCuePoint},
		{"DeviceName", TagDeviceName},
		{"EndOfTrack", TagEndOfTrack},
		{"InstrumentName", TagInstrumentName},
		{"KeySignature", TagKeySignature},
		{"Lyric", TagLyric},
		{"Marker", TagMarker},
		{"MIDIChannelPrefix", TagMIDIChannelPrefix},
		{"MIDIPort", TagMIDIPort},
		{"ProgramName", TagProgramName},
		{"SequenceNumber", TagSequenceNumber},
		{"SMPTEOffset", TagSMPTEOffset},
		{"Tempo", TagTempo},
		{"Text", TagText},
		{"TimeSignature", TagTimeSignature},
		{"TrackName", TagTrackName},
		{"ChannelPressure", TagChannelPressure},
		{"Controller", TagController},
		{"NoteOff", TagNoteOff},
		{"NoteOn", TagNoteOn},
		{"PitchBend", TagPitchBend},
		{"PolyphonicPressure", TagPolyphonicPressure},
		{"ProgramChange", TagProgramChange},
		{"SysExContinuation", TagSysExContinuation},
		{"SysExEscape", TagSysExEscape},
		{"SysExMessage", TagSysExMessage},
	}

	for _, test := range tests {
		var tag Tag
		if err := tag.UnmarshalText([]byte(test.text)); err != nil {
			t.Fatalf("Error unmarshalling %v (%v)", test.text, err)
		} else if tag != test.expected {
			t.Errorf("Incorrect tag for %q - expected:%v, got:%v", test.text, test.expected, tag)
		}
	}
}
