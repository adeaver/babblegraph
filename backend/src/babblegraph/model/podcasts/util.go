package podcasts

import (
	"babblegraph/model/content"
	"fmt"
)

func makePodcastEpisodeID(sourceID content.SourceID, episodeGUID string) EpisodeID {
	return EpisodeID(fmt.Sprintf("podcast_%s_%s", sourceID, episodeGUID))
}
