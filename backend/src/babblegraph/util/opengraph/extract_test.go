package opengraph

import (
	"babblegraph/util/ptr"
	"babblegraph/util/testutils"
	"testing"
)

func compareBasicMetadata(t *testing.T, testCaseIdx int, a, b BasicMetadata) {
	if err := testutils.CompareNullableString(a.Title, b.Title); err != nil {
		t.Errorf("Error comparing titles on test case %d: %s", testCaseIdx+1, err.Error())
	}
	if err := testutils.CompareNullableString(a.URL, b.URL); err != nil {
		t.Errorf("Error comparing urls on test case %d: %s", testCaseIdx+1, err.Error())
	}
	if err := testutils.CompareNullableString(a.ImageURL, b.ImageURL); err != nil {
		t.Errorf("Error comparing image urls on test case %d: %s", testCaseIdx+1, err.Error())
	}
	if err := testutils.CompareNullableString(a.Type, b.Type); err != nil {
		t.Errorf("Error comparing types on test case %d: %s", testCaseIdx+1, err.Error())
	}
}

func TestGetBasicMetadata(t *testing.T) {
	type testCase struct {
		input    map[string]string
		expected BasicMetadata
	}

	testCases := []testCase{
		{
			input: map[string]string{
				"og:type":  "article",
				"og:image": "www.babblegraph.com/image.jpeg",
				"og:title": "title",
				"og:url":   "www.babblegraph.com",
			},
			expected: BasicMetadata{
				Type:     ptr.String("article"),
				ImageURL: ptr.String("www.babblegraph.com/image.jpeg"),
				Title:    ptr.String("title"),
				URL:      ptr.String("www.babblegraph.com"),
			},
		}, {
			input: map[string]string{
				"og:image": "www.babblegraph.com/image.jpeg",
				"og:title": "title",
				"og:url":   "www.babblegraph.com",
			},
			expected: BasicMetadata{
				Type:     nil,
				ImageURL: ptr.String("www.babblegraph.com/image.jpeg"),
				Title:    ptr.String("title"),
				URL:      ptr.String("www.babblegraph.com"),
			},
		}, {
			input: map[string]string{
				"og:type":  "article",
				"og:image": "www.babblegraph.com/image.jpeg",
				"og:url":   "www.babblegraph.com",
			},
			expected: BasicMetadata{
				Type:     ptr.String("article"),
				ImageURL: ptr.String("www.babblegraph.com/image.jpeg"),
				Title:    nil,
				URL:      ptr.String("www.babblegraph.com"),
			},
		}, {
			input: map[string]string{
				"og:type":  "article",
				"og:image": "www.babblegraph.com/image.jpeg",
				"og:title": "title",
			},
			expected: BasicMetadata{
				Type:     ptr.String("article"),
				ImageURL: ptr.String("www.babblegraph.com/image.jpeg"),
				Title:    ptr.String("title"),
				URL:      nil,
			},
		}, {
			input: map[string]string{
				"og:type":  "article",
				"og:title": "title",
				"og:url":   "www.babblegraph.com",
			},
			expected: BasicMetadata{
				Type:     ptr.String("article"),
				ImageURL: nil,
				Title:    ptr.String("title"),
				URL:      ptr.String("www.babblegraph.com"),
			},
		}, {
			input: map[string]string{
				"og:type":  "article",
				"og:title": "title",
				"charset":  "utf-8",
				"og:url":   "www.babblegraph.com",
			},
			expected: BasicMetadata{
				Type:     ptr.String("article"),
				ImageURL: nil,
				Title:    ptr.String("title"),
				URL:      ptr.String("www.babblegraph.com"),
			},
		},
	}
	for idx, tc := range testCases {
		result := GetBasicMetadata(tc.input)
		compareBasicMetadata(t, idx, result, tc.expected)
	}
}
