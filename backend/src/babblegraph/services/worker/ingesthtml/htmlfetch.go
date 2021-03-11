package ingesthtml

import (
	"babblegraph/util/ptr"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

func fetchHTMLForURL(u string) (_body, _characterSet *string, _err error) {
	resp, err := http.Get(u)
	if err != nil {
		return nil, nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, nil, fmt.Errorf("Got status code for website: %d", resp.StatusCode)
	}
	cset := getCharacterSetForResponse(resp)
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, nil, err
	}
	return ptr.String(string(data)), ptr.String(cset), nil
}

func getCharacterSetForResponse(resp *http.Response) string {
	headers := resp.Header
	cset := "utf-8"
	if contentTypeHeaders, ok := headers["Content-Type"]; ok {
		joinedContentTypeHeaders := strings.Join(contentTypeHeaders, ";")
		if contentTypeHeaderParts := strings.Split(joinedContentTypeHeaders, "charset="); len(contentTypeHeaderParts) > 1 {
			cset = strings.Split(contentTypeHeaderParts[1], ";")[0]
		}
	}
	return cset
}
