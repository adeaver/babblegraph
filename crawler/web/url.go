package web

import "net/url"

func splitURL(urlStr string) (_host, _fileName *string, _err error) {
	parsedURL, err := url.Parse(urlStr)
	if err != nil {
		return nil, nil, err
	}
	return &parsedURL.Host, &parsedURL.Path, nil
}
