package ingesthtml

type ParsedHTMLPage struct {
	Links    []string
	BodyText string
	Language *string
	PageType *string
	Metadata map[string]string
}

func parseHTML(domain string, htmlBody string) (*ParsedHTMLPage, error) {
	return nil, nil
}
