package partofspeechtagger

var (
	mockSymbolFrequencyForA = map[string]int64{
		"X": 5,
		"Y": 5,
		"Z": 5,
	}
	mockSymbolFrequencyForB = map[string]int64{
		"X": 3,
		"Y": 10,
		"Z": 0,
	}
	mockSymbolFrequencyForC = map[string]int64{
		"X": 0,
		"Y": 0,
		"Z": 20,
	}

	mockTrigramFrequencyForA = map[Trigram]int64{
		Trigram{
			FirstSymbol:  "X",
			SecondSymbol: "Y",
			ThirdSymbol:  "Z",
		}: 5,
		Trigram{
			FirstSymbol:  "X",
			SecondSymbol: "X",
			ThirdSymbol:  "Z",
		}: 5,
	}
)

type mockHiddenMarkovModel struct{}

func (m mockHiddenMarkovModel) GetAllSymbols() []string {
	return []string{"X", "Y", "Z"}
}

func (m mockHiddenMarkovModel) GetSymbolCountsForToken(token string) map[string]int64 {
	switch token {
	case "A":
		return mockSymbolFrequencyForA
	case "B":
		return mockSymbolFrequencyForB
	case "C":
		return mockSymbolFrequencyForC
	}
}
