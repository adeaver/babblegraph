package ingesthtml

import (
	"fmt"
	"strings"

	"golang.org/x/net/html"
)

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

func isCurrentNodeTextNode(tagName string) bool {
	return len(tagName) > 0 && textTokenTagNames[string(tagName)]
}

func getLinksFromAnchor(node *html.Node, domain string) []string {
	var links []string
	for _, attr := range node.Attr {
		if attr.Key == "href" {
			u := attr.Val
			switch {
			case strings.HasPrefix(u, "/"):
				u = fmt.Sprintf("%s%s", domain, u)
			case strings.HasPrefix(u, "#"),
				strings.HasPrefix(u, "."),
				strings.Contains(u, ".jpeg"),
				strings.Contains(u, ".jpg"),
				strings.Contains(u, ".gif"),
				strings.Contains(u, ".png"):
				continue
			}
			links = append(links, u)
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
