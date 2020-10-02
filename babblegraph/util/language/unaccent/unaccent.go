package unaccent

const (
	lowercaseA rune = 97
	lowercaseC rune = 99
	lowercaseE rune = 101
	lowercaseI rune = 105
	lowercaseO rune = 111
	lowercaseU rune = 117
	lowercaseN rune = 110
	lowercaseY rune = 121

	capitalLetterFactor rune = 32
)

var accentMap = map[rune][]rune{
	lowercaseA: []rune{
		224,
		225,
		226,
		227,
		228,
		229,
		230,
	},
	lowercaseC: []rune{
		231,
	},
	lowercaseE: []rune{
		232,
		233,
		234,
		235,
	},
	lowercaseI: []rune{
		236,
		237,
		238,
		239,
	},
	lowercaseO: []rune{
		240,
		242,
		243,
		244,
		245,
		246,
		248,
	},
	lowercaseN: []rune{
		241,
	},
	lowercaseU: []rune{
		249,
		250,
		251,
		252,
	},
	lowercaseY: []rune{
		253,
		254,
		255,
	},
}

func UnaccentRune(in rune) rune {
	var capitalizationCorrection rune
	if in >= 192 && in <= 222 {
		capitalizationCorrection = 32
		in -= capitalLetterFactor
	}
	for out, mappings := range accentMap {
		for _, r := range mappings {
			if r == in {
				return out + capitalizationCorrection
			}
		}
	}
	return in
}
