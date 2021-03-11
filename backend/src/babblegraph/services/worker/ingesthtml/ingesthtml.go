package ingesthtml

import (
	"fmt"
	"strings"
)

func ProcessURL(u, domain string) (*ParsedHTMLPage, error) {
	if !strings.HasPrefix(u, "http://") && !strings.HasPrefix(u, "https://") {
		u = fmt.Sprintf("http://%s", u)
	}
	htmlStr, cset, err := fetchHTMLForURL(u)
	if err != nil {
		return nil, err
	}
	return parseHTML(domain, *htmlStr, *cset)
}
