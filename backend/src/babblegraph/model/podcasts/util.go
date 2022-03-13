package podcasts

import (
	"babblegraph/model/content"
	"encoding/base64"
	"fmt"
	"strings"
)

func makePodcastEpisodeID(sourceID content.SourceID, episodeGUID string) EpisodeID {
	// Make sure the ID is URL base64 encoded
	return EpisodeID(fmt.Sprintf("podcast_%s_%s", sourceID, base64.URLEncoding.EncodeToString([]byte(episodeGUID))))
}

func GetSourceIDFromEpisodeID(id EpisodeID) content.SourceID {
	withoutPrefix := strings.TrimPrefix(id.Str(), "podcast_")
	parts := strings.Split(withoutPrefix, "_")
	return content.SourceID(parts[0])
}
