package lib

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
)

type Hex []byte

func ParseHex(s string) (Hex, error) {
	v := regexp.MustCompile(`\s+`).ReplaceAllString(s, "")
	if bytes, err := hex.DecodeString(v); err != nil {
		return Hex([]byte{}), err
	} else {
		return Hex(bytes), nil
	}
}

func (bytes Hex) String() string {
	s := ""
	for _, b := range bytes {
		s += fmt.Sprintf("%02X ", b)
	}

	return strings.TrimSpace(s)
}

func (h Hex) MarshalBinary() ([]byte, error) {
	return vlf2bin([]byte(h)), nil
}

func (h Hex) MarshalJSON() ([]byte, error) {
	bytes := []uint{}
	for _, b := range h {
		bytes = append(bytes, uint(b))
	}

	return json.Marshal(bytes)
}

func (h *Hex) UnmarshalJSON(bytes []byte) error {
	v := []uint{}

	if err := json.Unmarshal(bytes, &v); err != nil {
		return err
	}

	buffer := []byte{}
	for _, b := range v {
		buffer = append(buffer, byte(b))
	}

	*h = buffer

	return nil
}
