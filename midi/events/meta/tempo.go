package metaevent

import (
	"encoding/json"
	"fmt"
	"math"
	"regexp"
	"strconv"

	lib "github.com/transcriptaze/midiasm/midi/types"
)

type Tempo struct {
	event
	Tempo uint32
}

func MakeTempo(tick uint64, delta uint32, tempo uint32, bytes ...byte) Tempo {
	return Tempo{
		event: event{
			tick:   tick,
			delta:  delta,
			bytes:  bytes,
			tag:    lib.TagTempo,
			Status: 0xff,
			Type:   lib.TypeTempo,
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

	event := MakeTempo(tick, delta, tempo, append([]byte{0x00, 0xff, 0x51, 0x03}, bytes...)...)

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

func (e *Tempo) UnmarshalText(bytes []byte) error {
	e.tick = 0
	e.delta = 0
	e.bytes = []byte{}
	e.tag = lib.TagTempo
	e.Status = 0xff
	e.Type = 0x51

	re := regexp.MustCompile(`(?i)delta:([0-9]+)(?:.*?)Tempo\s+tempo:([0-9]+)`)
	text := string(bytes)

	if match := re.FindStringSubmatch(text); match == nil || len(match) < 3 {
		return fmt.Errorf("invalid Tempo event (%v)", text)
	} else if delta, err := strconv.ParseUint(match[1], 10, 32); err != nil {
		return err
	} else if v, err := strconv.ParseUint(match[2], 10, 32); err != nil {
		return err
	} else {
		e.delta = uint32(delta)
		e.Tempo = uint32(v)
	}

	return nil
}

func (e Tempo) MarshalJSON() (encoded []byte, err error) {
	t := struct {
		Tag    string `json:"tag"`
		Delta  uint32 `json:"delta"`
		Status byte   `json:"status"`
		Type   byte   `json:"type"`
		Tempo  uint32 `json:"tempo"`
	}{
		Tag:    fmt.Sprintf("%v", e.tag),
		Delta:  e.delta,
		Status: byte(e.Status),
		Type:   byte(e.Type),
		Tempo:  e.Tempo,
	}

	return json.Marshal(t)
}

func (e *Tempo) UnmarshalJSON(bytes []byte) error {
	e.tick = 0
	e.delta = 0
	e.bytes = []byte{}
	e.Status = 0xff
	e.tag = lib.TagTempo
	e.Type = lib.TypeTempo

	t := struct {
		Tag   string `json:"tag"`
		Delta uint32 `json:"delta"`
		Tempo uint32 `json:"tempo"`
	}{}

	if err := json.Unmarshal(bytes, &t); err != nil {
		return err
	} else if t.Tag != "Tempo" {
		return fmt.Errorf("invalid %v event (%v)", e.tag, string(bytes))
	} else {
		e.delta = t.Delta
		e.Tempo = t.Tempo
	}

	return nil
}
