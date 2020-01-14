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
00 FF 03 09 45 78 61 6D 70 6C 65 20 31      tick:0          delta:0          03 TrackName        Example 1
00 FF 51 03 07 A1 20                        tick:0          delta:0          51 Tempo            tempo:500000
00 FF 2F 00                                 tick:0          delta:0          2F EndOfTrack

4D 54 72 6B 00 00 00 22…                    MTrk 1  length:34
00 FF 03 0F 41 63 6F 75 73 74 69 63 20 47…  tick:0          delta:0          03 TrackName        Acoustic Guitar
00 C0 19                                    tick:0          delta:0          C0 ProgramChange    channel:0 program:25
00 FF 58 04 04 02 18 08                     tick:0          delta:0          58 TimeSignature    4/4, 24 ticks per click, 8/32 per quarter
00 FF 2F 00                                 tick:0          delta:0          2F EndOfTrack


>>>>>>>>>>>>>>>>>>>>>>>>>

`

func TestPrintSMF(t *testing.T) {
	bytes := []byte{
		0x4d, 0x54, 0x68, 0x64, 0x00, 0x00, 0x00, 0x06, 0x00, 0x01, 0x00, 0x02, 0x01, 0xe0,
		0x4d, 0x54, 0x72, 0x6b, 0x00, 0x00, 0x00, 0x18,
		0x00, 0xFF, 0x03, 0x09, 0x45, 0x78, 0x61, 0x6D, 0x70, 0x6C, 0x65, 0x20, 0x31,
		0x00, 0xFF, 0x51, 0x03, 0x07, 0xA1, 0x20,
		0x00, 0xFF, 0x2F, 0x00,
		0x4d, 0x54, 0x72, 0x6b, 0x00, 0x00, 0x00, 0x22,
		0x00, 0xff, 0x03, 0x0f, 0x41, 0x63, 0x6f, 0x75, 0x73, 0x74, 0x69, 0x63, 0x20, 0x47, 0x75, 0x69, 0x74, 0x61, 0x72,
		0x00, 0xC0, 0x19,
		0x00, 0xFF, 0x58, 0x04, 0x04, 0x02, 0x18, 0x08,
		0x00, 0xFF, 0x2F, 0x00,
	}

	smf := midi.SMF{}
	if err := smf.UnmarshalBinary(bytes); err != nil {
		t.Fatalf("Error unmarshaling SMF: %v", err)
	}

	var s strings.Builder

	printer := Print{}

	printer.printWithTemplate(&smf, &s)

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
