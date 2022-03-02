package ingesthtml

import (
	"babblegraph/model/content"
	"babblegraph/util/urlparser"
	"fmt"
	"strings"
)

// TODO: get rid of this
func ProcessURLDEPRECATED(u, domain string) (*ParsedHTMLPage, error) {
	if !strings.HasPrefix(u, "http://") && !strings.HasPrefix(u, "https://") {
		u = fmt.Sprintf("http://%s", u)
	}
	_, _, err := fetchHTMLForURL(u)
	if err != nil {
		return nil, err
	}
	return nil, nil
}

type ProcessURLInput struct {
	URL          string
	Source       content.Source
	SourceFilter *content.SourceFilter
}

func ProcessURL(input ProcessURLInput) (*ParsedHTMLPage, error) {
	urlWithProtocol, err := urlparser.EnsureProtocol(input.URL)
	if err != nil {
		return nil, err
	}
	htmlStr, cset, err := fetchHTMLForURL(*urlWithProtocol)
	if err != nil {
		return nil, err
	}
	return parseHTML(parseHTMLInput{
		htmlStr:      *htmlStr,
		cset:         *cset,
		source:       input.Source,
		sourceFilter: input.SourceFilter,
	})
}
