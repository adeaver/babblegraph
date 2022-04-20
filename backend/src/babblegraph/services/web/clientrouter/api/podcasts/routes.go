package podcasts

import (
	"babblegraph/model/content"
	"babblegraph/model/podcasts"
	"babblegraph/model/useraccounts"
	"babblegraph/model/userpodcasts"
	"babblegraph/model/virtualfile"
	"babblegraph/services/web/clientrouter/clienterror"
	"babblegraph/services/web/clientrouter/routermiddleware"
	"babblegraph/services/web/router"
	"babblegraph/util/cache"
	"babblegraph/util/database"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
)

var Routes = router.RouteGroup{
	Prefix: "podcasts",
	Routes: []router.Route{
		{
			Path: "get_podcast_metadata_1",
			Handler: routermiddleware.WithNoBodyRequestLogger(
				getPodcastMetadata,
			),
		},
	},
}

type getPodcastMetadataRequest struct {
	UserPodcastID userpodcasts.ID `json:"user_podcast_id"`
}

type getPodcastMetadataResponse struct {
	Error    *clienterror.Error `json:"error,omitempty"`
	Metadata *podcastMetadata   `json:"metadata,omitempty"`
}

type podcastMetadata struct {
	PodcastTitle       string  `json:"podcast_title"`
	EpisodeTitle       string  `json:"episode_title"`
	EpisodeDescription string  `json:"episode_description"`
	PodcastURL         string  `json:"podcast_url"`
	AudioURL           string  `json:"audio_url"`
	ImageURL           *string `json:"image_url"`
}

func getPodcastMetadata(r *router.Request) (interface{}, error) {
	var req getPodcastMetadataRequest
	if err := r.GetJSONBody(&req); err != nil {
		return nil, err
	}
	var cErr *clienterror.Error
	var source *content.Source
	var showMetadata *podcasts.PodcastMetadata
	var userPodcast *userpodcasts.UserPodcast
	err := database.WithTx(func(tx *sqlx.Tx) error {
		var err error
		userPodcast, err = userpodcasts.GetByID(tx, req.UserPodcastID)
		if err != nil {
			return err
		}
		subscriptionLevel, err := useraccounts.LookupSubscriptionLevelForUser(tx, userPodcast.UserID)
		switch {
		case err != nil:
			return err
		case subscriptionLevel == nil,
			*subscriptionLevel == useraccounts.SubscriptionLevelLegacy:
			cErr = clienterror.ErrorNoAuth.Ptr()
			return nil
		default:
			// no-op
		}
		source, err = content.GetSource(tx, userPodcast.SourceID)
		if err != nil {
			return err
		}
		if err := userpodcasts.RegisterOpenedPodcast(tx, userPodcast.ID); err != nil {
			return err
		}
		showMetadata, err = podcasts.GetPodcastMetadataForSourceID(tx, userPodcast.SourceID)
		return err
	})
	switch {
	case err != nil:
		return nil, err
	case cErr != nil:
		return getPodcastMetadataResponse{
			Error: cErr,
		}, nil
	default:
		var episode podcasts.Episode
		if err := cache.WithCache(fmt.Sprintf("podcast-episode-%s", userPodcast.EpisodeID.Str()), &episode, 12*time.Hour, func() (interface{}, error) {
			episode, err := podcasts.GetEpisodeByID(source.LanguageCode, userPodcast.EpisodeID)
			if err != nil {
				return nil, err
			}
			return episode, nil
		}); err != nil {
			return nil, err
		}
		audioURL, err := virtualfile.GetVirtualFileURL(episode.ID.Str(), virtualfile.TypePodcast)
		if err != nil {
			return nil, err
		}
		var imageURL *string
		if showMetadata.ImageURL != nil {
			imageURL, err = virtualfile.GetVirtualFileURL(source.ID.Str(), virtualfile.TypePodcastImage)
			if err != nil {
				return nil, err
			}
		}
		return getPodcastMetadataResponse{
			Metadata: &podcastMetadata{
				PodcastTitle:       source.Title,
				EpisodeTitle:       episode.Title,
				EpisodeDescription: episode.Description,
				PodcastURL:         source.URL,
				AudioURL:           *audioURL,
				ImageURL:           imageURL,
			},
		}, nil
	}
}
