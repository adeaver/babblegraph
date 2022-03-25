package search

import (
	"babblegraph/util/ctx"
	"babblegraph/util/ptr"
	"fmt"
)

type Options struct {
	SupportedLanguages []string          `json:"supported_languages"`
	SupportedRegions   []SupportedRegion `json:"supported_regions"`
	Genres             []SupportedGenre  `json:"genres"`
}

type SupportedRegion struct {
	DisplayName string `json:"display_name"`
	APIValue    string `json:"api_value"`
}

type SupportedGenre struct {
	DisplayName string `json:"display_name"`
	APIValue    int64  `json:"api_value"`
}

func GetSearchOptions(c ctx.LogContext) (*Options, error) {
	supportedLanguages, err := getSupportedLanguages()
	if err != nil {
		return nil, err
	}
	supportedRegions, err := getSupportedRegions(c)
	if err != nil {
		return nil, err
	}
	genres, err := getGenres(c)
	if err != nil {
		return nil, err
	}
	return &Options{
		SupportedLanguages: supportedLanguages,
		SupportedRegions:   supportedRegions,
		Genres:             genres,
	}, nil
}

type Params struct {
	Language   string `json:"language"`
	Region     string `json:"region"`
	Genre      int64  `json:"genre"`
	PageNumber *int64 `json:"page_number,omitempty"`
}

type PodcastMetadata struct {
	ExternalID            string `json:"external_id"`
	Title                 string `json:"title"`
	Country               string `json:"country"`
	Description           string `json:"description"`
	Website               string `json:"website"`
	Language              string `json:"language"`
	Type                  string `json:"type"`
	TotalNumberOfEpisodes int64  `json:"total_number_of_episodes"`
	ListenNotesURL        string `json:"listen_notes_url"`
}

func SearchPodcasts(c ctx.LogContext, params Params) (_results []PodcastMetadata, _nextPageNumber *int64, _err error) {
	resp, err := getBestPodcastsForParams(c, params)
	if err != nil {
		return nil, nil, err
	}
	var nextPage *int64
	if resp.HasNext {
		nextPage = ptr.Int64(resp.NextPageNumber)
	}
	var out []PodcastMetadata
	for _, p := range resp.Podcasts {
		parsedWebsiteURL := MaybeParseURLForListenNotesWebsiteURL(p.Website)
		if parsedWebsiteURL == nil {
			c.Infof("Skipping url %s", p.Website)
			continue
		}
		out = append(out, PodcastMetadata{
			ExternalID:            fmt.Sprintf("%v", p.ExternalID),
			Title:                 p.Title,
			Country:               p.Country,
			Description:           p.Description,
			Website:               parsedWebsiteURL.URL,
			Language:              p.Language,
			Type:                  p.Type,
			TotalNumberOfEpisodes: p.TotalNumberOfEpisodes,
			ListenNotesURL:        p.ListenNotesURL,
		})
	}
	return out, nextPage, nil
}
