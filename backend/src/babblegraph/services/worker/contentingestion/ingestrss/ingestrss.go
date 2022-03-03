package ingestrss

import (
	"babblegraph/util/ctx"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
)

func GetPodcastDataForRSSFeed(c ctx.LogContext, url string) error {
	feed, err := fetchRSSForURL(url)
	if err != nil {
		return err
	}
	c.Debugf("Got feed %+v", *feed)
	return nil
}

func fetchRSSForURL(u string) (*podcastRSSFeed, error) {
	resp, err := http.Get(u)
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
	var feed podcastRSSFeed
	if err := xml.Unmarshal(data, &feed); err != nil {
		return nil, err
	}
	return &feed, nil
}
