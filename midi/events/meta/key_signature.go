package metaevent

import (
	"fmt"
	"github.com/twystd/midiasm/midi/context"
	"io"
)

type KeySignature struct {
	MetaEvent
	Accidentals int8
	KeyType     uint8
}

var major_keys = map[int8]string{
	0:  "C major",
	1:  "G major",
	2:  "D major",
	3:  "A major",
	4:  "E major",
	5:  "B major",
	6:  "F\u266f major",
	-1: "F major",
	-2: "B\u266d major",
	-3: "E\u266d major",
	-4: "A\u266d major",
	-5: "D\u266d major",
	-6: "G\u266d major",
}

var minor_keys = map[int8]string{
	0:  "A minor",
	1:  "E minor",
	2:  "B minor",
	3:  "F\u266f minor",
	4:  "C\u266f minor",
	5:  "G\u266f minor",
	6:  "D\u266f minor",
	-1: "D minor",
	-2: "G minor",
	-3: "C minor",
	-4: "F minor",
	-5: "B\u266d minor",
	-6: "E\u266d minor",
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
	keyType := data[1]
	if keyType != 0 && keyType != 1 {
		return nil, fmt.Errorf("Invalid KeySignature key type (%d): expected a value in the interval [0,1]", keyType)
	}

	return &KeySignature{
		MetaEvent:   *event,
		Accidentals: accidentals,
		KeyType:     keyType,
	}, nil
}

func (e *KeySignature) Render(ctx *context.Context, w io.Writer) {
	if e.Accidentals < 0 {
		ctx.Scale = context.Flats
	} else {
		ctx.Scale = context.Sharps
	}

	switch e.KeyType {
	case 0:
		if signature, ok := major_keys[e.Accidentals]; ok {
			fmt.Fprintf(w, "%s %-16s %s", e.MetaEvent, "KeySignature", signature)
			return
		}
	case 1:
		if signature, ok := minor_keys[e.Accidentals]; ok {
			fmt.Fprintf(w, "%s %-16s %s", e.MetaEvent, "KeySignature", signature)
			return
		}
	}

	fmt.Fprintf(w, "%s %-16s %s", e.MetaEvent, "KeySignature", "???")
}
