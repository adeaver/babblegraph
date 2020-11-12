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
		}, {
			// super ambiguous url
			input: "blog.musica.ar/some-page/v2?q=123",
			// The reason that this parses not as a sudomain for musica.ar
			// is that there *should* be more websites ending in multiple tlds
			// rather than websites with a single tld as the domain name
			expectedDomain:        "blog.musica.ar",
			expectedURLIdentifier: "blog.musica.ar|some-page/v2",
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
		domain            string
		expectedDomain    *string
		expectedSubdomain *string
	}
	testCases := []testCase{
		{
			domain:            "www.google.com",
			expectedDomain:    ptr.String("google.com"),
			expectedSubdomain: ptr.String(""),
		}, {
			domain:            "www.blog.google.com",
			expectedDomain:    ptr.String("google.com"),
			expectedSubdomain: ptr.String("blog"),
		}, {
			domain:            "www",
			expectedDomain:    nil,
			expectedSubdomain: nil,
		}, {
			domain:            "com",
			expectedDomain:    nil,
			expectedSubdomain: nil,
		}, {
			domain:            "mx.google.com",
			expectedDomain:    ptr.String("google.com"),
			expectedSubdomain: ptr.String("mx"),
		}, {
			domain:            "musica.ar",
			expectedDomain:    ptr.String("musica.ar"),
			expectedSubdomain: ptr.String(""),
		}, {
			domain:            "www.musica.ar",
			expectedDomain:    ptr.String("musica.ar"),
			expectedSubdomain: ptr.String(""),
		}, {
			domain:            "google.com",
			expectedDomain:    ptr.String("google.com"),
			expectedSubdomain: ptr.String(""),
		}, {
			domain:            "google.cn",
			expectedDomain:    nil,
			expectedSubdomain: nil,
		}, {
			domain:            "www.google.cn",
			expectedDomain:    nil,
			expectedSubdomain: nil,
		}, {
			domain:            "www.google.co.uk",
			expectedDomain:    nil,
			expectedSubdomain: nil,
		},
	}
	for idx, tc := range testCases {
		resultDomain, resultSubdomain := verifyDomain(tc.domain)
		testNullableString(t, idx, "domain", resultDomain, tc.expectedDomain)
		testNullableString(t, idx, "subdomain", resultSubdomain, tc.expectedSubdomain)
	}
}

func TestIsValidURL(t *testing.T) {
	type testCase struct {
		input   string
		isValid bool
	}
	testCases := []testCase{
		{
			input:   "www.google.com",
			isValid: true,
		}, {
			input:   "/some-page-on-website",
			isValid: false,
		}, {
			input:   "https://some-website",
			isValid: false,
		}, {
			input:   "https://some-website.com",
			isValid: true,
		}, {
			input:   "?q=param",
			isValid: false,
		},
	}
	for idx, tc := range testCases {
		result := IsValidURL(tc.input)
		if result != tc.isValid {
			t.Errorf("Error on test case %d: expected %t, but got %t", idx+1, tc.isValid, result)
		}
	}
}

func testNullableString(t *testing.T, testCaseIdx int, fieldName string, result, expected *string) {
	switch {
	case result == nil && expected != nil:
		t.Errorf("Error on test case %d for field %s, expected %s, but got null", testCaseIdx+1, fieldName, *expected)
	case result != nil && expected == nil:
		t.Errorf("Error on test case %d for field %s, expected null, but got %s", testCaseIdx+1, fieldName, *result)
	case result == nil && expected == nil:
		// no-op
	case *result != *expected:
		t.Errorf("Error on test case %d for field %s, expected %s, but got %s", testCaseIdx+1, fieldName, *expected, *result)
	}
}
