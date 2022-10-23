package types

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
		{TagCopyright, "Copyright"},
		{TagCuePoint, "CuePoint"},
		{TagDeviceName, "DeviceName"},
		{TagEndOfTrack, "EndOfTrack"},
		{TagInstrumentName, "InstrumentName"},
		{TagKeySignature, "KeySignature"},
		{TagLyric, "Lyric"},
		{TagMarker, "Marker"},
		{TagMIDIChannelPrefix, "MIDIChannelPrefix"},
		{TagMIDIPort, "MIDIPort"},
		{TagProgramName, "ProgramName"},
		{TagSequenceNumber, "SequenceNumber"},
		{TagSMPTEOffset, "SMPTEOffset"},
		{TagTempo, "Tempo"},
		{TagText, "Text"},
		{TagTimeSignature, "TimeSignature"},
		{TagTrackName, "TrackName"},
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
