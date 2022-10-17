package metaevent

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/transcriptaze/midiasm/midi/context"
	"github.com/transcriptaze/midiasm/midi/types"
)

type KeySignature struct {
	event
	Accidentals int8
	KeyType     types.KeyType
	Key         string
}

func NewKeySignature(ctx *context.Context, tick uint64, delta uint32, bytes []byte) (*KeySignature, error) {
	if len(bytes) != 2 {
		return nil, fmt.Errorf("Invalid KeySignature length (%d): expected '2'", len(bytes))
	}

	var accidentals = int8(bytes[0])
	var key = ""
	var keyType types.KeyType

	switch bytes[1] {
	case 0:
		keyType = types.Major
		if scale, ok := types.MajorScale(accidentals); !ok {
			return nil, fmt.Errorf("Invalid major key signature (%d accidentals): expected a value in the interval [-6,0]", accidentals)
		} else {
			key = scale.Name
		}

	case 1:
		keyType = types.Minor
		if scale, ok := types.MinorScale(accidentals); !ok {
			return nil, fmt.Errorf("Invalid minor key signature (%d accidentals): expected a value in the interval [-6,0]", accidentals)
		} else {
			key = scale.Name
		}

	default:
		return nil, fmt.Errorf("Invalid KeySignature key type (%d): expected a value in the interval [0,1]", keyType)
	}

	if ctx != nil {
		if accidentals < 0 {
			ctx.UseFlats()
		} else {
			ctx.UseSharps()
		}
	}

	return &KeySignature{
		event: event{
			tick:  tick,
			delta: delta,
			bytes: concat([]byte{0x00, 0xff, 0x59, 0x02}, bytes),

			Tag:    "KeySignature",
			Status: 0xff,
			Type:   0x59,
		},
		Accidentals: accidentals,
		KeyType:     types.KeyType(keyType),
		Key:         key,
	}, nil
}

func (k *KeySignature) Transpose(ctx *context.Context, steps int) {
	var scale types.Scale
	var ok bool

	switch k.KeyType {
	case 0:
		if scale, ok = types.MajorScale(k.Accidentals); !ok {
			return
		}

	case 1:
		if scale, ok = types.MinorScale(k.Accidentals); !ok {
			return
		}

	default:
		return
	}

	transposed := scale.Transpose(steps)

	k.Key = transposed.Name
	k.Accidentals = transposed.Accidentals

	if k.Accidentals < 0 {
		ctx.UseFlats()
	} else {
		ctx.UseSharps()
	}
}

func (k KeySignature) MarshalBinary() (encoded []byte, err error) {
	encoded = make([]byte, 5)

	encoded[0] = byte(k.Status)
	encoded[1] = byte(k.Type)
	encoded[2] = byte(2)
	encoded[3] = byte(k.Accidentals)
	encoded[4] = byte(k.KeyType)

	return
}

func (k *KeySignature) UnmarshalText(bytes []byte) error {
	k.tick = 0
	k.delta = 0
	k.bytes = []byte{}
	k.Status = 0xff
	k.Tag = "KeySignature"
	k.Type = 0x59

	re := regexp.MustCompile(`(?i)KeySignature\s+([ABCDEFG][♯♭]?)\s+(major|minor)`)
	text := string(bytes)

	if match := re.FindStringSubmatch(text); match == nil || len(match) < 3 {
		return fmt.Errorf("invalid KeySignature event (%v)", text)
	} else {
		k.Accidentals = 0
		k.KeyType = 0
		k.Key = ""

		key := strings.ToLower(fmt.Sprintf("%v %v", match[1], match[2]))

		for _, scale := range types.MAJOR_SCALES {
			if strings.ToLower(scale.Name) == key {
				k.Accidentals = scale.Accidentals
				k.KeyType = scale.Type
				k.Key = scale.Name
				break
			}
		}

		for _, scale := range types.MINOR_SCALES {
			if strings.ToLower(scale.Name) == key {
				k.Accidentals = scale.Accidentals
				k.KeyType = scale.Type
				k.Key = scale.Name
				break
			}
		}
	}

	return nil
}
