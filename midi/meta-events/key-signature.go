package metaevent

import (
	"fmt"
	"io"
)

type KeySignature struct {
	MetaEvent
	accidentals int8
	keyType     uint8
}

func NewKeySignature(event *MetaEvent, r io.ByteReader) (*KeySignature, error) {
	if event.Type != 0x59 {
		return nil, fmt.Errorf("Invalid KeySignature event type (%02x): expected '59'", event.Type)
	}

	data, err := read(r)
	if err != nil {
		return nil, err
	} else if len(data) != 2 {
		return nil, fmt.Errorf("Invalid KeySignature length (%d): expected '2'", len(data))
	}

	accidentals := int8(data[0])
	if accidentals < -7 || accidentals > +7 {
		return nil, fmt.Errorf("Invalid KeySignature accidentals (%d): expected a value in the interval [-7,+7]", accidentals)
	}

	keyType := data[1]
	if keyType != 0 && keyType != 1 {
		return nil, fmt.Errorf("Invalid KeySignature key type (%d): expectedi a value in the interval [0,1]", keyType)
	}

	return &KeySignature{
		MetaEvent:   *event,
		accidentals: accidentals,
		keyType:     keyType,
	}, nil
}

func (e *KeySignature) Render(w io.Writer) {
	fmt.Fprintf(w, "%s %-16s accidentals:%d key-type:%d", e.MetaEvent, "KeySignature", e.accidentals, e.keyType)
}
