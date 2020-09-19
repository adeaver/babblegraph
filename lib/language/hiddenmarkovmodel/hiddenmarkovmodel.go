package hiddenmarkovmodel

import "fmt"

type model struct{}

type observance struct {
	observedItem string
	frequency
}

type TokenToSymbolCount struct {
	Token string
	Count int64
}

type SymbolTrigramWithCount struct {
	FirstSymbol  string
	SecondSymbol string
	ThirdSymbols string
	Count        int64
}

func (s SymbolTrigramWithCount) makeBigramKey() string {
	return fmt.Sprintf("%s,%s", s.FirstSymbol, s.SecondSymbol)
}

func GetSymbolSequence(sequence []string) ([]string, error) {

}
