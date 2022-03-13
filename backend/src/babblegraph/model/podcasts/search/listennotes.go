package search

import (
	"babblegraph/util/cache"
	"babblegraph/util/ctx"
	"babblegraph/util/deref"
	"babblegraph/util/env"
	"encoding/json"
	"fmt"
	"strconv"
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
		for apiValue, displayName := range regions {
			displayNameStr, ok := displayName.(string)
			if !ok {
				c.Debugf("Display name %+v", displayNameStr)
				return nil, fmt.Errorf("Expected all api values to be string")
			}
			supportedRegions = append(supportedRegions, SupportedRegion{
				DisplayName: displayNameStr,
				APIValue:    apiValue,
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
			"top_level_only": "0",
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

type response struct {
	HasNext            bool        `json:"has_next"`
	HasPrevious        bool        `json:"has_previous"`
	ID                 interface{} `json:"id,string"`
	ListenNotesURL     string      `json:"listennotes_url"`
	Name               string      `json:"name"`
	NextPageNumber     int64       `json:"next_page_number"`
	PageNumber         int64       `json:"page_number"`
	ParentID           int64       `json:"parent_id"`
	PreviousPageNumber int64       `json:"previous_page_number"`
	Total              int64       `json:"total"`
	Podcasts           []podcast   `json:"podcasts"`
}

type podcast struct {
	ExternalID                  interface{}       `json:"id,string"`
	Country                     string            `json:"country"`
	Description                 string            `json:"description"`
	EarliestPubDateMilliseconds int64             `json:"earliest_pub_date_ms"`
	Email                       string            `json:"email"`
	ExplicitContent             bool              `json:"explicit_content"`
	Extra                       extraPodcastInfo  `json:"extra"`
	GenreIDs                    []int64           `json:"genre_ids"`
	ImageURL                    string            `json:"image"`
	IsClaimed                   bool              `json:"is_claimed"`
	ITunesID                    int64             `json:"itunes_id"`
	Language                    string            `json:"language"`
	LatestPubDateMilliseconds   int64             `json:"latest_pub_date_ms"`
	ListenScore                 interface{}       `json:"listen_score,string"`
	ListenScoreGlobalRank       interface{}       `json:"listen_score_global_rank,string"`
	ListenNotesURL              string            `json:"listennotes_url"`
	Publisher                   string            `json:"publisher"`
	RSS                         string            `json:"rss"`
	Thumbnail                   string            `json:"thumbnail"`
	Title                       string            `json:"title"`
	TotalNumberOfEpisodes       int64             `json:"total_episodes"`
	Type                        string            `json:"type"`
	Website                     string            `json:"website"`
	LookingFor                  podcastLookingFor `json:"looking_for"`
}

type extraPodcastInfo struct {
	FacebookHandle  string `json:"facebook_handle"`
	GoogleURL       string `json:"google_url"`
	InstagramHandle string `json:"instagram_handle"`
	LinkedInURL     string `json:"linkedin_url"`
	PatreonHandle   string `json:"patreon_handle"`
	SpotifyURL      string `json:"spotify_url"`
	TwitterHandle   string `json:"twitter_handle"`
	URL1            string `json:"url1"`
	URL2            string `json:"url2"`
	URL3            string `json:"url3"`
	WechatHandle    string `json:"wechat_handle"`
	YouTubeURL      string `json:"youtube_url"`
}

type podcastLookingFor struct {
	Cohosts        bool `json:"cohosts"`
	CrossPromotion bool `json:"cross_promotion"`
	Guests         bool `json:"guests"`
	Sponsors       bool `json:"sponsors"`
}

func getBestPodcastsForParams(c ctx.LogContext, params Params) (*response, error) {
	resp, err := client.FetchBestPodcasts(map[string]string{
		"page":      strconv.FormatInt(deref.Int64(params.PageNumber, 1), 10),
		"genre_id":  strconv.FormatInt(params.Genre, 10),
		"region":    params.Region,
		"language":  params.Language,
		"safe_mode": "1",
		"sort":      "listen_score",
	})
	if err != nil {
		return nil, err
	}
	var respJSON response
	if err := json.Unmarshal([]byte(resp.ToJSON()), &respJSON); err != nil {
		c.Debugf("Response %+v", resp.Data)
		return nil, err
	}
	return &respJSON, nil
}
