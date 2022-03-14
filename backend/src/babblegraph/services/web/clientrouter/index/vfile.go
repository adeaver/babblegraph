package index

import (
	"babblegraph/model/content"
	"babblegraph/model/podcasts"
	"babblegraph/model/virtualfile"
	"babblegraph/services/web/router"
	"babblegraph/util/cache"
	"babblegraph/util/database"
	"babblegraph/util/ptr"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
)

func handleVirtualFile(r *router.Request) (interface{}, error) {
	virtualFileName, err := r.GetRouteVar("fileName")
	if err != nil {
		return nil, err
	}
	objectID, virtualFileType, err := virtualfile.GetObjectIDAndType(*virtualFileName)
	if err != nil {
		return nil, err
	}
	switch *virtualFileType {
	case virtualfile.TypePodcast:
		podcastEpisodeID := podcasts.EpisodeID(*objectID)
		var episode podcasts.Episode
		if err := cache.WithCache(fmt.Sprintf("podcast-episode-%s", podcastEpisodeID.Str()), &episode, 12*time.Hour, func() (interface{}, error) {
			var source *content.Source
			if err := database.WithTx(func(tx *sqlx.Tx) error {
				var err error
				source, err = content.GetSource(tx, podcasts.GetSourceIDFromEpisodeID(podcastEpisodeID))
				return err
			}); err != nil {
				return nil, err
			}
			episode, err := podcasts.GetEpisodeByID(source.LanguageCode, podcastEpisodeID)
			if err != nil {
				return nil, err
			}
			return episode, nil
		}); err != nil {
			return nil, err
		}
		return ptr.String(episode.AudioFile.URL), nil
	case virtualfile.TypePodcastImage:
		sourceID := content.SourceID(*objectID)
		var podcastMetadata *podcasts.PodcastMetadata
		err := database.WithTx(func(tx *sqlx.Tx) error {
			var err error
			podcastMetadata, err = podcasts.GetPodcastMetadataForSourceID(tx, sourceID)
			return err
		})
		switch {
		case err != nil:
			return nil, err
		case podcastMetadata.ImageURL == nil:
			return nil, fmt.Errorf("Podcast %s does not have an image url", sourceID)
		default:
			return podcastMetadata.ImageURL, nil
		}
	default:
		return nil, fmt.Errorf("Got unsupported vfile type %s", *virtualFileType)
	}
}
