package lib

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"strings"
)

type Manufacturer struct {
	ID     []byte `json:"id"`
	Region string `json:"region"`
	Name   string `json:"name"`
}

func AddManufacturers(list map[string]Manufacturer) {
	for k, v := range list {
		manufacturers[k] = v
	}
}

func LookupManufacturer(id []byte) Manufacturer {
	manufacturer := Manufacturer{
		ID:     id,
		Region: "<unknown>",
		Name:   "<unknown>",
	}

	switch len(id) {
	case 1:
		key := fmt.Sprintf("%02X", id[0])
		if m, ok := manufacturers[key]; ok {
			return m
		}

	case 3:
		key := fmt.Sprintf("%02X%02X", id[1], id[2])
		if m, ok := manufacturers[key]; ok {
			return m
		}
	}

	return manufacturer
}

func FindManufacturer(s string) (Manufacturer, error) {
	for _, m := range manufacturers {
		if s == m.Name {
			return m, nil
		}
	}

	return Manufacturer{}, fmt.Errorf("Unrecognised manufacturer %q", s)
}

func LoadManufacturers(r io.Reader) (map[string]Manufacturer, error) {
	conf := struct {
		Manufacturers []Manufacturer `json:"manufacturers"`
	}{
		Manufacturers: []Manufacturer{},
	}

	bytes, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(bytes, &conf)
	if err != nil {
		return nil, err
	}

	list := map[string]Manufacturer{}
	for _, m := range conf.Manufacturers {
		k, v, err := makeManufacturer(m.ID, m.Region, m.Name)
		if err != nil {
			return nil, err
		}

		if v != nil {
			list[k] = *v
		}
	}

	return list, nil
}

func makeManufacturer(ID []byte, region string, name string) (string, *Manufacturer, error) {
	if len(ID) != 1 && len(ID) != 3 {
		return "", nil, fmt.Errorf("Invalid manufacturer ID: %v", ID)
	}

	if strings.Trim(name, " ") == "" {
		return "", nil, fmt.Errorf("Invalid manufacturer name: %s", name)
	}

	key := fmt.Sprintf("%02X", ID[0])
	if len(ID) == 3 {
		key = fmt.Sprintf("%02X%02X", ID[1], ID[2])
	}

	m := Manufacturer{
		ID:     make([]byte, len(ID)),
		Region: strings.Trim(region, " "),
		Name:   strings.Trim(name, " "),
	}

	copy(m.ID, ID)

	return key, &m, nil
}

var manufacturers = map[string]Manufacturer{
	// special purpose

	"7D": Manufacturer{ID: []byte{0x7d}, Region: "Special Purpose", Name: "Non-Commercial"},
	"7E": Manufacturer{ID: []byte{0x7e}, Region: "Special Purpose", Name: "Non-RealTime Extensions"},
	"7F": Manufacturer{ID: []byte{0x7e}, Region: "Special Purpose", Name: "RealTime Extensions"},

	// American
	"01": Manufacturer{ID: []byte{0x01}, Region: "American", Name: "Sequential Circuits"},
	"04": Manufacturer{ID: []byte{0x04}, Region: "American", Name: "Moog"},
	"05": Manufacturer{ID: []byte{0x05}, Region: "American", Name: "Passport Designs"},
	"06": Manufacturer{ID: []byte{0x06}, Region: "American", Name: "Lexicon"},
	"07": Manufacturer{ID: []byte{0x07}, Region: "American", Name: "Kurzweil"},
	"08": Manufacturer{ID: []byte{0x08}, Region: "American", Name: "Fender"},
	"0A": Manufacturer{ID: []byte{0x0a}, Region: "American", Name: "AKG Acoustics"},
	"0F": Manufacturer{ID: []byte{0x0f}, Region: "American", Name: "Ensoniq"},
	"10": Manufacturer{ID: []byte{0x10}, Region: "American", Name: "Oberheim"},
	"11": Manufacturer{ID: []byte{0x11}, Region: "American", Name: "Apple"},
	"13": Manufacturer{ID: []byte{0x13}, Region: "American", Name: "Digidesign"},
	"18": Manufacturer{ID: []byte{0x18}, Region: "American", Name: "Emu"},
	"1A": Manufacturer{ID: []byte{0x1a}, Region: "American", Name: "ART"},
	"1C": Manufacturer{ID: []byte{0x1c}, Region: "American", Name: "Eventide"},

	// European
	"22": Manufacturer{ID: []byte{0x22}, Region: "European", Name: "Synthaxe"},
	"24": Manufacturer{ID: []byte{0x24}, Region: "European", Name: "Hohner"},
	"29": Manufacturer{ID: []byte{0x29}, Region: "European", Name: "PPG"},
	"2B": Manufacturer{ID: []byte{0x2b}, Region: "European", Name: "SSL"},
	"2D": Manufacturer{ID: []byte{0x2d}, Region: "European", Name: "Hinton Instruments"},
	"2F": Manufacturer{ID: []byte{0x2f}, Region: "European", Name: "Elka / General Music"},
	"30": Manufacturer{ID: []byte{0x30}, Region: "European", Name: "Dynacord"},
	"33": Manufacturer{ID: []byte{0x33}, Region: "European", Name: "Clavia (Nord)"},
	"36": Manufacturer{ID: []byte{0x36}, Region: "European", Name: "Cheetah"},
	"3E": Manufacturer{ID: []byte{0x3e}, Region: "European", Name: "Waldorf Electronics Gmbh"},

	// Japanese
	"40": Manufacturer{ID: []byte{0x40}, Region: "Japanese", Name: "Kawai"},
	"41": Manufacturer{ID: []byte{0x41}, Region: "Japanese", Name: "Roland"},
	"42": Manufacturer{ID: []byte{0x42}, Region: "Japanese", Name: "Korg"},
	"43": Manufacturer{ID: []byte{0x43}, Region: "Japanese", Name: "Yamaha"},
	"44": Manufacturer{ID: []byte{0x44}, Region: "Japanese", Name: "Casio"},
	"47": Manufacturer{ID: []byte{0x47}, Region: "Japanese", Name: "Akai"},
	"48": Manufacturer{ID: []byte{0x48}, Region: "Japanese", Name: "Japan Victor (JVC)"},
	"4C": Manufacturer{ID: []byte{0x4c}, Region: "Japanese", Name: "Sony"},
	"4E": Manufacturer{ID: []byte{0x4e}, Region: "Japanese", Name: "Teac Corporation"},
	"51": Manufacturer{ID: []byte{0x51}, Region: "Japanese", Name: "Fostex"},
	"52": Manufacturer{ID: []byte{0x52}, Region: "Japanese", Name: "Zoom"},

	// American
	"0007": Manufacturer{ID: []byte{0x00, 0x00, 0x07}, Region: "American", Name: "Digital Music Corporation"},
	"0009": Manufacturer{ID: []byte{0x00, 0x00, 0x09}, Region: "American", Name: "New England Digital"},
	"000E": Manufacturer{ID: []byte{0x00, 0x00, 0x0e}, Region: "American", Name: "Alesis"},
	"0015": Manufacturer{ID: []byte{0x00, 0x00, 0x15}, Region: "American", Name: "KAT"},
	"0016": Manufacturer{ID: []byte{0x00, 0x00, 0x16}, Region: "American", Name: "Opcode"},
	"001A": Manufacturer{ID: []byte{0x00, 0x00, 0x1a}, Region: "American", Name: "Allen & Heath Brenell"},
	"001B": Manufacturer{ID: []byte{0x00, 0x00, 0x1b}, Region: "American", Name: "Peavey Electronics"},
	"001C": Manufacturer{ID: []byte{0x00, 0x00, 0x1c}, Region: "American", Name: "360 Systems"},
	"001F": Manufacturer{ID: []byte{0x00, 0x00, 0x1f}, Region: "American", Name: "Zeta Systems"},
	"0020": Manufacturer{ID: []byte{0x00, 0x00, 0x20}, Region: "American", Name: "Axxes"},
	"003B": Manufacturer{ID: []byte{0x00, 0x00, 0x3b}, Region: "American", Name: "Mark Of The Unicorn (MOTU)"},
	"004D": Manufacturer{ID: []byte{0x00, 0x00, 0x4d}, Region: "American", Name: "Studio Electronics"},
	"0050": Manufacturer{ID: []byte{0x00, 0x00, 0x50}, Region: "American", Name: "MIDI Solutions Inc"},
	"0137": Manufacturer{ID: []byte{0x00, 0x01, 0x37}, Region: "American", Name: "Roger Linn Design"},
	"0172": Manufacturer{ID: []byte{0x00, 0x01, 0x72}, Region: "American", Name: "Kilpatrick Audio"},
	"0173": Manufacturer{ID: []byte{0x00, 0x01, 0x73}, Region: "American", Name: "iConnectivity"},
	"0214": Manufacturer{ID: []byte{0x00, 0x02, 0x14}, Region: "American", Name: "Intellijel Designs Inc"},

	// // European
	"2011": Manufacturer{ID: []byte{0x00, 0x20, 0x11}, Region: "European", Name: "Forefront Technology"},
	"2013": Manufacturer{ID: []byte{0x00, 0x20, 0x13}, Region: "European", Name: "Kenton Electronics"},
	"201F": Manufacturer{ID: []byte{0x00, 0x20, 0x1f}, Region: "European", Name: "TC Electronic"},
	"2020": Manufacturer{ID: []byte{0x00, 0x20, 0x20}, Region: "European", Name: "Doepfer"},
	"2027": Manufacturer{ID: []byte{0x00, 0x20, 0x27}, Region: "European", Name: "Acorn Computer"},
	"2029": Manufacturer{ID: []byte{0x00, 0x20, 0x29}, Region: "European", Name: "Focusrite / Novation"},
	"2032": Manufacturer{ID: []byte{0x00, 0x20, 0x32}, Region: "European", Name: "Behringer"},
	"2033": Manufacturer{ID: []byte{0x00, 0x20, 0x33}, Region: "European", Name: "Access Music Electronics"},
	"203C": Manufacturer{ID: []byte{0x00, 0x20, 0x3c}, Region: "European", Name: "Elektron"},
	"204D": Manufacturer{ID: []byte{0x00, 0x20, 0x4d}, Region: "European", Name: "Vermona"},
	"2052": Manufacturer{ID: []byte{0x00, 0x20, 0x52}, Region: "European", Name: "Analogue Systems"},
	"2069": Manufacturer{ID: []byte{0x00, 0x20, 0x69}, Region: "European", Name: "Elby Designs"},
	"206B": Manufacturer{ID: []byte{0x00, 0x20, 0x6b}, Region: "European", Name: "Arturia"},
	"2076": Manufacturer{ID: []byte{0x00, 0x20, 0x76}, Region: "European", Name: "Teenage Engineering"},
	"2102": Manufacturer{ID: []byte{0x00, 0x21, 0x02}, Region: "European", Name: "Mutable Instruments"},
	"2109": Manufacturer{ID: []byte{0x00, 0x21, 0x09}, Region: "European", Name: "Native Instruments"},
	"2110": Manufacturer{ID: []byte{0x00, 0x21, 0x10}, Region: "European", Name: "ROLI Ltd"},
	"211A": Manufacturer{ID: []byte{0x00, 0x21, 0x1a}, Region: "European", Name: "IK Multimedia"},
	"211C": Manufacturer{ID: []byte{0x00, 0x21, 0x1c}, Region: "European", Name: "Modor Music"},
	"211D": Manufacturer{ID: []byte{0x00, 0x21, 0x1d}, Region: "European", Name: "Ableton"},
	"2127": Manufacturer{ID: []byte{0x00, 0x21, 0x27}, Region: "European", Name: "Expert Sleepers"},
}
