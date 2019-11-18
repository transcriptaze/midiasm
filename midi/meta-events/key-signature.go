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

func NewKeySignature(event MetaEvent, data []byte) (*KeySignature, error) {
	if event.eventType != 0x59 {
		return nil, fmt.Errorf("Invalid KeySignature event type (%02x): expected '59'", event.eventType)
	}

	if event.length != 2 {
		return nil, fmt.Errorf("Invalid KeySignature length (%d): expected '2'", event.length)
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
		MetaEvent:   event,
		accidentals: accidentals,
		keyType:     keyType,
	}, nil
}

func (e *KeySignature) Render(w io.Writer) {
	fmt.Fprintf(w, "   ")
	for _, b := range e.bytes {
		fmt.Fprintf(w, "%02x ", b)
	}
	fmt.Fprintf(w, "                                  ")

	fmt.Fprintf(w, "%02x/%-16s %s accidentals:%d key-type:%d\n", e.eventType, "KeySignature", e.MetaEvent.Event, e.accidentals, e.keyType)
}
