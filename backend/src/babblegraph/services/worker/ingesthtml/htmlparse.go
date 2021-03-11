package ingesthtml

import (
	"babblegraph/util/ptr"
	"fmt"
	"io"
	"log"
	"strings"

	"golang.org/x/net/html"
	"golang.org/x/text/encoding/htmlindex"
)

type ParsedHTMLPage struct {
	Links    []string
	BodyText string
	Language *string
	PageType *string
	Metadata map[string]string
}

func parseHTML(domain, htmlStr, cset string) (*ParsedHTMLPage, error) {
	var body io.Reader = strings.NewReader(htmlStr)
	e, err := htmlindex.Get(cset)
	if err != nil {
		return nil, err
	}
	// Ignoring the error here since HTML pages
	// could potentially have garbage charsets
	if name, _ := htmlindex.Name(e); name != "utf-8" {
		body = e.NewDecoder().Reader(body)
	}
	htmlDoc, err := html.Parse(body)
	if err != nil {
		return nil, err
	}
	var isParseInTextNodeType bool
	var bodyText, links []string
	var language *string
	metadata := make(map[string]string)
	var f func(node *html.Node)
	f = func(node *html.Node) {
		switch node.Type {
		case html.ElementNode:
			// The way that this works is that we encounter
			// the marker of a text node (i.e. the p tag)
			// prior to encountering the text for that node
			// if we encounter a text node, then we need to
			// collect the text in the next node
			isParseInTextNodeType = isCurrentNodeTextNode(node.Data)
			switch node.Data {
			case "a":
				links = append(links, getLinksFromAnchor(node, domain)...)
			case "meta":
				if name, value := getKeyValuePairFromMetaTag(node); name != nil && value != nil {
					metadata[*name] = *value
				}
			case "html":
				for _, attr := range node.Attr {
					if attr.Key == "lang" {
						language = ptr.String(attr.Val)
					}
				}
			}
		case html.TextNode:
			if isParseInTextNodeType {
				bodyText = append(bodyText, node.Data)
			}
		case html.ErrorNode:
			log.Println(fmt.Sprintf("Error: %s", node.Data))
		case html.DocumentNode, html.CommentNode, html.DoctypeNode:
			// no-op
		default:
			log.Fatal(fmt.Sprintf("Unrecognized node type: %d", node.Type))
		}
		for c := node.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(htmlDoc)
	bodyTextStr := strings.Join(bodyText, "\n")
	var pageType *string
	if ogType, ok := metadata["og:type"]; ok {
		pageType = ptr.String(ogType)
	}
	return &ParsedHTMLPage{
		Links:    links,
		BodyText: bodyTextStr,
		Language: language,
		PageType: pageType,
		Metadata: metadata,
	}, nil
}
