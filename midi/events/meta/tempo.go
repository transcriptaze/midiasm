package metaevent

import (
	"fmt"
	"math"
	"regexp"
	"strconv"

	"github.com/transcriptaze/midiasm/midi/types"
)

type Tempo struct {
	event
	Tempo uint32
}

func MakeTempo(tick uint64, delta uint32, tempo uint32) Tempo {
	return Tempo{
		event: event{
			tick:   tick,
			delta:  delta,
			bytes:  append([]byte{0x00, 0xff, 0x51, 03}, byte(tempo>>16&0x0ff), byte(tempo>>8&0x0ff), byte(tempo>>0&0x0ff)),
			tag:    types.TagTempo,
			Status: 0xff,
			Type:   types.TypeTempo,
		},
		Tempo: tempo,
	}
}

func UnmarshalTempo(tick uint64, delta uint32, bytes []byte) (*Tempo, error) {
	if len(bytes) != 3 {
		return nil, fmt.Errorf("Invalid Tempo length (%d): expected '3'", len(bytes))
	}

	tempo := uint32(0)
	for _, b := range bytes {
		tempo <<= 8
		tempo += uint32(b)
	}

	event := MakeTempo(tick, delta, tempo)

	return &event, nil
}

func (t Tempo) String() string {
	bpm := uint(math.Round(60.0 * 1000000.0 / float64(t.Tempo)))

	return fmt.Sprintf("%v", bpm)
}

func (t Tempo) MarshalBinary() (encoded []byte, err error) {
	encoded = make([]byte, 6)

	encoded[0] = byte(t.Status)
	encoded[1] = byte(t.Type)
	encoded[2] = byte(3)
	encoded[3] = byte(t.Tempo >> 16 & 0xff)
	encoded[4] = byte(t.Tempo >> 8 & 0xff)
	encoded[5] = byte(t.Tempo >> 0 & 0xff)

	return
}

func (t *Tempo) UnmarshalText(bytes []byte) error {
	t.tick = 0
	t.delta = 0
	t.bytes = []byte{}
	t.tag = types.TagTempo
	t.Status = 0xff
	t.Type = 0x51

	re := regexp.MustCompile(`(?i)delta:([0-9]+)(?:.*?)Tempo\s+tempo:([0-9]+)`)
	text := string(bytes)

	if match := re.FindStringSubmatch(text); match == nil || len(match) < 3 {
		return fmt.Errorf("invalid Tempo event (%v)", text)
	} else if delta, err := strconv.ParseUint(match[1], 10, 32); err != nil {
		return err
	} else if v, err := strconv.ParseUint(match[2], 10, 32); err != nil {
		return err
	} else {
		t.delta = uint32(delta)
		t.Tempo = uint32(v)
	}

	return nil
}
