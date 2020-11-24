package opengraph

import (
	"babblegraph/util/ptr"
	"babblegraph/util/testutils"
	"testing"
	"time"
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
	if err := testutils.CompareNullableString(a.Description, b.Description); err != nil {
		t.Errorf("Error comparing descriptions on test case %d: %s", testCaseIdx+1, err.Error())
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
		}, {
			input: map[string]string{
				"og:type":        "article",
				"og:description": "A test",
				"charset":        "utf-8",
				"og:url":         "www.babblegraph.com",
			},
			expected: BasicMetadata{
				Type:        ptr.String("article"),
				ImageURL:    nil,
				Title:       nil,
				URL:         ptr.String("www.babblegraph.com"),
				Description: ptr.String("A test"),
			},
		},
	}
	for idx, tc := range testCases {
		result := GetBasicMetadata(tc.input)
		compareBasicMetadata(t, idx, result, tc.expected)
	}
}

func TestGetPublicationTime(t *testing.T) {
	type testCase struct {
		input    map[string]string
		expected *time.Time
	}
	testCases := []testCase{
		{
			input: map[string]string{
				"charset": "utf-8",
			},
			expected: nil,
		}, {
			input: map[string]string{
				"charset":                   "utf-8",
				"og:article:published_time": "2020-11-24T15:30:06+00:00",
			},
			expected: ptr.Time(time.Date(2020, time.November, 24, 15, 30, 6, 0, time.UTC)),
		}, {
			input: map[string]string{
				"charset":                   "utf-8",
				"og:article:published_time": "2020-11-24 15:30:06+00:00",
			},
			expected: nil,
		},
	}
	for idx, tc := range testCases {
		publicationTime := lookupPublicationTime(tc.input)
		switch {
		case tc.expected == nil && publicationTime == nil:
			// no-op
		case tc.expected == nil && publicationTime != nil:
			t.Errorf("Error on test case %d: Expected null publication time, but got %v", idx+1, *publicationTime)
		case tc.expected != nil && publicationTime == nil:
			t.Errorf("Error on test case %d: Expected publication time of %v, but got null", idx+1, *tc.expected)
		case !(tc.expected.Equal(*publicationTime)):
			t.Errorf("Error on test case %d: Expected publication time of %v, but got %v", idx+1, *tc.expected, *publicationTime)
		default:
			// no-op

		}
	}
}
