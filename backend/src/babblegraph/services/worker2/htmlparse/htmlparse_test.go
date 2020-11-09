package htmlparse

import (
	"testing"
)

func TestParseNormalHTMLDocument(t *testing.T) {
	htmlDoc := `<html lang="es">
	<head>
		<title>Title</title>
		<meta name="viewport" content="width=device-width, initial-scale=1.0, minimum-scale=1.0">
		<meta content="Baltimore educators are tracking down students to fight low virtual attendance" property="og:title">
	</head>
	<body>
		<a href="#content">Local Link</a>
		<div>
			<div>
				<span>This is text</span>
			</div>
		</div>
		<p>This text is <strong>strong</strong></p>
		<a href="www.example.com">Real link</a>
	</body>`
	parsedDoc, err := parseHTMLDocument(htmlDoc)
	if err != nil {
		t.Errorf("Got error: %s", err.Error())
	}
	title, hasTitle := parsedDoc.Metadata["og:title"]
	if !hasTitle {
		t.Errorf("Expected to receive og:title key")
	}
	expectedTitle := "Baltimore educators are tracking down students to fight low virtual attendance"
	if title != expectedTitle {
		t.Errorf("Expected title property %s, but got %s", expectedTitle, title)
	}
	if len(parsedDoc.Links) != 1 {
		t.Errorf("Expected 1 link, but got %d", len(parsedDoc.Links))
	}
	expectedLink := "www.example.com"
	if parsedDoc.Links[0] != expectedLink {
		t.Errorf("Expected link %s, got %s", expectedLink, parsedDoc.Links[0])
	}
	if parsedDoc.LanguageValue == nil {
		t.Errorf("Expected a language value, did not receive one")
		return
	}
	if *parsedDoc.LanguageValue != "es" {
		t.Errorf("Expected a value of es, but got: %s", *parsedDoc.LanguageValue)
	}
}
