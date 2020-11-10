package ingesthtml

func ProcessURL(u, domain string) (*ParsedHTMLPage, error) {
	htmlStr, err := fetchHTMLForURL(u)
	if err != nil {
		return nil, err
	}
	return parseHTML(domain, *htmlStr)
}
