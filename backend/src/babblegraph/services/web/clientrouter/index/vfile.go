package index

import (
	"babblegraph/model/content"
	"babblegraph/model/podcasts"
	"babblegraph/model/virtualfile"
	"babblegraph/util/cache"
	"babblegraph/util/database"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
)

func HandleVirtualFile(w http.ResponseWriter, r *http.Request) {
	routeVars := mux.Vars(r)
	virtualFileName, ok := routeVars["fileName"]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	objectID, virtualFileType, err := virtualfile.GetObjectIDAndType(virtualFileName)
	if err != nil {
		if !ok {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
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
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		http.Redirect(w, r, episode.AudioFile.URL, http.StatusFound)
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
			w.WriteHeader(http.StatusBadRequest)
			return
		case podcastMetadata.ImageURL == nil:
			w.WriteHeader(http.StatusBadRequest)
			return
		default:
			http.Redirect(w, r, *podcastMetadata.ImageURL, http.StatusFound)
		}
	default:
		w.WriteHeader(http.StatusBadRequest)
	}
}
