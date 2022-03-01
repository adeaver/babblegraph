package podcasts

import (
	"babblegraph/model/admin"
	podcastsearch "babblegraph/model/podcasts/search"
	"babblegraph/services/web/adminrouter/middleware"
	"babblegraph/services/web/router"
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
