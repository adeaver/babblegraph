package htmlparse

import (
	"babblegraph/worker/storage"
	"fmt"
	"log"
	"strings"

	"golang.org/x/net/html"
)

func ParseAndStoreFileText(filename storage.FileIdentifier) (*storage.FileIdentifier, []string, error) {
	htmlBytes, err := storage.ReadFile(filename)
	if err != nil {
		return nil, nil, err
	}
	text, links, err := getTextAndLinksForHTML(string(htmlBytes))
	if err != nil {
		return nil, nil, err
	}
	id, err := storage.WriteFile("txt", *text)
	if err != nil {
		return nil, nil, err
	}
	return id, links, nil
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

func getTextAndLinksForHTML(htmlStr string) (*string, []string, error) {
	htmlReader := strings.NewReader(htmlStr)
	htmlDoc, err := html.Parse(htmlReader)
	if err != nil {
		return nil, nil, err
	}
	var links []string
	var bodyText []string
	var shouldCollectText bool
	var f func(n *html.Node)
	f = func(n *html.Node) {
		switch n.Type {
		case html.ElementNode:
			_, shouldCollectText = textTokenTagNames[n.Data]
			if n.Data == "a" {
				for _, attr := range n.Attr {
					if attr.Key == "href" && isWeblink(attr.Val) {
						links = append(links, attr.Val)
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
			log.Fatal("Unrecognized node type: %d", n.Type)
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(htmlDoc)
	body := strings.Join(bodyText, "\n")
	return &body, links, nil
}
