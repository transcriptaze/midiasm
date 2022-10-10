package assemble

import (
	"bufio"
	"bytes"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/transcriptaze/midiasm/midi"
)

type TextAssembler struct {
}

func NewTextAssembler() TextAssembler {
	return TextAssembler{}
}

func (a TextAssembler) Assemble(source []byte) ([]byte, error) {
	lines, err := a.read(source)
	if err != nil {
		return nil, err
	}

	// ... MThd

	var mthd *midi.MThd

	for _, line := range lines {
		if strings.Contains(line, "MThd") {
			var format uint16
			var ppqn uint16

			if match := regexp.MustCompile(`format:(0|1|2)`).FindStringSubmatch(line); match == nil || len(match) < 2 {
				return nil, fmt.Errorf("missing or invalid 'format' field in MThd")
			} else if v, err := strconv.ParseUint(match[1], 10, 16); err != nil {
				return nil, err
			} else {
				format = uint16(v)
			}

			if match := regexp.MustCompile(`metrical(?:[ -])?time:([0-9]+)\s*ppqn`).FindStringSubmatch(line); match == nil || len(match) < 2 {
				return nil, fmt.Errorf("missing 'metrical-time' field in MThd")
			} else if v, err := strconv.ParseUint(match[1], 10, 16); err != nil {
				return nil, err
			} else {
				ppqn = uint16(v)
			}

			if mthd, err = midi.NewMThd(format, 0, ppqn); err != nil {
				return nil, err
			}

			break
		}
	}

	tracks := make([]*midi.MTrk, 0)

	smf := midi.SMF{
		MThd:   mthd,
		Tracks: tracks,
	}

	return smf.MarshalBinary()
}

func (a TextAssembler) read(source []byte) ([]string, error) {
	lines := []string{}

	r := bytes.NewBuffer(source)
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return lines, nil
}
