package text

import "testing"

func TestNormalize(t *testing.T) {
	testCases := []struct {
		input    string
		expected string
	}{
		{
			input:    "En España",
			expected: "en españa",
		},
	}
	for _, tc := range testCases {
		out := Normalize(tc.input)
		if out != tc.expected {
			t.Errorf("Expected %s, got %s", tc.expected, out)
		}
	}
}
