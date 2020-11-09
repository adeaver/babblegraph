package htmlparse

import (
	"fmt"
	"log"
	"strings"

	"babblegraph/util/storage"

	"golang.org/x/net/html"
)

type ParsedDocumentInStorage struct {
	BodyTextFilename storage.FileIdentifier
	Links            []string
	Metadata         map[string]string
	LanguageValue    *string
}

func ParseAndStoreFileText(filename storage.FileIdentifier) (*ParsedDocumentInStorage, error) {
	htmlBytes, err := storage.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	parsedDocument, err := parseHTMLDocument(string(htmlBytes))
	if err != nil {
		return nil, err
	}
	id, err := storage.WriteFile("txt", parsedDocument.BodyText)
	if err != nil {
		return nil, err
	}
	return &ParsedDocumentInStorage{
		BodyTextFilename: *id,
		Links:            parsedDocument.Links,
		Metadata:         parsedDocument.Metadata,
		LanguageValue:    parsedDocument.LanguageValue,
	}, nil
}

var textTokenTagNames = map[string]bool{
	"a":       true,
	"abbr":    true,
	"address": true,
	"b":       true,
	"center":  true,
	"h1":      true,
	"h2":      true,
	"h3":      true,
	"h4":      true,
	"h5":      true,
	"h6":      true,
	"li":      true,
	"p":       true,
	"span":    true,
	"strong":  true,
	"td":      true,
}

func isTagNameForTextToken(tagName []byte) bool {
	return len(tagName) > 0 && textTokenTagNames[string(tagName)]
}

func isWeblink(href string) bool {
	switch {
	case strings.HasPrefix(href, "#"),
		strings.HasPrefix(href, "/"),
		strings.HasPrefix(href, "."):
		return false
	case strings.Contains(href, ".jpeg"),
		strings.Contains(href, ".jpg"),
		strings.Contains(href, ".gif"),
		strings.Contains(href, ".png"):
		return false
	}
	return true
}

type parsedDocument struct {
	BodyText      string
	Links         []string
	Metadata      map[string]string
	LanguageValue *string
}

func parseHTMLDocument(htmlStr string) (*parsedDocument, error) {
	htmlReader := strings.NewReader(htmlStr)
	htmlDoc, err := html.Parse(htmlReader)
	if err != nil {
		return nil, err
	}
	var language *string
	metadata := make(map[string]string)
	var links, bodyText []string
	var shouldCollectText bool
	var f func(n *html.Node)
	f = func(n *html.Node) {
		switch n.Type {
		case html.ElementNode:
			_, shouldCollectText = textTokenTagNames[n.Data]
			switch n.Data {
			case "a":
				links = append(links, getLinksFromAnchor(n)...)
			case "meta":
				if name, value := getKeyValuePairFromMetaTag(n); name != nil && value != nil {
					metadata[*name] = *value
				}
			case "html":
				for _, attr := range n.Attr {
					if attr.Key == "lang" {
						language = &attr.Val
					}
				}
			}
		case html.TextNode:
			if shouldCollectText {
				bodyText = append(bodyText, n.Data)
			}
		case html.ErrorNode:
			log.Println(fmt.Sprintf("Error: %s", n.Data))
		case html.DocumentNode, html.CommentNode, html.DoctypeNode:
			// no-op
		default:
			log.Fatal(fmt.Sprintf("Unrecognized node type: %d", n.Type))
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(htmlDoc)
	body := strings.Join(bodyText, "\n")
	return &parsedDocument{
		BodyText:      body,
		Links:         links,
		Metadata:      metadata,
		LanguageValue: language,
	}, nil
}

func getLinksFromAnchor(n *html.Node) []string {
	var links []string
	for _, attr := range n.Attr {
		if attr.Key == "href" && isWeblink(attr.Val) {
			links = append(links, attr.Val)
		}
	}
	return links
}

func getKeyValuePairFromMetaTag(n *html.Node) (_key, _value *string) {
	var name, value string
	var foundName, foundVal bool
	for _, attr := range n.Attr {
		if attr.Key == "name" || attr.Key == "property" {
			name = attr.Val
			foundName = true
		}
		if attr.Key == "content" {
			value = attr.Val
			foundVal = true
		}
	}
	if !foundName || !foundVal {
		return nil, nil
	}
	return &name, &value
}
