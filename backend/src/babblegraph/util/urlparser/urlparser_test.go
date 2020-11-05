package urlparser

import (
	"babblegraph/util/ptr"
	"testing"
)

func TestParseURL(t *testing.T) {
	type testCase struct {
		input                 string
		expectedDomain        string
		expectedURLIdentifier string
	}
	testCases := []testCase{
		{
			// normal https
			input:                 "https://google.com",
			expectedDomain:        "google.com",
			expectedURLIdentifier: "google.com",
		}, {
			// normal http
			input:                 "http://google.com",
			expectedDomain:        "google.com",
			expectedURLIdentifier: "google.com",
		}, {
			// http with slash
			input:                 "http://google.com/",
			expectedDomain:        "google.com",
			expectedURLIdentifier: "google.com",
		}, {
			// https with query params
			input:                 "http://google.com/?q=value",
			expectedDomain:        "google.com",
			expectedURLIdentifier: "google.com",
		}, {
			// subdomain with query params
			input:                 "http://blog.google.com/?q=value",
			expectedDomain:        "google.com",
			expectedURLIdentifier: "blog|google.com",
		}, {
			// www identifier
			input:                 "http://www.google.com/?q=value",
			expectedDomain:        "google.com",
			expectedURLIdentifier: "google.com",
		}, {
			// no https
			input:                 "www.google.com/?q=value",
			expectedDomain:        "google.com",
			expectedURLIdentifier: "google.com",
		}, {
			// no http, www, with page
			input:                 "www.google.com/some-page",
			expectedDomain:        "google.com",
			expectedURLIdentifier: "google.com|some-page",
		}, {
			// no http, www, with multilayer page
			input:                 "www.google.com/some-page/v2",
			expectedDomain:        "google.com",
			expectedURLIdentifier: "google.com|some-page/v2",
		}, {
			// no http, www, with multilayer page and query params
			input:                 "www.google.com/some-page/v2?q=123",
			expectedDomain:        "google.com",
			expectedURLIdentifier: "google.com|some-page/v2",
		}, {
			// ambiguous url
			input:                 "www.musica.ar/some-page/v2?q=123",
			expectedDomain:        "musica.ar",
			expectedURLIdentifier: "musica.ar|some-page/v2",
		},
	}
	for idx, tc := range testCases {
		result := ParseURL(tc.input)
		if result == nil {
			t.Errorf("Error on test case %d: got null result", idx+1)
			continue
		}
		if result.Domain != tc.expectedDomain {
			t.Errorf("Error on test case %d: expected domain %s, but got %s", idx+1, tc.expectedDomain, result.Domain)
		}
		if result.URLIdentifier != tc.expectedURLIdentifier {
			t.Errorf("Error on test case %d: expected identifier %s, but got %s", idx+1, tc.expectedURLIdentifier, result.URLIdentifier)
		}
	}
}

func TestFindURLParts(t *testing.T) {
	type testCase struct {
		rawURL          string
		expectedWebsite string
		expectedPage    string
		expectedParams  string
	}
	testCases := []testCase{
		{
			rawURL:          "http://www.google.com/",
			expectedWebsite: "www.google.com",
			expectedPage:    "",
			expectedParams:  "",
		}, {
			rawURL:          "http://www.google.com/some-page",
			expectedWebsite: "www.google.com",
			expectedPage:    "some-page",
			expectedParams:  "",
		}, {
			rawURL:          "http://google.com/some-page",
			expectedWebsite: "google.com",
			expectedPage:    "some-page",
			expectedParams:  "",
		}, {
			rawURL:          "www.google.com/some-page",
			expectedWebsite: "www.google.com",
			expectedPage:    "some-page",
			expectedParams:  "",
		}, {
			rawURL:          "www.google.com/some-page.html",
			expectedWebsite: "www.google.com",
			expectedPage:    "some-page.html",
			expectedParams:  "",
		},
	}
	for idx, tc := range testCases {
		result := findURLParts(tc.rawURL)
		switch {
		case result.Website != tc.expectedWebsite:
			t.Errorf("Error on test case %d: expected website %s, got %s", idx+1, tc.expectedWebsite, result.Website)
		case result.Page != tc.expectedPage:
			t.Errorf("Error on test case %d: expected page %s, got %s", idx+1, tc.expectedPage, result.Page)
		case result.Params != tc.expectedParams:
			t.Errorf("Error on test case %d: expected params %s, got %s", idx+1, tc.expectedParams, result.Params)
		}
	}
}

func TestVerifyDomain(t *testing.T) {
	type testCase struct {
		domain         string
		verifiedOutput *string
	}
	testCases := []testCase{
		{
			domain:         "www.google.com",
			verifiedOutput: ptr.String("google.com"),
		}, {
			domain:         "www.blog.google.com",
			verifiedOutput: ptr.String("blog.google.com"),
		}, {
			domain:         "www",
			verifiedOutput: nil,
		}, {
			domain:         "com",
			verifiedOutput: nil,
		}, {
			domain:         "mx.google.com",
			verifiedOutput: ptr.String("mx.google.com"),
		}, {
			domain:         "musica.ar",
			verifiedOutput: ptr.String("musica.ar"),
		}, {
			domain:         "www.musica.ar",
			verifiedOutput: ptr.String("musica.ar"),
		}, {
			domain:         "google.com",
			verifiedOutput: ptr.String("google.com"),
		}, {
			domain:         "google.cn",
			verifiedOutput: nil,
		}, {
			domain:         "www.google.cn",
			verifiedOutput: nil,
		}, {
			domain:         "www.google.co.uk",
			verifiedOutput: nil,
		},
	}
	for idx, tc := range testCases {
		result := verifyDomain(tc.domain)
		switch {
		case result == nil && tc.verifiedOutput != nil:
			t.Errorf("Error on test case %d, expected %s, but got null", idx+1, *tc.verifiedOutput)
		case result != nil && tc.verifiedOutput == nil:
			t.Errorf("Error on test case %d, expected null, but got %s", idx+1, *result)
		case result == nil && tc.verifiedOutput == nil:
			// no-op
		case *result != *tc.verifiedOutput:
			t.Errorf("Error on test case %d, expected %s, but got %s", idx+1, *tc.verifiedOutput, *result)
		}
	}
}
