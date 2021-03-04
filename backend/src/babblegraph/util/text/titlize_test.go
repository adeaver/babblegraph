package text

import (
	"babblegraph/wordsmith"
	"testing"
)

func TestToTitleCase(t *testing.T) {
	type testCase struct {
		inputText     string
		inputLanguage wordsmith.LanguageCode
		expectedText  string
	}
	testCases := []testCase{
		{
			inputText:     "noticias de ciencia",
			inputLanguage: wordsmith.LanguageCodeSpanish,
			expected:      "Noticias de Cienca",
		}, {
			inputText:     "de ciencia",
			inputLanguage: wordsmith.LanguageCodeSpanish,
			expected:      "De Cienca",
		}, {
			inputText:     "noticias del futuro",
			inputLanguage: wordsmith.LanguageCodeSpanish,
			expected:      "Noticias del Futuro",
		}, {
			inputText:     "noticIAS del futuro",
			inputLanguage: wordsmith.LanguageCodeSpanish,
			expected:      "Noticias del Futuro",
		}, {
			inputText:     "vamos al futuro",
			inputLanguage: wordsmith.LanguageCodeSpanish,
			expected:      "Vamos al Futuro",
		}, {
			inputText:     "al futuro",
			inputLanguage: wordsmith.LanguageCodeSpanish,
			expected:      "Al Futuro",
		},
	}
	for idx, tc := range testCases {
		result := ToTitleCaseForLanguage(tc.inputText, tc.inputLanguage)
		if result != tc.expected {
			t.Errorf("Error on test case %d: expected %s but got %s", idx+1, tc.expected, result)
		}
	}
}
