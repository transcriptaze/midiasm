package lib

import (
	"fmt"
)

// Ref. https://computermusicresource.com/midikeys.html
type MiddleC int

const (
	C3 MiddleC = iota
	C4
)

func (c MiddleC) String() string {
	return [...]string{"C3", "C4"}[c]
}

func (c *MiddleC) Set(s string) error {
	switch s {
	case "C3", "c3":
		*c = C3
		return nil

	case "C4", "c4":
		*c = C4
		return nil

	}

	return fmt.Errorf("invalid middle C convention - expected C3 or C4, got:%v", s)
}
