package documents

import (
	"babblegraph/util/urlparser"
	"testing"
)

func TestMakeDocumentIndexForURL(t *testing.T) {
	type testCase struct {
		urlA      string
		urlB      string
		hasSameID bool
	}

	testCases := []testCase{
		{
			urlA:      "www.google.com/",
			urlB:      "www.google.com",
			hasSameID: true,
		}, {
			urlA:      "google.com/?q=value",
			urlB:      "www.google.com",
			hasSameID: true,
		}, {
			urlA:      "google.com/some-page",
			urlB:      "google.com",
			hasSameID: false,
		}, {
			urlA:      "mx.google.com",
			urlB:      "google.com",
			hasSameID: false,
		}, {
			urlA:      "babblegraph.com",
			urlB:      "google.com",
			hasSameID: false,
		}, {
			urlA:      "https://google.com/?q=value",
			urlB:      "www.google.com",
			hasSameID: true,
		}, {
			urlA:      "https://google.com/#fragment?q=value",
			urlB:      "https://google.com/?q=value",
			hasSameID: true,
		}, {
			urlA:      "https://google.com/#fragment?q=value",
			urlB:      "https://google.com/",
			hasSameID: true,
		}, {
			urlA:      "https://google.com/#fragment?q=value",
			urlB:      "https://google.com/?q=value#fragment",
			hasSameID: true,
		}, {
			urlA:      "https://google.com",
			urlB:      "https://google.com/?q=value#fragment",
			hasSameID: true,
		}, {
			urlA:      "https://google.com",
			urlB:      "https://google.com/#fragment",
			hasSameID: true,
		}, {
			urlA:      "https://google.com",
			urlB:      "google.com/#fragment",
			hasSameID: true,
		}, {
			urlA:      "https://google.com/page-name/page-part",
			urlB:      "google.com/page-name/page-part/#fragment",
			hasSameID: true,
		},
	}
	for idx, tc := range testCases {
		parsedURLA := urlparser.ParseURL(tc.urlA)
		if parsedURLA == nil {
			t.Errorf("Error on test case %d: URL A parses to null URL", idx+1)
			continue
		}
		parsedURLB := urlparser.ParseURL(tc.urlB)
		if parsedURLB == nil {
			t.Errorf("Error on test case %d: URL B parses to null URL", idx+1)
			continue
		}
		idA := makeDocumentIndexForURL(*parsedURLA)
		idB := makeDocumentIndexForURL(*parsedURLB)
		isDocIDTheSame := string(idA) == string(idB)
		if isDocIDTheSame != tc.hasSameID {
			t.Errorf("Error on test case %d: expected %t, but got %t", idx+1, tc.hasSameID, isDocIDTheSame)
		}
	}
}
