package ingesthtml

import (
	"babblegraph/model/content"
	"babblegraph/util/deref"
	"babblegraph/util/ptr"
	"fmt"
	"io"
	"log"
	"strings"

	"golang.org/x/net/html"
	"golang.org/x/text/encoding/htmlindex"
)

type ParsedHTMLPage struct {
	Links       []string
	BodyText    string
	Language    *string
	PageType    *string
	Metadata    map[string]string
	IsPaywalled bool
}

type parseHTMLInput struct {
	htmlStr      string
	cset         string
	source       content.Source
	sourceFilter *content.SourceFilter
}

func parseHTML(input parseHTMLInput) (*ParsedHTMLPage, error) {
	var body io.Reader = strings.NewReader(input.htmlStr)
	e, err := htmlindex.Get(input.cset)
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
	var isParseInTextNodeType, isParseLDJSON, isPaywalled bool
	var bodyText, links []string
	var language *string
	metadata := make(map[string]string)
	var f func(node *html.Node)
	f = func(node *html.Node) {
		switch node.Type {
		case html.ElementNode:
			if input.sourceFilter != nil {
				switch {
				case deref.Bool(input.sourceFilter.UseLDJSONValidation, false):
					if node.Data == "script" {
						for _, attr := range node.Attr {
							if attr.Key == "type" && attr.Val == "application/ld+json" {
								isParseLDJSON = true
							}
						}
					}
				case len(input.sourceFilter.PaywallClasses) != 0:
					isPaywalled = isPaywalled || processPaywallFromClasses(node, input.sourceFilter.PaywallClasses)
				case len(input.sourceFilter.PaywallIDs) != 0:
					isPaywalled = isPaywalled || processPaywallFromIDs(node, input.sourceFilter.PaywallIDs)
				default:
					log.Println("Paywall validation is not null, but no paywall validation type is specified")
				}
			}
			// The way that this works is that we encounter
			// the marker of a text node (i.e. the p tag)
			// prior to encountering the text for that node
			// if we encounter a text node, then we need to
			// collect the text in the next node
			isParseInTextNodeType = isCurrentNodeTextNode(node.Data)
			switch node.Data {
			case "a":
				links = append(links, getLinksFromAnchor(node, input.source.URL)...)
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
			switch {
			case isParseInTextNodeType:
				bodyText = append(bodyText, node.Data)
			case isParseLDJSON:
				ldJSONPaywall, err := processPaywallFromLDJSON(node.Data)
				if err != nil {
					log.Println(fmt.Sprintf("Error unmarshalling ld+json: %s", err.Error()))
				}
				isPaywalled = isPaywalled || ldJSONPaywall
				isParseLDJSON = false
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
		Links:       links,
		BodyText:    bodyTextStr,
		Language:    language,
		PageType:    pageType,
		Metadata:    metadata,
		IsPaywalled: isPaywalled,
	}, nil
}
