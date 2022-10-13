package metaevent

import (
	"fmt"
	"math"
	"regexp"
	"strconv"

	"github.com/transcriptaze/midiasm/midi/types"
)

type Tempo struct {
	Tag    string
	Status types.Status
	Type   types.MetaEventType
	Tempo  uint32
}

func NewTempo(bytes []byte) (*Tempo, error) {
	if len(bytes) != 3 {
		return nil, fmt.Errorf("Invalid Tempo length (%d): expected '3'", len(bytes))
	}

	tempo := uint32(0)
	for _, b := range bytes {
		tempo <<= 8
		tempo += uint32(b)
	}

	return &Tempo{
		Tag:    "Tempo",
		Status: 0xff,
		Type:   0x51,
		Tempo:  tempo,
	}, nil
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
	//	e.tick = 0
	//	e.delta = 0
	//	e.bytes = []byte{}
	t.Status = 0xff
	t.Tag = "Tempo"
	t.Type = 0x51

	re := regexp.MustCompile(`(?i)Tempo\s+tempo:([0-9]+)`)
	text := string(bytes)

	if match := re.FindStringSubmatch(text); match == nil || len(match) < 2 {
		return fmt.Errorf("invalid Tempo event (%v)", text)
	} else if v, err := strconv.ParseUint(match[1], 10, 32); err != nil {
		return err
	} else {
		t.Tempo = uint32(v)
	}

	return nil
}
