package htmlfetch

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"babblegraph/util/storage"
)

func FetchAndStoreHTMLForURL(url string) (*storage.FileIdentifier, error) {
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
	return storage.WriteFile("html", string(data))
}
