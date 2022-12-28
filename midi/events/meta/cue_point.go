package metaevent

import (
	"encoding/json"
	"fmt"
	"regexp"

	"github.com/transcriptaze/midiasm/midi/lib"
)

type CuePoint struct {
	event
	CuePoint string
}

func MakeCuePoint(tick uint64, delta lib.Delta, cuepoint string, bytes ...byte) CuePoint {
	return CuePoint{
		event: event{
			tick:   tick,
			delta:  delta,
			bytes:  bytes,
			tag:    lib.TagCuePoint,
			Status: 0xff,
			Type:   lib.TypeCuePoint,
		},
		CuePoint: cuepoint,
	}
}

func (e *CuePoint) unmarshal(tick uint64, delta lib.Delta, status byte, data []byte, bytes ...byte) error {
	cuepoint := string(data)

	*e = MakeCuePoint(tick, delta, cuepoint, bytes...)

	return nil
}

func (e CuePoint) MarshalBinary() (encoded []byte, err error) {
	return append([]byte{
		byte(e.Status),
		byte(e.Type),
		byte(len(e.CuePoint)),
	},
		[]byte(e.CuePoint)...), nil
}

func (e *CuePoint) UnmarshalBinary(bytes []byte) error {
	if delta, remaining, err := delta(bytes); err != nil {
		return err
	} else if len(remaining) < 2 {
		return fmt.Errorf("Invalid event (%v)", remaining)
	} else if remaining[0] != 0xff {
		return fmt.Errorf("Invalid %v status (%02X)", lib.TagCuePoint, remaining[0])
	} else if !equals(remaining[1], lib.TypeCuePoint) {
		return fmt.Errorf("Invalid %v event type (%02X)", lib.TagCuePoint, remaining[1])
	} else if cuepoint, err := vlf(remaining[2:]); err != nil {
		return err
	} else {
		*e = MakeCuePoint(0, delta, string(cuepoint), bytes...)
	}

	return nil
}

func (e *CuePoint) UnmarshalText(bytes []byte) error {
	re := regexp.MustCompile(`(?i)delta:([0-9]+)(?:.*?)CuePoint\s+(.*)`)
	text := string(bytes)

	if match := re.FindStringSubmatch(text); match == nil || len(match) < 3 {
		return fmt.Errorf("invalid CuePoint event (%v)", text)
	} else if delta, err := lib.ParseDelta(match[1]); err != nil {
		return err
	} else {
		*e = MakeCuePoint(0, delta, match[2], []byte{}...)
	}

	return nil
}

func (e CuePoint) MarshalJSON() (encoded []byte, err error) {
	t := struct {
		Tag      string    `json:"tag"`
		Delta    lib.Delta `json:"delta"`
		Status   byte      `json:"status"`
		Type     byte      `json:"type"`
		CuePoint string    `json:"cuepoint"`
	}{
		Tag:      fmt.Sprintf("%v", e.tag),
		Delta:    e.delta,
		Status:   byte(e.Status),
		Type:     byte(e.Type),
		CuePoint: e.CuePoint,
	}

	return json.Marshal(t)
}

func (e *CuePoint) UnmarshalJSON(bytes []byte) error {
	t := struct {
		Tag      string    `json:"tag"`
		Delta    lib.Delta `json:"delta"`
		CuePoint string    `json:"cuepoint"`
	}{}

	if err := json.Unmarshal(bytes, &t); err != nil {
		return err
	} else if !equal(t.Tag, lib.TagCuePoint) {
		return fmt.Errorf("invalid %v event (%v)", e.tag, string(bytes))
	} else {
		*e = MakeCuePoint(0, t.Delta, t.CuePoint, []byte{}...)
	}

	return nil
}
