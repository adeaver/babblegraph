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
	"babblegraph/util/ptr"
	"babblegraph/util/urlparser"
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
		}, {
			Path: "add_podcast_1",
			Handler: middleware.WithPermission(
				admin.PermissionEditContentSources,
				addPodcast,
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
	Podcasts       []podcastMetadataWithSourceInfo `json:"podcasts"`
}

type podcastMetadataWithSourceInfo struct {
	Metadata podcastsearch.PodcastMetadata `json:"metadata"`
	SourceID *content.SourceID             `json:"source_id,omitempty"`
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
	var out []podcastMetadataWithSourceInfo
	if err := database.WithTx(func(tx *sqlx.Tx) error {
		for _, p := range podcasts {
			r.Infof("Website %s", p.Website)
			parsedURL := urlparser.ParseURL(p.Website)
			if parsedURL == nil {
				continue
			}
			var sourceID *content.SourceID
			sourceID, err := content.LookupSourceIDForParsedURL(tx, *parsedURL)
			if err != nil {
				return err
			}
			out = append(out, podcastMetadataWithSourceInfo{
				Metadata: p,
				SourceID: sourceID,
			})
		}
		return nil
	}); err != nil {
		return nil, err
	}
	return searchPodcastsResponse{
		NextPageNumber: nextPageNumber,
		Podcasts:       out,
	}, nil
}

type addPodcastRequest struct {
	CountryCode  string            `json:"country_code"`
	LanguageCode string            `json:"language_code"`
	TopicIDs     []content.TopicID `json:"topic_ids"`
	RSSFeedURL   string            `json:"rss_feed_url"`
	WebsiteURL   string            `json:"website_url"`
	Title        string            `json:"title"`
}

type addPodcastResponse struct {
	Error *string `json:"error,omitempty"`
}

func addPodcast(adminID admin.ID, r *router.Request) (interface{}, error) {
	var req addPodcastRequest
	if err := r.GetJSONBody(&req); err != nil {
		return nil, err
	}
	languageCode, err := wordsmith.GetLanguageCodeFromString(req.LanguageCode)
	if err != nil {
		return addPodcastResponse{
			Error: ptr.String(err.Error()),
		}, nil
	}
	countryCode, err := geo.GetCountryCodeFromString(req.CountryCode)
	if err != nil {
		return addPodcastResponse{
			Error: ptr.String(err.Error()),
		}, nil
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
		return addPodcastResponse{
			Error: ptr.String(err.Error()),
		}, nil
	}
	return addPodcastResponse{
		Error: nil,
	}, nil
}
