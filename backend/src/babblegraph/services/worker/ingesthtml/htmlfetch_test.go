package ingesthtml

import (
	"net/http"
	"testing"
)

func TestGetCharacterSet(t *testing.T) {
	type testCase struct {
		contentTypeHeaders   map[string][]string
		expectedCharacterSet string
	}
	testCases := []testCase{
		{
			contentTypeHeaders: map[string][]string{
				"text/html; charset=latin1",
			},
			expectedCharacterSet: "latin1",
		}, {
			contentTypeHeaders: map[string][]string{
				"something",
				"text/html; charset=latin1",
			},
			expectedCharacterSet: "latin1",
		}, {
			contentTypeHeaders:   map[string][]string{},
			expectedCharacterSet: "utf-8",
		},
	}
	for idx, tc := range testCases {
		result := getCharacterSetForResponse(&http.Response{
			Header: tc.contentTypeHeaders,
		})
		if result != tc.expectedCharacterSet {
			t.Errorf("Error on test case %d: expected %s, but got %s", idx+1, tc.expectedCharacterSet, result)
		}
	}
}
