package lib

type Controller struct {
	ID   byte   `json:"id"`
	Name string `json:"name"`
}

func LookupController(id byte) Controller {
	if c, ok := controllers[id]; ok {
		return c
	}
	return Controller{
		ID:   id,
		Name: "<unknown>",
	}
}

var controllers = map[byte]Controller{
	// High resolution continuous controllers (MSB)
	0:  Controller{0, "Bank Select (MSB)"},
	1:  Controller{1, "Modulation Wheel (MSB)"},
	2:  Controller{2, "Breath Controller (MSB)"},
	4:  Controller{4, "Foot Controller (MSB)"},
	5:  Controller{5, "Portamento Time (MSB)"},
	6:  Controller{6, "Data Entry (MSB)"},
	7:  Controller{7, "Channel Volume (MSB)"},
	8:  Controller{8, "Balance (MSB)"},
	10: Controller{10, "Pan (MSB)"},
	11: Controller{11, "Expression Controller (MSB)"},
	12: Controller{12, "Effect Control 1 (MSB)"},
	13: Controller{13, "Effect Control 2 (MSB)"},
	16: Controller{16, "General Purpose Controller 1 (MSB)"},
	17: Controller{17, "General Purpose Controller 2 (MSB)"},
	18: Controller{18, "General Purpose Controller 3 (MSB)"},
	19: Controller{19, "General Purpose Controller 4 (MSB)"},

	// High resolution continuous controllers (LSB)
	32: Controller{32, "Bank Select (LSB)"},
	33: Controller{33, "Modulation Wheel (LSB)"},
	34: Controller{34, "Breath Controller (LSB)"},
	36: Controller{36, "Foot Controller (LSB)"},
	37: Controller{37, "Portamento Time (LSB)"},
	38: Controller{38, "Data Entry (LSB)"},
	39: Controller{39, "Channel Volume (LSB)"},
	40: Controller{40, "Balance (LSB)"},
	42: Controller{42, "Pan (LSB)"},
	43: Controller{43, "Expression Controller (LSB)"},
	44: Controller{44, "Effect Control 1 (LSB)"},
	45: Controller{45, "Effect Control 2 (LSB)"},
	48: Controller{48, "General Purpose Controller 1 (LSB)"},
	49: Controller{49, "General Purpose Controller 2 (LSB)"},
	50: Controller{50, "General Purpose Controller 3 (LSB)"},
	51: Controller{51, "General Purpose Controller 4 (LSB)"},

	// Switches
	64: Controller{64, "Sustain On/Off"},
	65: Controller{65, "Portamento On/Off"},
	66: Controller{66, "Sostenuto On/Off"},
	67: Controller{67, "Soft Pedal On/Off"},
	68: Controller{68, "Legato On/Off"},
	69: Controller{69, "Hold 2 On/Off"},

	// Low resolution continuous controllers
	70: Controller{70, "Sound Controller 1  (TG: Sound Variation;  FX: Exciter On/Off)"},
	71: Controller{71, "Sound Controller 2  (TG: Harmonic Content; FX: Compressor On/Off)"},
	72: Controller{72, "Sound Controller 3  (TG: Release Time;     FX: Distortion On/Off)"},
	73: Controller{73, "Sound Controller 4  (TG: Attack Time;      FX: EQ On/Off)"},
	74: Controller{74, "Sound Controller 5  (TG: Brightness;       FX: Expander On/Off)"},
	75: Controller{75, "Sound Controller 6  (TG: Decay Time;       FX: Reverb On/Off)"},
	76: Controller{76, "Sound Controller 7  (TG: Vibrato Rate;     FX: Delay On/Off)"},
	77: Controller{77, "Sound Controller 8  (TG: Vibrato Depth;    FX: Pitch Transpose On/Off)"},
	78: Controller{78, "Sound Controller 9  (TG: Vibrato Delay;    FX: Flange/Chorus On/Off)"},
	79: Controller{79, "Sound Controller 10 (TG: Undefined;        FX: Special Effects On/Off)"},
	80: Controller{80, "General Purpose Controller 5"},
	81: Controller{81, "General Purpose Controller 6"},
	82: Controller{82, "General Purpose Controller 7"},
	83: Controller{83, "General Purpose Controller 8"},
	84: Controller{84, "Portamento Control"},
	88: Controller{88, "High Resolution Velocity Prefix"},
	91: Controller{91, "Effects 1 Depth (Reverb Send Level)"},
	92: Controller{92, "Effects 2 Depth (Tremelo Depth)"},
	93: Controller{93, "Effects 3 Depth (Chorus Send Level)"},
	94: Controller{94, "Effects 4 Depth (Celeste Depth)"},
	95: Controller{95, "Effects 5 Depth (Phaser Depth)"},

	// RPNs / NRPNs
	96:  Controller{96, "Data Increment"},
	97:  Controller{97, "Data Decrement"},
	98:  Controller{98, "Non-Registered Parameter Number (LSB)"},
	99:  Controller{99, "Non-Registered Parameter Number (MSB)"},
	100: Controller{100, "Registered Parameter Number (LSB)"},
	101: Controller{101, "Registered Parameter Number (MSB)"},

	// Channel Mode messages
	120: Controller{120, "All Sound Off"},
	121: Controller{121, "Reset All Controllers"},
	122: Controller{122, "Local Control On/Off"},
	123: Controller{123, "All Notes Off"},
	124: Controller{124, "Omni Mode Off"},
	125: Controller{125, "Omni Mode On "},
	126: Controller{126, "Mono Mode On"},
	127: Controller{127, "Poly Mode On"},
}
