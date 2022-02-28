package search

import (
	"babblegraph/util/cache"
	"babblegraph/util/ctx"
	"babblegraph/util/env"
	"encoding/json"
	"fmt"
	"time"

	listennotes "github.com/ListenNotes/podcast-api-go"
)

const (
	supportedLanguagesTTL = 3 * 24 * time.Hour
	supportedRegionsTTL   = 3 * 24 * time.Hour
	supportedGenresTTL    = 3 * 24 * time.Hour
)

var client listennotes.HTTPClient

func init() {
	apiKey := env.GetEnvironmentVariableOrDefault("LISTEN_NOTES_API_KEY", "")
	client = listennotes.NewClient(apiKey)
}

func getSupportedLanguages() ([]string, error) {
	var supportedLanguages struct {
		Languages []string `json:"languages"`
	}
	if err := cache.WithCache("podcast_search_supported_languages", &supportedLanguages, supportedLanguagesTTL, func() (interface{}, error) {
		resp, err := client.FetchPodcastLanguages(nil)
		if err != nil {
			return nil, err
		}
		return resp.Data, nil
	}); err != nil {
		return nil, err
	}
	return supportedLanguages.Languages, nil
}

func getSupportedRegions(c ctx.LogContext) ([]SupportedRegion, error) {
	type supportedRegionsResponse struct {
		Regions []SupportedRegion `json:"regions"`
	}
	var supportedRegions supportedRegionsResponse
	if err := cache.WithCache("podcast_search_supported_regions", &supportedRegions, supportedRegionsTTL, func() (interface{}, error) {
		resp, err := client.FetchPodcastRegions(nil)
		if err != nil {
			return nil, err
		}
		regionsMap, ok := resp.Data["regions"]
		if !ok {
			c.Debugf("Regions response looked like: %+v", resp.Data)
			return nil, fmt.Errorf("Expected regions in response, but was not there")
		}
		regions, ok := regionsMap.(map[string]interface{})
		if !ok {
			c.Debugf("Got regions map: %+v", regionsMap)
			return nil, fmt.Errorf("Regions map did not cast into map")
		}
		var supportedRegions []SupportedRegion
		for displayName, apiValue := range regions {
			apiValueStr, ok := apiValue.(string)
			if !ok {
				c.Debugf("API Value %+v", apiValue)
				return nil, fmt.Errorf("Expected all api values to be string")
			}
			supportedRegions = append(supportedRegions, SupportedRegion{
				DisplayName: displayName,
				APIValue:    apiValueStr,
			})
		}
		return supportedRegionsResponse{
			Regions: supportedRegions,
		}, nil
	}); err != nil {
		return nil, err
	}
	return supportedRegions.Regions, nil
}

func getGenres(c ctx.LogContext) ([]SupportedGenre, error) {
	type supportedGenresResponse struct {
		Genres []SupportedGenre `json:"genres"`
	}
	var supportedGenres supportedGenresResponse
	if err := cache.WithCache("podcast_search_supported_genres", &supportedGenres, supportedGenresTTL, func() (interface{}, error) {
		resp, err := client.FetchPodcastGenres(map[string]string{
			"top_level_only": "1",
		})
		if err != nil {
			return nil, err
		}
		var respJSON struct {
			Genres []struct {
				ID       int    `json:"id"`
				Name     string `json:"name"`
				ParentID *int   `json:"parent_id"`
			} `json:"genres"`
		}
		if err := json.Unmarshal([]byte(resp.ToJSON()), &respJSON); err != nil {
			c.Debugf("Response %+v", resp.Data)
			return nil, err
		}
		var genres []SupportedGenre
		for _, genre := range respJSON.Genres {
			genres = append(genres, SupportedGenre{
				APIValue:    int64(genre.ID),
				DisplayName: genre.Name,
			})
		}
		return supportedGenresResponse{
			Genres: genres,
		}, nil
	}); err != nil {
		return nil, err
	}
	return supportedGenres.Genres, nil
}
