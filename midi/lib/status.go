package lib

import (
	"fmt"
)

type Status byte

func (t Status) String() string {
	return fmt.Sprintf("%02X", byte(t))
}
