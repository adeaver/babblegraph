package ingesthtml

import (
	"fmt"
	"strings"
)

func ProcessURL(u, domain string) (*ParsedHTMLPage, error) {
	if !strings.HasPrefix(u, "http://") && !strings.HasPrefix(u, "https://") {
		u = fmt.Sprintf("https://%s", u)
	}
	htmlStr, err := fetchHTMLForURL(u)
	if err != nil {
		return nil, err
	}
	return parseHTML(domain, *htmlStr)
}
