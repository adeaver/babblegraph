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
