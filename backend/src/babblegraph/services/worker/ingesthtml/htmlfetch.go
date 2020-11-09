package ingesthtml

import (
	"babblegraph/util/ptr"
	"fmt"
	"io/ioutil"
	"net/http"
)

func fetchHTMLForURL(u string) (*string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Got status code for website: %d", resp.StatusCode)
	}
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return ptr.String(string(data)), nil
}
