package web

type PageData struct {
	URLHost     string
	URLFileName string
	RawHTML     string
	BodyText    string
	Links       []string
}

func GetPageDataForURL(urlStr string) (*PageData, error) {
	host, fileName, err := splitURL(urlStr)
	if err != nil {
		return nil, err
	}
	htmlStr, err := getHTML(urlStr)
	if err != nil {
		return nil, err
	}
	bodyText, links, err := getBodyTextAndLinksForHTML(*htmlStr)
	if err != nil {
		return nil, err
	}
	return &PageData{
		URLHost:     *host,
		URLFileName: *fileName,
		RawHTML:     *htmlStr,
		BodyText:    *bodyText,
		Links:       links,
	}, nil
}
