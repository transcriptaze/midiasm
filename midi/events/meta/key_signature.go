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

func (e *KeySignature) unmarshal(tick uint64, delta lib.Delta, status byte, data []byte, bytes ...byte) error {
	if len(data) != 2 {
		return fmt.Errorf("Invalid KeySignature length (%d): expected '2'", len(data))
	}

	var accidentals = int8(data[0])
	var keyType lib.KeyType

	switch data[1] {
	case 0:
		keyType = lib.Major
		if _, ok := lib.MajorScale(accidentals); !ok {
			return fmt.Errorf("Invalid major key signature (%d accidentals): expected a value in the interval [-6..0]", accidentals)
		}

	case 1:
		keyType = lib.Minor
		if _, ok := lib.MinorScale(accidentals); !ok {
			return fmt.Errorf("Invalid minor key signature (%d accidentals): expected a value in the interval [-6..0]", accidentals)
		}

	default:
		return fmt.Errorf("Invalid KeySignature key type (%d): expected a value in the interval [0..1]", keyType)
	}

	// if ctx != nil {
	// 	if accidentals < 0 {
	// 		ctx.UseFlats()
	// 	} else {
	// 		ctx.UseSharps()
	// 	}
	// }

	*e = MakeKeySignature(tick, delta, accidentals, keyType, bytes...)

	return nil
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

func (e *KeySignature) UnmarshalBinary(bytes []byte) error {
	if delta, remaining, err := delta(bytes); err != nil {
		return err
	} else if len(remaining) < 2 {
		return fmt.Errorf("Invalid event (%v)", remaining)
	} else if remaining[0] != 0xff {
		return fmt.Errorf("Invalid %v status (%02X)", lib.TagKeySignature, remaining[0])
	} else if !equals(remaining[1], lib.TypeKeySignature) {
		return fmt.Errorf("Invalid %v event type (%02X)", lib.TagKeySignature, remaining[1])
	} else if data, err := vlf(remaining[2:]); err != nil {
		return err
	} else if len(data) < 2 {
		return fmt.Errorf("Invalid KeySignature data")
	} else {
		var accidentals = int8(data[0])
		var keyType lib.KeyType

		switch data[1] {
		case 0:
			keyType = lib.Major
			if _, ok := lib.MajorScale(accidentals); !ok {
				return fmt.Errorf("Invalid major key signature (%d accidentals): expected a value in the interval [-6..0]", accidentals)
			}

		case 1:
			keyType = lib.Minor
			if _, ok := lib.MinorScale(accidentals); !ok {
				return fmt.Errorf("Invalid minor key signature (%d accidentals): expected a value in the interval [-6..0]", accidentals)
			}

		default:
			return fmt.Errorf("Invalid KeySignature key type (%d): expected a value in the interval [0..1]", keyType)
		}

		*e = MakeKeySignature(0, delta, accidentals, keyType, bytes...)
	}

	return nil
}

func (e *KeySignature) UnmarshalText(text []byte) error {
	re := regexp.MustCompile(`(?i)delta:([0-9]+)(?:.*?)KeySignature\s+([ABCDEFG][♯♭]?)\s+(major|minor)`)

	if match := re.FindStringSubmatch(string(text)); match == nil || len(match) < 4 {
		return fmt.Errorf("invalid KeySignature event (%v)", text)
	} else if delta, err := lib.ParseDelta(match[1]); err != nil {
		return err
	} else {
		accidentals := int8(0)
		keyType := lib.KeyType(0)
		key := strings.ToLower(fmt.Sprintf("%v %v", match[2], match[3]))

		for _, scale := range lib.MAJOR_SCALES {
			if strings.ToLower(scale.Name) == key {
				accidentals = scale.Accidentals
				keyType = scale.Type
				break
			}
		}

		for _, scale := range lib.MINOR_SCALES {
			if strings.ToLower(scale.Name) == key {
				accidentals = scale.Accidentals
				keyType = scale.Type
				break
			}
		}

		*e = MakeKeySignature(0, delta, accidentals, keyType, []byte{}...)
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
		*e = MakeKeySignature(0, t.Delta, t.Accidentals, t.KeyType, []byte{}...)
	}

	return nil
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
