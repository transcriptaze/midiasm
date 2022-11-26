package metaevent

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/transcriptaze/midiasm/midi/context"
	"github.com/transcriptaze/midiasm/midi/lib"
)

type KeySignature struct {
	event
	Accidentals int8
	KeyType     lib.KeyType
	Key         string
}

func MakeKeySignature(tick uint64, delta lib.Delta, accidentals int8, keyType lib.KeyType, key string, bytes ...byte) KeySignature {
	return KeySignature{
		event: event{
			tick:   tick,
			delta:  delta,
			bytes:  bytes,
			tag:    lib.TagKeySignature,
			Status: 0xff,
			Type:   0x59,
		},
		Accidentals: accidentals,
		KeyType:     keyType,
		Key:         key,
	}
}

func UnmarshalKeySignature(tick uint64, delta lib.Delta, bytes []byte) (*KeySignature, error) {
	if len(bytes) != 2 {
		return nil, fmt.Errorf("Invalid KeySignature length (%d): expected '2'", len(bytes))
	}

	var accidentals = int8(bytes[0])
	var key = ""
	var keyType lib.KeyType

	switch bytes[1] {
	case 0:
		keyType = lib.Major
		if scale, ok := lib.MajorScale(accidentals); !ok {
			return nil, fmt.Errorf("Invalid major key signature (%d accidentals): expected a value in the interval [-6,0]", accidentals)
		} else {
			key = scale.Name
		}

	case 1:
		keyType = lib.Minor
		if scale, ok := lib.MinorScale(accidentals); !ok {
			return nil, fmt.Errorf("Invalid minor key signature (%d accidentals): expected a value in the interval [-6,0]", accidentals)
		} else {
			key = scale.Name
		}

	default:
		return nil, fmt.Errorf("Invalid KeySignature key type (%d): expected a value in the interval [0,1]", keyType)
	}

	event := MakeKeySignature(tick, delta, accidentals, keyType, key, append([]byte{0x00, 0xff, 0x59, 0x02}, bytes...)...)

	return &event, nil
}

func (k *KeySignature) Transpose(ctx *context.Context, steps int) {
	var scale lib.Scale
	var ok bool

	switch k.KeyType {
	case 0:
		if scale, ok = lib.MajorScale(k.Accidentals); !ok {
			return
		}

	case 1:
		if scale, ok = lib.MinorScale(k.Accidentals); !ok {
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
	k.tag = lib.TagKeySignature
	k.Type = 0x59

	re := regexp.MustCompile(`(?i)delta:([0-9]+)(?:.*?)KeySignature\s+([ABCDEFG][♯♭]?)\s+(major|minor)`)
	text := string(bytes)

	if match := re.FindStringSubmatch(text); match == nil || len(match) < 4 {
		return fmt.Errorf("invalid KeySignature event (%v)", text)
	} else if delta, err := lib.ParseDelta(match[1]); err != nil {
		return err
	} else {
		k.delta = delta
		k.Accidentals = 0
		k.KeyType = 0
		k.Key = ""

		key := strings.ToLower(fmt.Sprintf("%v %v", match[2], match[3]))

		for _, scale := range lib.MAJOR_SCALES {
			if strings.ToLower(scale.Name) == key {
				k.Accidentals = scale.Accidentals
				k.KeyType = scale.Type
				k.Key = scale.Name
				break
			}
		}

		for _, scale := range lib.MINOR_SCALES {
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
