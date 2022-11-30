package metaevent

import (
	"encoding/json"
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

func MakeKeySignature(tick uint64, delta lib.Delta, accidentals int8, keyType lib.KeyType, bytes ...byte) KeySignature {
	key := ""

	switch keyType {
	case lib.Major:
		if scale, ok := lib.MajorScale(accidentals); !ok {
			panic(fmt.Errorf("Invalid major key signature (%d accidentals): expected a value in the interval [-6,0]", accidentals))
		} else {
			key = scale.Name
		}

	case lib.Minor:
		if scale, ok := lib.MinorScale(accidentals); !ok {
			panic(fmt.Errorf("Invalid minor key signature (%d accidentals): expected a value in the interval [-6,0]", accidentals))
		} else {
			key = scale.Name
		}

	default:
		panic(fmt.Errorf("Invalid KeySignature key type (%d): expected a value in the interval [0..1]", keyType))
	}

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

func UnmarshalKeySignature(tick uint64, delta lib.Delta, data ...byte) (*KeySignature, error) {
	if len(data) != 2 {
		return nil, fmt.Errorf("Invalid KeySignature length (%d): expected '2'", len(data))
	}

	var accidentals = int8(data[0])
	var keyType lib.KeyType

	switch data[1] {
	case 0:
		keyType = lib.Major
		if _, ok := lib.MajorScale(accidentals); !ok {
			return nil, fmt.Errorf("Invalid major key signature (%d accidentals): expected a value in the interval [-6..0]", accidentals)
		}

	case 1:
		keyType = lib.Minor
		if _, ok := lib.MinorScale(accidentals); !ok {
			return nil, fmt.Errorf("Invalid minor key signature (%d accidentals): expected a value in the interval [-6..0]", accidentals)
		}

	default:
		return nil, fmt.Errorf("Invalid KeySignature key type (%d): expected a value in the interval [0..1]", keyType)
	}

	event := MakeKeySignature(tick, delta, accidentals, keyType)

	return &event, nil
}

func (e KeySignature) Transpose(ctx *context.Context, steps int) KeySignature {
	var scale lib.Scale
	var ok bool

	switch e.KeyType {
	case 0:
		if scale, ok = lib.MajorScale(e.Accidentals); !ok {
			panic(fmt.Sprintf("Unknown major scale (%v)", e.Accidentals))
		}

	case 1:
		if scale, ok = lib.MinorScale(e.Accidentals); !ok {
			panic(fmt.Sprintf("Unknown minor scale (%v)", e.Accidentals))
		}

	default:
		panic(fmt.Sprintf("Unknown key type (%v)", e.KeyType))
	}

	transposed := scale.Transpose(steps)

	keyType := e.KeyType
	key := transposed.Name
	accidentals := transposed.Accidentals

	if accidentals < 0 {
		ctx.UseFlats()
	} else {
		ctx.UseSharps()
	}

	return KeySignature{
		event: event{
			tick:   e.tick,
			delta:  e.delta,
			tag:    lib.TagKeySignature,
			Status: 0xff,
			Type:   lib.TypeKeySignature,
		},
		Accidentals: accidentals,
		KeyType:     keyType,
		Key:         key,
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

func (e KeySignature) MarshalJSON() (encoded []byte, err error) {
	t := struct {
		Tag         string      `json:"tag"`
		Delta       lib.Delta   `json:"delta"`
		Status      byte        `json:"status"`
		Type        byte        `json:"type"`
		Accidentals int8        `json:"accidentals"`
		KeyType     lib.KeyType `json:"key-type"`
		Key         string      `json:"key"`
	}{
		Tag:         fmt.Sprintf("%v", e.tag),
		Delta:       e.delta,
		Status:      byte(e.Status),
		Type:        byte(e.Type),
		Accidentals: e.Accidentals,
		KeyType:     e.KeyType,
		Key:         e.Key,
	}

	return json.Marshal(t)
}

func (e *KeySignature) UnmarshalJSON(bytes []byte) error {
	t := struct {
		Tag         string      `json:"tag"`
		Delta       lib.Delta   `json:"delta"`
		Accidentals int8        `json:"accidentals"`
		KeyType     lib.KeyType `json:"key-type"`
	}{}

	if err := json.Unmarshal(bytes, &t); err != nil {
		return err
	} else if !equal(t.Tag, lib.TagKeySignature) {
		return fmt.Errorf("invalid %v event (%v)", e.tag, string(bytes))
	} else {
		key := ""

		switch t.KeyType {
		case lib.Major:
			if scale, ok := lib.MajorScale(t.Accidentals); !ok {
				return fmt.Errorf("Invalid major key signature (%d accidentals): expected a value in the interval [-6,0]", t.Accidentals)
			} else {
				key = scale.Name
			}

		case lib.Minor:
			if scale, ok := lib.MinorScale(t.Accidentals); !ok {
				return fmt.Errorf("Invalid minor key signature (%d accidentals): expected a value in the interval [-6,0]", t.Accidentals)
			} else {
				key = scale.Name
			}

		default:
			return fmt.Errorf("Invalid key type (%v)", t.KeyType)
		}

		e.tick = 0
		e.delta = t.Delta
		e.bytes = []byte{}
		e.Status = 0xff
		e.tag = lib.TagKeySignature
		e.Type = lib.TypeKeySignature
		e.Accidentals = t.Accidentals
		e.KeyType = t.KeyType
		e.Key = key
	}

	return nil
}
