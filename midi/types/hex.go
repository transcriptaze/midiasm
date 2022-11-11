package types

import (
	"encoding/hex"
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
