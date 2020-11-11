package linkprocessing

import (
	"babblegraph/util/urlparser"
	"testing"
)

func TestShouldFilterOutURL(t *testing.T) {
	type testCase struct {
		input          urlparser.ParsedURL
		expectedResult bool
	}

	testCases := []testCase{
		{
			input: urlparser.ParsedURL{
				Domain:        "google.com",
				URLIdentifier: "google.com",
				URL:           "google.com",
			},
			expectedResult: true,
		},
	}

	for idx, tc := range testCases {
		result := shouldFilterOutURL(tc.input)
		if result != tc.expectedResult {
			t.Errorf("Error on test case %d. Expected %t, but got %t", idx+1, tc.expectedResult, result)
		}
	}
}
