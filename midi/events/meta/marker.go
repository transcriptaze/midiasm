package metaevent

import (
	"fmt"
	"regexp"

	lib "github.com/transcriptaze/midiasm/midi/types"
)

type Marker struct {
	event
	Marker string
}

func MakeMarker(tick uint64, delta lib.Delta, marker string) Marker {
	return Marker{
		event: event{
			tick:   tick,
			delta:  delta,
			bytes:  append([]byte{0x00, 0xff, 0x06, byte(len(marker))}, []byte(marker)...),
			tag:    lib.TagMarker,
			Status: 0xff,
			Type:   lib.TypeMarker,
		},
		Marker: marker,
	}
}

func UnmarshalMarker(tick uint64, delta lib.Delta, bytes []byte) (*Marker, error) {
	marker := string(bytes)
	event := MakeMarker(tick, delta, marker)

	return &event, nil
}

func (m Marker) MarshalBinary() (encoded []byte, err error) {
	return append([]byte{
		byte(m.Status),
		byte(m.Type),
		byte(len(m.Marker)),
	},
		[]byte(m.Marker)...), nil
}

func (m *Marker) UnmarshalText(bytes []byte) error {
	m.tick = 0
	m.delta = 0
	m.bytes = []byte{}
	m.tag = lib.TagMarker
	m.Status = 0xff
	m.Type = lib.TypeMarker

	re := regexp.MustCompile(`(?i)delta:([0-9]+)(?:.*?)Marker\s+(.*)`)
	text := string(bytes)

	if match := re.FindStringSubmatch(text); match == nil || len(match) < 3 {
		return fmt.Errorf("invalid Marker event (%v)", text)
	} else if delta, err := lib.ParseDelta(match[1]); err != nil {
		return err
	} else {
		m.delta = delta
		m.Marker = string(match[2])
	}

	return nil
}
