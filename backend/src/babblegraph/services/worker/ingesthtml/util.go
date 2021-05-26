package ingesthtml

import (
	"encoding/json"
	"fmt"
	"log"
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

func processPaywallFromLDJSON(ldJSONData string) (bool, error) {
	var ldJSON map[string]interface{}
	if err := json.Unmarshal([]byte(ldJSONData), &ldJSON); err != nil {
		return false, err
	}
	isAccessibleInterface, ok := ldJSON["isAccessibleForFree"]
	if !ok {
		log.Println("LD+JSON does not contain isAccessibleForFree, assuming not paywalled")
		return false, nil
	}
	isAccessibleForFree, ok := isAccessibleInterface.(bool)
	if !ok {
		log.Println("Could not convert isAccessibleForFree key to bool, trying string...")
		isAccessibleForFreeStr, ok := isAccessibleInterface.(string)
		if !ok {
			return false, fmt.Errorf("Could not convert isAccessibleForFree to bool or string")
		}
		isAccessibleForFree = strings.ToLower(isAccessibleForFreeStr) != "false"
	}
	return !isAccessibleForFree, nil
}

func processPaywallFromClasses(node *html.Node, paywallValidationClasses []string) bool {
	return processPaywallFromAttr(node, "class", paywallValidationClasses)
}

func processPaywallFromIDs(node *html.Node, paywallValidationIDs []string) bool {
	return processPaywallFromAttr(node, "id", paywallValidationIDs)
}

func processPaywallFromAttr(node *html.Node, attrName string, checkValues []string) bool {
	for _, attr := range node.Attr {
		if attr.Key == attrName {
			values := strings.Split(attr.Val, " ")
			for _, v := range values {
				for _, paywallAttrVal := range checkValues {
					if v == paywallAttrVal {
						return true
					}
				}
			}
		}
	}
	return false
}
