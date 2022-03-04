package podcasts

import (
	"babblegraph/model/content"
	"babblegraph/util/ctx"
	"babblegraph/util/elastic"
	"babblegraph/util/ptr"
	"babblegraph/wordsmith"
	"fmt"
	"time"
)

type podcastIndex struct {
	languageCode wordsmith.LanguageCode
}

func getPodcastIndexForLanguageCode(languageCode wordsmith.LanguageCode) *podcastIndex {
	return &podcastIndex{
		languageCode: languageCode,
	}
}

func (p *podcastIndex) GetName() string {
	return fmt.Sprintf("podcasts_episodes_%s", p.languageCode.Str())
}

func (p *podcastIndex) ValidateDocument(document interface{}) error {
	if _, ok := document.(Episode); !ok {
		return fmt.Errorf("could not validate interface %+v to be of type podcasts.Episode", document)
	}
	return nil
}

func (p *podcastIndex) GenerateIDForDocument(document interface{}) (*string, error) {
	episode, ok := document.(Episode)
	if !ok {
		return nil, fmt.Errorf("could not validate interface %+v, to be of type podcasts.Episode", document)
	}
	id := makePodcastEpisodeID(episode.SourceID, episode.GUID)
	return ptr.String(string(id)), nil
}

type IndexPodcastEpisodeInput struct {
	Title           string
	Description     string
	PublicationDate time.Time
	EpisodeType     string
	Duration        time.Duration
	IsExplicit      bool
	AudioFile       AudioFile
	GUID            string

	Version Version

	LanguageCode wordsmith.LanguageCode
	TopicIDs     []content.TopicID
	SourceID     content.SourceID
}

func AssignIDAndIndexPodcastEpisode(c ctx.LogContext, input IndexPodcastEpisodeInput) (*EpisodeID, error) {
	podcastIndex := getPodcastIndexForLanguageCode(input.LanguageCode)
	episodeID := makePodcastEpisodeID(input.SourceID, input.GUID)
	if err := elastic.IndexDocument(c, podcastIndex, Episode{
		ID:              episodeID,
		Title:           input.Title,
		Description:     input.Description,
		PublicationDate: input.PublicationDate,
		EpisodeType:     input.EpisodeType,
		Duration:        input.Duration,
		IsExplicit:      input.IsExplicit,
		AudioFile:       input.AudioFile,
		GUID:            input.GUID,
		Version:         input.Version,
		LanguageCode:    input.LanguageCode,
		TopicIDs:        input.TopicIDs,
		SourceID:        input.SourceID,
	}); err != nil {
		return nil, err
	}
	return &episodeID, nil
}
