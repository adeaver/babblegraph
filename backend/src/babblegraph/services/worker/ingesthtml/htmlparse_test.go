package ingesthtml

import (
	"babblegraph/util/ptr"
	"babblegraph/util/testutils"
	"testing"
)

var (
	normalHTMLPage = `<html lang="es">
<head>
		<title>Page Title</title>
		<meta property="og:type" content="article" />
		<meta property="og:title" content="This is an Article" />
		<meta data-ue-u="og:title" content="Not an article" />
		<link rel="stylesheet" href="stylesheet.css" />
</head>
<body>
		<p>Some body text</p>
		<p>Text with a <a href="www.google.com">link</a></p>
		<p>Text with <strong>styling</strong> and a <span>span</span></p>
		<div>Text in a div</div>
		<script>
			some random javascript
		</script>
		<a href="/relative-link">relative link</a>
</body>`
	expectedParsedPage = ParsedHTMLPage{
		Links: []string{
			"www.google.com",
			"babblegraph.com/relative-link",
		},
		BodyText: "Some body text Text with a link Text with styling and a span Text in a div relative link",
		Language: ptr.String("es"),
		PageType: ptr.String("article"),
		Metadata: map[string]string{
			"og:type":  "article",
			"og:title": "This is an article",
		},
	}
)

func TestParseHTML(t *testing.T) {
	parsed, err := parseHTML("babblegraph.com", normalHTMLPage)
	if err != nil {
		t.Errorf("Not expecting error, but got one: %s", err.Error())
		return
	}
	if err := testutils.CompareStringLists(expectedParsedPage.Links, parsed.Links); err != nil {
		t.Errorf("Error on links: %s", err.Error())
	}
	if parsed.BodyText != expectedParsedPage.BodyText {
		t.Errorf("Error on body test. Expected %s, but got %s", expectedParsedPage.BodyText, parsed.BodyText)
	}
	if err := testutils.CompareNullableString(parsed.Language, expectedParsedPage.Language); err != nil {
		t.Errorf("Error on Language: %s", err.Error())
	}
	if err := testutils.CompareNullableString(parsed.PageType, expectedParsedPage.PageType); err != nil {
		t.Errorf("Error on page type: %s", err.Error())
	}
	if err := testutils.CompareStringMap(parsed.Metadata, expectedParsedPage.Metadata); err != nil {
		t.Errorf("Error on metadata: %s", err.Error())
	}
}
