package podcasts

import (
	"babblegraph/model/content"
	"babblegraph/util/elastic/esquery"
	"babblegraph/util/math/decimal"
	"babblegraph/wordsmith"
	"encoding/json"
	"fmt"
	"time"

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

type QueryEpisodesInput struct {
	SeenPodcastIDs          []EpisodeID
	ValidSourceIDs          []content.SourceID
	TopicID                 content.TopicID
	IncludeExplicitPodcasts bool
	MaxDurationNanoseconds  *time.Duration
	MinDurationNanoseconds  *time.Duration
}

type ScoredEpisode struct {
	Score   decimal.Number
	Episode Episode
}

func QueryEpisodes(languageCode wordsmith.LanguageCode, input QueryEpisodesInput) ([]ScoredEpisode, error) {
	podcastIndex := getPodcastIndexForLanguageCode(languageCode)
	queryBuilder := esquery.NewBoolQueryBuilder()
	queryBuilder.AddMust(esquery.MatchPhrase("source_id", input.ValidSourceIDs))
	queryBuilder.AddMust(esquery.MatchPhrase("topic_ids", input.TopicID))
	versionRangeQueryBuilder := esquery.NewRangeQueryBuilderForFieldName("version")
	versionRangeQueryBuilder.GreaterThanOrEqualToInt64(Version1.Int64())
	versionRangeQueryBuilder.LessThanOrEqualToInt64(Version1.Int64())
	queryBuilder.AddMust(versionRangeQueryBuilder.BuildRangeQuery())
	if !input.IncludeExplicitPodcasts {
		queryBuilder.AddMust(esquery.Match("is_explicit", false))
	}
	if input.MaxDurationNanoseconds != nil || input.MinDurationNanoseconds != nil {
		durationQueryBuilder := esquery.NewRangeQueryBuilderForFieldName("duration_nanoseconds")
		if input.MaxDurationNanoseconds != nil {
			durationQueryBuilder.LessThanOrEqualToInt64(int64(*input.MaxDurationNanoseconds))
		}
		if input.MinDurationNanoseconds != nil {
			durationQueryBuilder.GreaterThanOrEqualToInt64(int64(*input.MinDurationNanoseconds))
		}
		queryBuilder.AddMust(durationQueryBuilder.BuildRangeQuery())
	}
	queryBuilder.AddMustNot(esquery.Match("id", input.SeenPodcastIDs))
	var out []ScoredEpisode
	if err := esquery.ExecuteSearch(podcastIndex, queryBuilder.BuildBoolQuery(), nil, func(source []byte, score decimal.Number) error {
		var episode Episode
		if err := json.Unmarshal(source, &episode); err != nil {
			return err
		}
		out = append(out, ScoredEpisode{
			Score:   score,
			Episode: episode,
		})
		return nil
	}); err != nil {
		return nil, err
	}
	return out, nil
}
