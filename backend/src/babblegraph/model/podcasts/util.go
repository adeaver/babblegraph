package podcasts

import (
	"babblegraph/model/content"
	"fmt"
	"strings"
)

func makePodcastEpisodeID(sourceID content.SourceID, episodeGUID string) EpisodeID {
	return EpisodeID(fmt.Sprintf("podcast_%s_%s", sourceID, episodeGUID))
}

func GetSourceIDFromEpisodeID(id EpisodeID) content.SourceID {
	withoutPrefix := strings.TrimPrefix(id.Str(), "podcast_")
	parts := strings.Split(withoutPrefix, "_")
	return content.SourceID(parts[0])
}
