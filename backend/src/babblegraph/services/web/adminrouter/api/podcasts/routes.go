package podcasts

import (
	"babblegraph/model/admin"
	"babblegraph/model/content"
	"babblegraph/model/podcasts"
	podcastsearch "babblegraph/model/podcasts/search"
	"babblegraph/services/web/adminrouter/middleware"
	"babblegraph/services/web/router"
	"babblegraph/util/database"
	"babblegraph/util/geo"
	"babblegraph/wordsmith"

	"github.com/jmoiron/sqlx"
)

var Routes = router.RouteGroup{
	Prefix: "podcasts",
	Routes: []router.Route{
		{
			Path: "get_podcast_search_options_1",
			Handler: middleware.WithPermission(
				admin.PermissionPodcastSearch,
				getPodcastSearchOptions,
			),
		}, {
			Path: "search_podcasts_1",
			Handler: middleware.WithPermission(
				admin.PermissionPodcastSearch,
				searchPodcasts,
			),
		},
	},
}

type getPodcastSearchOptionsRequest struct{}

type getPodcastSearchOptionsResponse struct {
	Options podcastsearch.Options `json:"options"`
}

func getPodcastSearchOptions(adminID admin.ID, r *router.Request) (interface{}, error) {
	options, err := podcastsearch.GetSearchOptions(r)
	if err != nil {
		return nil, err
	}
	return getPodcastSearchOptionsResponse{
		Options: *options,
	}, nil
}

type searchPodcastsRequest struct {
	Params podcastsearch.Params `json:"params"`
}

type searchPodcastsResponse struct {
	NextPageNumber *int64                          `json:"next_page_number,omitempty"`
	Podcasts       []podcastsearch.PodcastMetadata `json:"podcasts"`
}

func searchPodcasts(adminID admin.ID, r *router.Request) (interface{}, error) {
	var req searchPodcastsRequest
	if err := r.GetJSONBody(&req); err != nil {
		return nil, err
	}
	podcasts, nextPageNumber, err := podcastsearch.SearchPodcasts(r, req.Params)
	if err != nil {
		return nil, err
	}
	return searchPodcastsResponse{
		NextPageNumber: nextPageNumber,
		Podcasts:       podcasts,
	}, nil
}

type addPodcastRequest struct {
	CountryCode  string            `json:"country_code"`
	LanguageCode string            `json:"language_code"`
	TopicIDs     []content.TopicID `json:"topic_ids"`
	RSSFeedURL   string            `json:"rss_feed_url"`
	WebsiteURL   string            `json:"wesbite_url"`
	Title        string            `json:"title"`
}

type addPodcastResponse struct {
	Success bool `json:"success"`
}

func addPodcast(adminID admin.ID, r *router.Request) (interface{}, error) {
	var req addPodcastRequest
	if err := r.GetJSONBody(&req); err != nil {
		return nil, err
	}
	languageCode, err := wordsmith.GetLanguageCodeFromString(req.LanguageCode)
	if err != nil {
		return nil, err
	}
	countryCode, err := geo.GetCountryCodeFromString(req.CountryCode)
	if err != nil {
		return nil, err
	}
	if err := database.WithTx(func(tx *sqlx.Tx) error {
		return podcasts.AddPodcast(tx, podcasts.AddPodcastInput{
			CountryCode:  *countryCode,
			LanguageCode: *languageCode,
			WebsiteURL:   req.WebsiteURL,
			Title:        req.Title,
			TopicIDs:     req.TopicIDs,
			RSSFeedURL:   req.RSSFeedURL,
		})
	}); err != nil {
		return nil, err
	}
	return addPodcastResponse{
		Success: true,
	}, nil
}
