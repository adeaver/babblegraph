package textprocessing

import (
	"babblegraph/model/documents"
	"babblegraph/util/testutils"
	"babblegraph/wordsmith"
	"fmt"
	"testing"
)

func TestTokenizeText(t *testing.T) {
	type testCase struct {
		normalizedText string
		expected       []string
	}
	testCases := []testCase{
		{
			normalizedText: "uno dos tres",
			expected:       []string{"uno", "dos", "tres"},
		}, {
			normalizedText: "uno    dos",
			expected:       []string{"uno", "dos"},
		}, {
			normalizedText: `
            uno uno

            dos   dos
            tres            tres


            cuatro`,
			expected: []string{"uno", "uno", "dos", "dos", "tres", "tres", "cuatro"},
		},
	}
	for idx, tc := range testCases {
		result := tokenizeText(tc.normalizedText)
		if err := testutils.CompareStringLists(result, tc.expected); err != nil {
			t.Errorf("Error on test case %d: %s", idx+1, err.Error())
		}
	}
}

func TestGetUniqueWordsForTest(t *testing.T) {
	type testCase struct {
		tokenizedText []string
		expected      []string
	}
	testCases := []testCase{
		{
			tokenizedText: []string{"uno", "uno", "uno", "dos", "uno", "dos", "dos", "tres"},
			expected:      []string{"uno", "dos", "tres"},
		}, {
			tokenizedText: []string{"dos uno uno uno uno uno dos"},
			expected:      []string{"dos", "uno"},
		}, {
			tokenizedText: []string{""},
			expected:      []string{},
		}, {
			tokenizedText: []string{"uno uno dos dos dos tres dos"},
			expected:      []string{"uno", "dos", "tres"},
		}, {
			tokenizedText: []string{"uno", "dos", "uno", "uno"},
			expected:      []string{"uno", "dos"},
		}, {
			tokenizedText: []string{"unó", "dos", "uno", "uno"},
			expected:      []string{"unó", "uno", "dos"},
		},
	}
	for idx, tc := range testCases {
		result := getUniqueWordsForText(tc.tokenizedText)
		if err := testutils.CompareStringLists(result, tc.expected); err != nil {
			t.Errorf("Error on test case %d: %s", idx+1, err.Error())
		}
	}
}

func TestExtractWordExclusionsFromRankings(t *testing.T) {
	type testCase struct {
		input    []wordsmith.WordRanking
		expected []documents.WordExclusion
	}
	testCases := []testCase{
		// TODO: write some test cases
	}
	t.Errorf("Write some test cases")
	for idx, tc := range testCases {
		result := extractWordExclusionsFromRankings(tc.input)
		if err := compareWordExclusionLists(result, tc.expected); err != nil {
			t.Errorf("Error on test case %d: %s", idx+1, err.Error())
		}
	}
}

func compareWordExclusionLists(result, expected []documents.WordExclusion) error {
	if len(result) != len(expected) {
		return fmt.Errorf("lists not of equal size. Expected %d, but got %d", len(expected), len(result))
	}
	for idx, expectedExclusion := range expected {
		switch {
		case expectedExclusion.LeastFrequentRankingWithoutWord != result[idx].LeastFrequentRankingWithoutWord:
			return fmt.Errorf("Expected least frequent ranking on case %d is %d, but got %d", idx+1, expectedExclusion.LeastFrequentRankingWithoutWord, result[idx].LeastFrequentRankingWithoutWord)
		case expectedExclusion.WordText != result[idx].WordText:
			return fmt.Errorf("Expected least frequent ranking on case %d is %s, but got %s", idx+1, expectedExclusion.WordText, result[idx].WordText)
		}
	}
	return nil
}
