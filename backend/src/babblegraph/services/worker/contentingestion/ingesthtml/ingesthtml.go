package ingesthtml

import (
	"babblegraph/model/content"
	"babblegraph/util/urlparser"
)

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
