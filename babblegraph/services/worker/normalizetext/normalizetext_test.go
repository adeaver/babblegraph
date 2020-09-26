package normalizetext

import "testing"

func TestNormalizeText(t *testing.T) {
	testCases := []struct {
		input    []byte
		expected string
	}{
		{
			input:    []byte("En España"),
			expected: "en españa",
		},
	}
	for _, tc := range testCases {
		out := normalizeText(tc.input)
		if out != tc.expected {
			t.Errorf("Expected %s, got %s", tc.expected, out[0])
		}
	}
}