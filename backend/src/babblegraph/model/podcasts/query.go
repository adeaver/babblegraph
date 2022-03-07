package podcasts

import (
	"babblegraph/model/content"
	"babblegraph/util/elastic/esquery"
	"babblegraph/util/math/decimal"
	"babblegraph/wordsmith"
	"encoding/json"
	"fmt"

	"github.com/jmoiron/sqlx"
)

const (
	upsertPodcastMetadataQuery = `INSERT INTO content_podcast_metadata (
        content_id, image_url
    ) VALUES (
        $1, $2
    ) ON CONFLICT (content_id) DO UPDATE
    SET image_url=$2`
	getPodcastMetadataBySourceIDQuery = "SELECT * FROM content_podcast_metadata WHERE content_id = $1"
)

func UpsertPodcastMetadata(tx *sqlx.Tx, sourceID content.SourceID, imageURL string) error {
	if _, err := tx.Exec(upsertPodcastMetadataQuery, sourceID, imageURL); err != nil {
		return err
	}
	return nil
}

func GetPodcastMetadataForSourceID(tx *sqlx.Tx, sourceID content.SourceID) (*PodcastMetadata, error) {
	var matches []dbPodcastMetadata
	err := tx.Select(&matches, getPodcastMetadataBySourceIDQuery, sourceID)
	switch {
	case err != nil:
		return nil, err
	case len(matches) == 0,
		len(matches) > 1:
		return nil, fmt.Errorf("Expected exactly one match for podcast metadata for source ID %s, but got %d", sourceID, len(matches))
	default:
		out := matches[0].ToNonDB()
		return &out, nil
	}
}

func GetEpisodeByID(languageCode wordsmith.LanguageCode, id EpisodeID) (*Episode, error) {
	podcastIndex := getPodcastIndexForLanguageCode(languageCode)
	var episodes []Episode
	err := esquery.ExecuteSearch(podcastIndex, esquery.MatchPhrase("id", id), nil, func(source []byte, score decimal.Number) error {
		var episode Episode
		if err := json.Unmarshal(source, &episode); err != nil {
			return err
		}
		episodes = append(episodes, episode)
		return nil
	})
	switch {
	case err != nil:
		return nil, err
	case len(episodes) == 0,
		len(episodes) > 1:
		return nil, fmt.Errorf("Expected exactly one episode with id %s, but got %d", id, len(episodes))
	default:
		out := episodes[0]
		return &out, nil
	}
}
