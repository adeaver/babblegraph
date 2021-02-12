package text

import (
	"babblegraph/util/testutils"
	"testing"
)

func TestTokenize(t *testing.T) {
	type testCase struct {
		input    string
		expected []string
	}
	testCases := []testCase{
		{
			input:    "a great thing",
			expected: []string{"a", "great", "thing"},
		}, {
			input:    "a                    great thing",
			expected: []string{"a", "great", "thing"},
		}, {
			input:    "",
			expected: []string{},
		}, {
			input:    "                                       ",
			expected: []string{},
		},
	}
	for idx, tc := range testCases {
		result := Tokenize(tc.input)
		if err := testutils.CompareStringLists(tc.expected, result); err != nil {
			t.Errorf("Error on test case %d: %s", idx+1, err.Error())
		}
	}
}

func TestTokenizeUnique(t *testing.T) {
	type testCase struct {
		input    string
		expected []string
	}
	testCases := []testCase{
		{
			input:    "a great great thing",
			expected: []string{"a", "great", "thing"},
		}, {
			input:    "a great great great great thing",
			expected: []string{"a", "great", "thing"},
		}, {
			input:    "a great                great great              great thing thing             thing a      a",
			expected: []string{"a", "great", "thing"},
		}, {
			input:    "",
			expected: []string{},
		}, {
			input:    "                            ",
			expected: []string{},
		},
	}
	for idx, tc := range testCases {
		result := TokenizeUnique(tc.input)
		if err := testutils.CompareStringLists(tc.expected, result); err != nil {
			t.Errorf("Error on test case %d: %s", idx+1, err.Error())
		}
	}
}
