package search

import "babblegraph/util/ctx"

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
