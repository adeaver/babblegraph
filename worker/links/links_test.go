package links

import "testing"

func TestGetDomain(t *testing.T) {
	testCases := []struct {
		input    string
		expected string
	}{
		{
			input:    "https://www.google.com/blah?q=123",
			expected: "google.com",
		}, {
			input:    "www.google.com/blah",
			expected: "google.com",
		}, {
			input:    "https://12345678:900",
			expected: "12345678",
		},
	}
	for i, tc := range testCases {
		result, _, err := getDomainAndCleanURL(tc.input)
		if err != nil {
			t.Errorf("Got error %s on test case %d", err.Error(), i+1)
		}
		if (*result).Str() != tc.expected {
			t.Errorf("Unexpected result on test case %d. Expected %s, but got %s", i+1, tc.expected, (*result).Str())
		}
	}
}

func TestGetCleanURL(t *testing.T) {
	testCases := []struct {
		input    string
		expected string
	}{
		{
			input:    "https://www.google.com/blah?q=123",
			expected: "https://www.google.com/blah",
		}, {
			input:    "www.google.com/blah",
			expected: "https://www.google.com/blah",
		}, {
			input:    "https://12345678:900",
			expected: "https://12345678",
		}, {
			input:    "google.com/blah",
			expected: "https://google.com/blah",
		},
	}
	for i, tc := range testCases {
		_, result, err := getDomainAndCleanURL(tc.input)
		if err != nil {
			t.Errorf("Got error %s on test case %d", err.Error(), i+1)
		}
		if *result != tc.expected {
			t.Errorf("Unexpected result on test case %d. Expected %s, but got %s", i+1, tc.expected, *result)
		}
	}
}
