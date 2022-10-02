package metaevent

import (
	"fmt"

	"github.com/transcriptaze/midiasm/midi/context"
	"github.com/transcriptaze/midiasm/midi/types"
)

type KeyType uint8

func (k KeyType) String() string {
	if k == 1 {
		return "minor"
	}

	return "major"
}

type KeySignature struct {
	Tag         string
	Status      types.Status
	Type        types.MetaEventType
	Accidentals int8
	KeyType     KeyType
	Key         string
}

func NewKeySignature(ctx *context.Context, bytes []byte) (*KeySignature, error) {
	if len(bytes) != 2 {
		return nil, fmt.Errorf("Invalid KeySignature length (%d): expected '2'", len(bytes))
	}

	accidentals := int8(bytes[0])
	keyType := bytes[1]
	if keyType != 0 && keyType != 1 {
		return nil, fmt.Errorf("Invalid KeySignature key type (%d): expected a value in the interval [0,1]", keyType)
	}

	key := ""
	switch keyType {
	case 0:
		if scale, ok := types.MajorScale(accidentals); !ok {
			return nil, fmt.Errorf("Invalid major key signature (%d accidentals): expected a value in the interval [-6,0]", accidentals)
		} else {
			key = scale.Name
		}
	case 1:
		if scale, ok := types.MinorScale(accidentals); !ok {
			return nil, fmt.Errorf("Invalid minor key signature (%d accidentals): expected a value in the interval [-6,0]", accidentals)
		} else {
			key = scale.Name
		}
	}

	if accidentals < 0 {
		ctx.UseFlats()
	} else {
		ctx.UseSharps()
	}

	return &KeySignature{
		Tag:         "KeySignature",
		Status:      0xff,
		Type:        0x59,
		Accidentals: accidentals,
		KeyType:     KeyType(keyType),
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
