package lib

import (
	"fmt"
)

type TEventType interface {
	MetaEventType | MidiEventType | SysExEventType

	Equals(byte) bool
}

type MetaEventType byte

const (
	TypeSequenceNumber         MetaEventType = 0x00
	TypeText                   MetaEventType = 0x01
	TypeCopyright              MetaEventType = 0x02
	TypeTrackName              MetaEventType = 0x03
	TypeInstrumentName         MetaEventType = 0x04
	TypeLyric                  MetaEventType = 0x05
	TypeMarker                 MetaEventType = 0x06
	TypeCuePoint               MetaEventType = 0x07
	TypeProgramName            MetaEventType = 0x08
	TypeDeviceName             MetaEventType = 0x09
	TypeMIDIChannelPrefix      MetaEventType = 0x20
	TypeMIDIPort               MetaEventType = 0x21
	TypeTempo                  MetaEventType = 0x51
	TypeSMPTEOffset            MetaEventType = 0x54
	TypeTimeSignature          MetaEventType = 0x58
	TypeKeySignature           MetaEventType = 0x59
	TypeEndOfTrack             MetaEventType = 0x2f
	TypeSequencerSpecificEvent MetaEventType = 0x7f
)

func (t MetaEventType) String() string {
	return fmt.Sprintf("%02X", byte(t))
}

func (t MetaEventType) Equals(b byte) bool {
	return (b & 0x7f) == byte(t)
}

type MidiEventType byte

const (
	TypeNoteOff            MidiEventType = 0x80
	TypeNoteOn             MidiEventType = 0x90
	TypePolyphonicPressure MidiEventType = 0xa0
	TypeController         MidiEventType = 0xb0
	TypeProgramChange      MidiEventType = 0xc0
	TypeChannelPressure    MidiEventType = 0xd0
	TypePitchBend          MidiEventType = 0xe0
)

func (t MidiEventType) String() string {
	return fmt.Sprintf("%02X", byte(t))
}

func (t MidiEventType) Equals(b byte) bool {
	return (b & 0xf0) == byte(t)
}

type SysExEventType byte

const (
	TypeSysExMessage SysExEventType = 0xf0
)

func (t SysExEventType) String() string {
	return fmt.Sprintf("%02X", byte(t))
}

func (t SysExEventType) Equals(b byte) bool {
	return (b & 0xf0) == byte(t)
}
