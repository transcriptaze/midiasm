package types

import (
	"fmt"
	"regexp"
	"strconv"
)

type Delta uint32

func (d Delta) String() string {
	return fmt.Sprintf("%d", uint32(d))
}

func (d *Delta) UnmarshalText(bytes []byte) error {
	re := regexp.MustCompile("delta:([0-9]+)")

	if match := re.FindStringSubmatch(string(bytes)); match == nil || len(match) != 2 {
		return fmt.Errorf("missing delta field")
	} else if delta, err := strconv.ParseUint(match[1], 10, 32); err != nil {
		return err
	} else {
		*d = Delta(delta)
	}

	return nil
}
