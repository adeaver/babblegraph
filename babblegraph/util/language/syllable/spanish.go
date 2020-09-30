package syllable

import "fmt"

const accentFactor rune = 128

var strongVowels = map[rune]bool{
	97:  true, // a
	101: true, // e
	111: true, // o
}

var weakVowels = map[rune]bool{
	105: true, // i
	117: true, // u
}

func countSyllablesForSpanish(word string) (*int64, error) {
	wordAsRunes := []rune(word)
	var currentSyllable []rune
	var syllableCount int64
	for idx, r := range wordAsRunes {
		if r >= 223 && r <= 246 {
			r -= accentFactor
		}
		_, isStrongVowel := strongVowels[r]
		_, isWeakVowel := weakVowels[r]
		switch {
		case isStrongVowel:
			if containsStrongVowel(currentSyllable) {
				syllableCount++
				currentSyllable = []rune{}
			}
			currentSyllable = append(currentSyllable, r)
		case isWeakVowel:
			var isPreviousRuneStrongVowel, isNextRuneStrongVowel bool
			if len(currentSyllable) > 0 {
				previousRune := currentSyllable[len(currentSyllable)-1]
				_, isPreviousRuneStrongVowel = strongVowels[previousRune]
			}
			if idx < len(wordAsRunes)-1 {
				nextRune := wordAsRunes[idx+1]
				_, isNextRuneStrongVowel = strongVowels[nextRune]
			}
			currentSyllableContainsVowel := containsStrongVowel(currentSyllable) || containsWeakVowel(currentSyllable)
			switch {
			case !isNextRuneStrongVowel && !isPreviousRuneStrongVowel && currentSyllableContainsVowel:
				// runes are not vowels
				syllableCount++
				currentSyllable = append([]rune{}, r)
			case isPreviousRuneStrongVowel:
				// no-op
			default:
				currentSyllable = append(currentSyllable, r)
			}
		case r >= 97 && r <= 122:
			if containsStrongVowel(currentSyllable) || containsWeakVowel(currentSyllable) {
				syllableCount++
				currentSyllable = []rune{}
			}
			currentSyllable = append(currentSyllable, r)
		default:
			return nil, fmt.Errorf("expected lowercase word")
		}
	}
	if containsStrongVowel(currentSyllable) || containsWeakVowel(currentSyllable) {
		syllableCount++
	}
	return &syllableCount, nil
}

func containsStrongVowel(runes []rune) bool {
	for _, r := range runes {
		if _, isStrongVowel := strongVowels[r]; isStrongVowel {
			return true
		}
	}
	return false
}

func containsWeakVowel(runes []rune) bool {
	for _, r := range runes {
		if _, isWeakVowel := weakVowels[r]; isWeakVowel {
			return true
		}
	}
	return false
}
