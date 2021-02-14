package spanishprocessing

import (
	"babblegraph/util/math/decimal"
	"babblegraph/wordsmith"
	"fmt"
	"testing"
)

func makeSampleWordsmithWord(text string, id int) wordsmith.Word {
	return wordsmith.Word{
		ID:             wordsmith.WordID(fmt.Sprintf("%s-%d", text, id)),
		Language:       wordsmith.LanguageCodeSpanish,
		CorpusID:       wordsmith.SpanishUPCWikiCorpus,
		PartOfSpeechID: wordsmith.PartOfSpeechID(fmt.Sprintf("%s-%d", text, id)),
		LemmaID:        wordsmith.LemmaID(fmt.Sprintf("%s-%d", text, id)),
		WordText:       text,
	}
}

func makeSampleWordsmithWordBigramCount(word1, word2 wordsmith.Word, count int64) wordsmith.WordBigramCount {
	return wordsmith.WordBigramCount{
		ID:       wordsmith.WordBigramCountID(fmt.Sprintf("%s-%s", word1.ID, word2.ID)),
		Language: wordsmith.LanguageCodeSpanish,
		CorpusID: wordsmith.SpanishUPCWikiCorpus,
		FirstWord: wordsmith.BigramWord{
			Text:    word1.WordText,
			LemmaID: word1.LemmaID,
		},
		SecondWord: wordsmith.BigramWord{
			Text:    word2.WordText,
			LemmaID: word2.LemmaID,
		},
		Count: count,
	}
}

func TestPickBestWordUsingBigrams(t *testing.T) {
	type testCase struct {
		input    pickBestWordUsingBigramInput
		expected wordsmith.WordID
	}
	testCases := []testCase{
		{
			input: pickBestWordUsingBigramInput{
				wordChoices: []wordsmith.Word{
					makeSampleWordsmithWord("hola", 1),
					makeSampleWordsmithWord("hola", 2),
				},
				bigramCountsEndingInToken: []wordsmith.WordBigramCount{
					makeSampleWordsmithWordBigramCount(makeSampleWordsmithWord("adios", 1), makeSampleWordsmithWord("hola", 1), 10),
					makeSampleWordsmithWordBigramCount(makeSampleWordsmithWord("adios", 2), makeSampleWordsmithWord("hola", 1), 10),
					makeSampleWordsmithWordBigramCount(makeSampleWordsmithWord("adios", 1), makeSampleWordsmithWord("hola", 2), 1),
				},
				bigramCountsStartingInToken: []wordsmith.WordBigramCount{
					makeSampleWordsmithWordBigramCount(makeSampleWordsmithWord("hola", 1), makeSampleWordsmithWord("placer", 1), 10),
					makeSampleWordsmithWordBigramCount(makeSampleWordsmithWord("hola", 2), makeSampleWordsmithWord("placer", 1), 10),
				},
			},
			expected: wordsmith.WordID("hola-1"),
		}, {
			input: pickBestWordUsingBigramInput{
				wordChoices: []wordsmith.Word{
					makeSampleWordsmithWord("hola", 1),
					makeSampleWordsmithWord("hola", 2),
				},
				bigramCountsEndingInToken: []wordsmith.WordBigramCount{},
				bigramCountsStartingInToken: []wordsmith.WordBigramCount{
					makeSampleWordsmithWordBigramCount(makeSampleWordsmithWord("hola", 1), makeSampleWordsmithWord("placer", 1), 10),
					makeSampleWordsmithWordBigramCount(makeSampleWordsmithWord("hola", 2), makeSampleWordsmithWord("placer", 1), 10),
				},
			},
			expected: wordsmith.WordID("hola-1"),
		}, {
			input: pickBestWordUsingBigramInput{
				wordChoices: []wordsmith.Word{
					makeSampleWordsmithWord("hola", 1),
					makeSampleWordsmithWord("hola", 2),
				},
				bigramCountsEndingInToken: []wordsmith.WordBigramCount{},
				bigramCountsStartingInToken: []wordsmith.WordBigramCount{
					makeSampleWordsmithWordBigramCount(makeSampleWordsmithWord("hola", 2), makeSampleWordsmithWord("placer", 1), 10),
				},
			},
			expected: wordsmith.WordID("hola-2"),
		}, {
			input: pickBestWordUsingBigramInput{
				wordChoices: []wordsmith.Word{
					makeSampleWordsmithWord("hola", 1),
					makeSampleWordsmithWord("hola", 2),
				},
				bigramCountsEndingInToken:   []wordsmith.WordBigramCount{},
				bigramCountsStartingInToken: []wordsmith.WordBigramCount{},
			},
			expected: wordsmith.WordID("hola-1"),
		}, {
			input: pickBestWordUsingBigramInput{
				wordChoices: []wordsmith.Word{
					makeSampleWordsmithWord("hola", 1),
				},
				bigramCountsEndingInToken:   []wordsmith.WordBigramCount{},
				bigramCountsStartingInToken: []wordsmith.WordBigramCount{},
			},
			expected: wordsmith.WordID("hola-1"),
		},
	}
	for idx, tc := range testCases {
		result := pickBestWordUsingBigrams(tc.input)
		if result.ID != tc.expected {
			t.Errorf("Error on test case %d: expected %s, but got %s", idx+1, tc.expected, result.ID)
		}
	}
}

func TestCalculateProbabilityOfEndingInToken(t *testing.T) {
	type testCase struct {
		word             wordsmith.Word
		wordBigramCounts []wordsmith.WordBigramCount
		expected         decimal.Number
	}
	testCases := []testCase{
		{
			word: makeSampleWordsmithWord("hola", 1),
			wordBigramCounts: []wordsmith.WordBigramCount{
				makeSampleWordsmithWordBigramCount(makeSampleWordsmithWord("adios", 1), makeSampleWordsmithWord("hola", 1), 25),
				makeSampleWordsmithWordBigramCount(makeSampleWordsmithWord("adios", 2), makeSampleWordsmithWord("hola", 1), 25),
				makeSampleWordsmithWordBigramCount(makeSampleWordsmithWord("adios", 1), makeSampleWordsmithWord("hola", 2), 50),
			},
			expected: decimal.FromFloat64(51.0 / 100.0),
		}, {
			word: makeSampleWordsmithWord("hola", 2),
			wordBigramCounts: []wordsmith.WordBigramCount{
				makeSampleWordsmithWordBigramCount(makeSampleWordsmithWord("adios", 1), makeSampleWordsmithWord("hola", 1), 25),
				makeSampleWordsmithWordBigramCount(makeSampleWordsmithWord("adios", 2), makeSampleWordsmithWord("hola", 1), 25),
			},
			expected: decimal.FromFloat64(1.0 / 50.0),
		}, {
			word:             makeSampleWordsmithWord("hola", 2),
			wordBigramCounts: []wordsmith.WordBigramCount{},
			expected:         decimal.FromFloat64(1),
		},
	}
	for idx, tc := range testCases {
		result := calculateProbabilityOfEndingInToken(tc.word, tc.wordBigramCounts)
		if !tc.expected.EqualTo(result) {
			t.Errorf("Error on test case %d: expected probability of %f, but got %f", idx+1, tc.expected.ToFloat64(), result.ToFloat64())
		}
	}
}

func TestCalculateProbabilityOfStartingInToken(t *testing.T) {
	type testCase struct {
		word             wordsmith.Word
		wordBigramCounts []wordsmith.WordBigramCount
		expected         decimal.Number
	}
	testCases := []testCase{
		{
			word: makeSampleWordsmithWord("hola", 1),
			wordBigramCounts: []wordsmith.WordBigramCount{
				makeSampleWordsmithWordBigramCount(makeSampleWordsmithWord("hola", 1), makeSampleWordsmithWord("adios", 1), 25),
				makeSampleWordsmithWordBigramCount(makeSampleWordsmithWord("hola", 1), makeSampleWordsmithWord("adios", 2), 25),
				makeSampleWordsmithWordBigramCount(makeSampleWordsmithWord("hola", 2), makeSampleWordsmithWord("adios", 1), 50),
			},
			expected: decimal.FromFloat64(51.0 / 100.0),
		}, {
			word: makeSampleWordsmithWord("hola", 2),
			wordBigramCounts: []wordsmith.WordBigramCount{
				makeSampleWordsmithWordBigramCount(makeSampleWordsmithWord("hola", 1), makeSampleWordsmithWord("adios", 1), 25),
				makeSampleWordsmithWordBigramCount(makeSampleWordsmithWord("hola", 1), makeSampleWordsmithWord("adios", 2), 25),
			},
			expected: decimal.FromFloat64(1.0 / 50.0),
		}, {
			word:             makeSampleWordsmithWord("hola", 2),
			wordBigramCounts: []wordsmith.WordBigramCount{},
			expected:         decimal.FromFloat64(1),
		},
	}
	for idx, tc := range testCases {
		result := calculateProbabilityOfStartingWithToken(tc.word, tc.wordBigramCounts)
		if !tc.expected.EqualTo(result) {
			t.Errorf("Error on test case %d: expected probability of %f, but got %f", idx+1, tc.expected.ToFloat64(), result.ToFloat64())
		}
	}
}
