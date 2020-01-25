package types

import (
	"fmt"
)

type Manufacturer struct {
	ID     []byte
	Region string
	Name   string
}

func LookupManufacturer(id []byte) Manufacturer {
	manufacturer := Manufacturer{
		ID:     id,
		Region: "<unknown>",
		Name:   "<unknown>",
	}

	switch len(id) {
	case 1:
		idx := fmt.Sprintf("%02X", id[0])
		if m, ok := manufacturers[idx]; ok {
			manufacturer.Name = m.Name
			manufacturer.Region = m.Region
		}

	case 3:
		idx := fmt.Sprintf("%02X%02X", id[1], id[2])
		if m, ok := manufacturers[idx]; ok {
			manufacturer.Name = m.Name
			manufacturer.Region = m.Region
		}
	}

	return manufacturer
}

var manufacturers = map[string]Manufacturer{
	// American
	"01": Manufacturer{Region: "American", Name: "Sequential Circuits"},
	"04": Manufacturer{Region: "American", Name: "Moog"},
	"05": Manufacturer{Region: "American", Name: "Passport Designs"},
	"06": Manufacturer{Region: "American", Name: "Lexicon"},
	"07": Manufacturer{Region: "American", Name: "Kurzweil"},
	"08": Manufacturer{Region: "American", Name: "Fender"},
	"0A": Manufacturer{Region: "American", Name: "AKG Acoustics"},
	"0F": Manufacturer{Region: "American", Name: "Ensoniq"},
	"10": Manufacturer{Region: "American", Name: "Oberheim"},
	"11": Manufacturer{Region: "American", Name: "Apple"},
	"13": Manufacturer{Region: "American", Name: "Digidesign"},
	"18": Manufacturer{Region: "American", Name: "Emu"},
	"1A": Manufacturer{Region: "American", Name: "ART"},
	"1C": Manufacturer{Region: "American", Name: "Eventide"},

	// European
	"22": Manufacturer{Region: "European", Name: "Synthaxe"},
	"24": Manufacturer{Region: "European", Name: "Hohner"},
	"29": Manufacturer{Region: "European", Name: "PPG"},
	"2B": Manufacturer{Region: "European", Name: "SSL"},
	"2D": Manufacturer{Region: "European", Name: "Hinton Instruments"},
	"2F": Manufacturer{Region: "European", Name: "Elka / General Music"},
	"30": Manufacturer{Region: "European", Name: "Dynacord"},
	"33": Manufacturer{Region: "European", Name: "Clavia (Nord)"},
	"36": Manufacturer{Region: "European", Name: "Cheetah"},
	"3E": Manufacturer{Region: "European", Name: "Waldorf Electronics Gmbh"},

	// Japanese
	"40": Manufacturer{Region: "Japanese", Name: "Kawai"},
	"41": Manufacturer{Region: "Japanese", Name: "Roland"},
	"42": Manufacturer{Region: "Japanese", Name: "Korg"},
	"43": Manufacturer{Region: "Japanese", Name: "Yamaha"},
	"44": Manufacturer{Region: "Japanese", Name: "Casio"},
	"47": Manufacturer{Region: "Japanese", Name: "Akai"},
	"48": Manufacturer{Region: "Japanese", Name: "Japan Victor (JVC)"},
	"4C": Manufacturer{Region: "Japanese", Name: "Sony"},
	"4E": Manufacturer{Region: "Japanese", Name: "Teac Corporation"},
	"51": Manufacturer{Region: "Japanese", Name: "Fostex"},
	"52": Manufacturer{Region: "Japanese", Name: "Zoom"},

	// American
	"0007": Manufacturer{Region: "American", Name: "Digital Music Corporation"},
	"0009": Manufacturer{Region: "American", Name: "New England Digital"},
	"000E": Manufacturer{Region: "American", Name: "Alesis"},
	"0015": Manufacturer{Region: "American", Name: "KAT"},
	"0016": Manufacturer{Region: "American", Name: "Opcode"},
	"001A": Manufacturer{Region: "American", Name: "Allen & Heath Brenell"},
	"001B": Manufacturer{Region: "American", Name: "Peavey Electronics"},
	"001C": Manufacturer{Region: "American", Name: "360 Systems"},
	"001F": Manufacturer{Region: "American", Name: "Zeta Systems"},
	"0020": Manufacturer{Region: "American", Name: "Axxes"},
	"003B": Manufacturer{Region: "American", Name: "Mark Of The Unicorn (MOTU)"},
	"004D": Manufacturer{Region: "American", Name: "Studio Electronics"},
	"0050": Manufacturer{Region: "American", Name: "MIDI Solutions Inc"},
	"0137": Manufacturer{Region: "American", Name: "Roger Linn Design"},
	"0172": Manufacturer{Region: "American", Name: "Kilpatrick Audio"},
	"0173": Manufacturer{Region: "American", Name: "iConnectivity"},
	"0214": Manufacturer{Region: "American", Name: "Intellijel Designs Inc"},

	// // European
	"2011": Manufacturer{Region: "European", Name: "Forefront Technology"},
	"2013": Manufacturer{Region: "European", Name: "Kenton Electronics"},
	"201F": Manufacturer{Region: "European", Name: "TC Electronic"},
	"2020": Manufacturer{Region: "European", Name: "Doepfer"},
	"2027": Manufacturer{Region: "European", Name: "Acorn Computer"},
	"2029": Manufacturer{Region: "European", Name: "Focusrite / Novation"},
	"2032": Manufacturer{Region: "European", Name: "Behringer"},
	"2033": Manufacturer{Region: "European", Name: "Access Music Electronics"},
	"203C": Manufacturer{Region: "European", Name: "Elektron"},
	"204D": Manufacturer{Region: "European", Name: "Vermona"},
	"2052": Manufacturer{Region: "European", Name: "Analogue Systems"},
	"2069": Manufacturer{Region: "European", Name: "Elby Designs"},
	"206B": Manufacturer{Region: "European", Name: "Arturia"},
	"2076": Manufacturer{Region: "European", Name: "Teenage Engineering"},
	"2102": Manufacturer{Region: "European", Name: "Mutable Instruments"},
	"2109": Manufacturer{Region: "European", Name: "Native Instruments"},
	"2110": Manufacturer{Region: "European", Name: "ROLI Ltd"},
	"211A": Manufacturer{Region: "European", Name: "IK Multimedia"},
	"211C": Manufacturer{Region: "European", Name: "Modor Music"},
	"211D": Manufacturer{Region: "European", Name: "Ableton"},
	"2127": Manufacturer{Region: "European", Name: "Expert Sleepers"},
}
