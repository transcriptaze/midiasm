package assemble

import (
	"bufio"
	"fmt"
	"io"
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

func (a TextAssembler) Assemble(r io.Reader) ([]byte, error) {
	chunks, err := a.read(r)
	if err != nil {
		return nil, err
	}

	// ... parse chunks

	smf := midi.SMF{}

	for _, chunk := range chunks {
		for _, line := range chunk {
			switch {
			case strings.Contains(line, "MThd"):
				if mthd, err := a.parseMThd(chunk); err != nil {
					return nil, err
				} else {
					smf.MThd = mthd
				}

				break

			case strings.Contains(line, "MTrk"):
				if mtrk, err := a.parseMTrk(chunk); err != nil {
					return nil, err
				} else {
					smf.MThd.Tracks += 1
					smf.Tracks = append(smf.Tracks, mtrk)
				}

				break
			}
		}
	}

	// ... 'k, done

	return smf.MarshalBinary()
}

func (a TextAssembler) read(r io.Reader) ([][]string, error) {
	scanner := bufio.NewScanner(r)
	lines := make(chan string)
	chunks := make(chan []string)

	go a.scan(scanner, lines)
	go a.chunkify(lines, chunks)

	list := [][]string{}
	for chunk := range chunks {
		list = append(list, chunk)
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return list, nil
}

func (a TextAssembler) scan(scanner *bufio.Scanner, lines chan string) {
	for scanner.Scan() {
		lines <- scanner.Text()
	}

	close(lines)
}

func (a TextAssembler) chunkify(lines chan string, chunks chan []string) {
	tags := regexp.MustCompile("(MThd)|(MTrk)")

	clone := func(slice []string) []string {
		chunk := make([]string, len(slice))
		copy(chunk, slice)
		return chunk
	}

	var chunk []string
	for line := range lines {
		if strings.Contains(line, "MThd") {
			chunk = []string{line}
			break
		}
	}

	for line := range lines {
		if match := tags.FindStringSubmatch(line); match != nil {
			chunks <- clone(chunk)
			chunk = []string{line}
		} else {
			chunk = append(chunk, line)
		}
	}

	if len(chunk) > 0 {
		chunks <- clone(chunk)
	}

	close(chunks)
}

func (a TextAssembler) parseMThd(chunk []string) (*midi.MThd, error) {
	for _, line := range chunk {
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

			if mthd, err := midi.NewMThd(format, 0, ppqn); err != nil {
				return nil, err
			} else {
				return mthd, nil
			}
		}
	}

	return nil, fmt.Errorf("invalid MThd")
}

func (a TextAssembler) parseMTrk(chunk []string) (*midi.MTrk, error) {
	for _, line := range chunk {
		if strings.Contains(line, "MTrk") {
			if mtrk, err := midi.NewMTrk(); err != nil {
				return nil, err
			} else {
				return mtrk, nil
			}
		}
	}

	return nil, fmt.Errorf("invalid MTrk")
}
