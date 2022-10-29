package metaevent

import (
	"fmt"
	"regexp"
	"strconv"

	"github.com/transcriptaze/midiasm/midi/types"
)

type InstrumentName struct {
	event
	Name string
}

func MakeInstrumentName(tick uint64, delta uint32, name string) InstrumentName {
	return InstrumentName{
		event: event{
			tick:   tick,
			delta:  delta,
			bytes:  append([]byte{0x00, 0xff, 0x04, byte(len(name))}, []byte(name)...),
			tag:    types.TagInstrumentName,
			Status: 0xff,
			Type:   types.TypeInstrumentName,
		},
		Name: name,
	}
}

func UnmarshalInstrumentName(tick uint64, delta uint32, bytes []byte) (*InstrumentName, error) {
	name := string(bytes)
	event := MakeInstrumentName(tick, delta, name)

	return &event, nil
}

func (n InstrumentName) MarshalBinary() (encoded []byte, err error) {
	return append([]byte{
		byte(n.Status),
		byte(n.Type),
		byte(len(n.Name)),
	},
		[]byte(n.Name)...), nil
}

func (n *InstrumentName) UnmarshalText(bytes []byte) error {
	n.tick = 0
	n.delta = 0
	n.bytes = []byte{}
	n.tag = types.TagInstrumentName
	n.Status = 0xff
	n.Type = types.TypeInstrumentName

	re := regexp.MustCompile(`(?i)delta:([0-9]+)(?:.*?)InstrumentName\s+(.*)`)
	text := string(bytes)

	if match := re.FindStringSubmatch(text); match == nil || len(match) < 3 {
		return fmt.Errorf("invalid InstrumentName event (%v)", text)
	} else if delta, err := strconv.ParseUint(match[1], 10, 32); err != nil {
		return err
	} else {
		n.delta = uint32(delta)
		n.Name = string(match[2])
	}

	return nil
}
