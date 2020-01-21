package operations

import (
	"github.com/twystd/midiasm/midi"
	"strings"
	"testing"
)

const expected string = `
>>>>>>>>>>>>>>>>>>>>>>>>>
4D 54 68 64 00 00 00 06 00 01 00 02 01 E0   MThd length:6, format:1, tracks:2, metrical time:480 ppqn

4D 54 72 6B 00 00 00 18…                    MTrk 0  length:24
00 FF 03 09 45 78 61 6D 70 6C 65 20 31      tick:0          delta:0          03 TrackName          Example 1
00 FF 51 03 07 A1 20                        tick:0          delta:0          51 Tempo              tempo:500000
00 FF 2F 00                                 tick:0          delta:0          2F EndOfTrack

4D 54 72 6B 00 00 00 B7…                    MTrk 1  length:183
00 FF 00 02 00 17                           tick:0          delta:0          00 SequenceNumber     23
00 FF 01 0D 54 68 69 73 20 61 6E 64 20 54…  tick:0          delta:0          01 Text               This and That
00 FF 02 04 54 68 65 6D                     tick:0          delta:0          02 Copyright          Them
00 FF 03 0F 41 63 6F 75 73 74 69 63 20 47…  tick:0          delta:0          03 TrackName          Acoustic Guitar
00 FF 04 0A 44 69 64 67 65 72 69 64 6F 6F   tick:0          delta:0          04 InstrumentName     Didgeridoo
00 FF 05 08 4C 61 2D 6C 61 2D 6C 61         tick:0          delta:0          05 Lyric              La-la-la
00 FF 06 0F 48 65 72 65 20 42 65 20 44 72…  tick:0          delta:0          06 Marker             Here Be Dragons
00 FF 07 0C 4D 6F 72 65 20 63 6F 77 62 65…  tick:0          delta:0          07 CuePoint           More cowbell
00 FF 08 06 45 73 63 61 70 65               tick:0          delta:0          08 ProgramName        Escape
00 FF 09 08 54 68 65 54 68 69 6E 67         tick:0          delta:0          09 DeviceName         TheThing
00 FF 20 01 0D                              tick:0          delta:0          20 MIDIChannelPrefix  13
00 FF 58 04 04 02 18 08                     tick:0          delta:0          58 TimeSignature      4/4, 24 ticks per click, 8/32 per quarter
00 FF 59 02 00 01                           tick:0          delta:0          59 KeySignature       A minor
00 C0 19                                    tick:0          delta:0          C0 ProgramChange      channel:0  program:25
00 B0 65 00                                 tick:0          delta:0          B0 Controller         channel:0  controller:101, value:0
00 A0 64                                    tick:0          delta:0          A0 PolyphonicPressure channel:0  pressure:100
00 D0 07                                    tick:0          delta:0          D0 ChannelPressure    channel:0  pressure:7
00 90 30 48                                 tick:0          delta:0          90 NoteOn             channel:0  note:C2, velocity:72
81 70 E0 00 08                              tick:240        delta:240        E0 PitchBend          channel:0  bend:8
83 60 80 30 40                              tick:720        delta:480        80 NoteOff            channel:0  note:C2, velocity:64
00 FF 2F 00                                 tick:720        delta:0          2F EndOfTrack


>>>>>>>>>>>>>>>>>>>>>>>>>

`

func TestPrintSMF(t *testing.T) {
	bytes := []byte{
		0x4d, 0x54, 0x68, 0x64, 0x00, 0x00, 0x00, 0x06, 0x00, 0x01, 0x00, 0x02, 0x01, 0xe0,
		0x4d, 0x54, 0x72, 0x6b, 0x00, 0x00, 0x00, 0x18,
		0x00, 0xff, 0x03, 0x09, 0x45, 0x78, 0x61, 0x6d, 0x70, 0x6c, 0x65, 0x20, 0x31,
		0x00, 0xff, 0x51, 0x03, 0x07, 0xa1, 0x20,
		0x00, 0xff, 0x2f, 0x00,
		0x4d, 0x54, 0x72, 0x6b, 0x00, 0x00, 0x00, 0xb7,
		0x00, 0xff, 0x00, 0x02, 0x00, 0x17,
		0x00, 0xff, 0x01, 0x0d, 0x54, 0x68, 0x69, 0x73, 0x20, 0x61, 0x6e, 0x64, 0x20, 0x54, 0x68, 0x61, 0x74,
		0x00, 0xff, 0x02, 0x04, 0x54, 0x68, 0x65, 0x6d,
		0x00, 0xff, 0x03, 0x0f, 0x41, 0x63, 0x6f, 0x75, 0x73, 0x74, 0x69, 0x63, 0x20, 0x47, 0x75, 0x69, 0x74, 0x61, 0x72,
		0x00, 0xff, 0x04, 0x0a, 0x44, 0x69, 0x64, 0x67, 0x65, 0x72, 0x69, 0x64, 0x6f, 0x6f,
		0x00, 0xff, 0x05, 0x08, 0x4c, 0x61, 0x2d, 0x6c, 0x61, 0x2d, 0x6c, 0x61,
		0x00, 0xff, 0x06, 0x0f, 0x48, 0x65, 0x72, 0x65, 0x20, 0x42, 0x65, 0x20, 0x44, 0x72, 0x61, 0x67, 0x6f, 0x6e, 0x73,
		0x00, 0xff, 0x07, 0x0c, 0x4d, 0x6f, 0x72, 0x65, 0x20, 0x63, 0x6f, 0x77, 0x62, 0x65, 0x6c, 0x6c,
		0x00, 0xff, 0x08, 0x06, 0x45, 0x73, 0x63, 0x61, 0x70, 0x65,
		0x00, 0xff, 0x09, 0x08, 0x54, 0x68, 0x65, 0x54, 0x68, 0x69, 0x6e, 0x67,
		0x00, 0xff, 0x20, 0x01, 0x0d,
		0x00, 0xff, 0x58, 0x04, 0x04, 0x02, 0x18, 0x08,
		0x00, 0xff, 0x59, 0x02, 0x00, 0x01,
		0x00, 0xc0, 0x19,
		0x00, 0xb0, 0x65, 0x00,
		0x00, 0xa0, 0x64,
		0x00, 0xd0, 0x07,
		0x00, 0x90, 0x30, 0x48,
		0x81, 0x70, 0xe0, 0x00, 0x08,
		0x83, 0x60, 0x80, 0x30, 0x40,
		0x00, 0xff, 0x2f, 0x00,
	}

	smf := midi.SMF{}
	if err := smf.UnmarshalBinary(bytes); err != nil {
		t.Fatalf("Error unmarshaling SMF: %v", err)
	}

	var s strings.Builder

	printer := Print{}

	printer.PrintWithTemplate(&smf, &s)

	if s.String() != expected {
		p, q := diff(expected, s.String())
		t.Errorf("Output does not match expected:\n%s\n-----\n%s\n%s\n-----\n", s.String(), p, q)
		diff(expected, s.String())
	}
}

func diff(p, q string) (string, string) {
	s := strings.Split(p, "\n")
	t := strings.Split(q, "\n")

	for ix := 0; ix < len(s) && ix < len(t); ix++ {
		if s[ix] != t[ix] {
			u := []rune(s[ix])
			v := []rune(t[ix])
			for jx := 0; jx < len(u) && jx < len(v); jx++ {
				if u[jx] != v[jx] {
					break
				}
				u[jx] = '.'
				v[jx] = '.'
			}

			return string(u), string(v)
		}
	}

	return "?", "?"
}
