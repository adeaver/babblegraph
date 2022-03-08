package newsletter

import (
	"babblegraph/model/content"
	"babblegraph/model/podcasts"
	"babblegraph/model/usernewsletterpreferences"
	"babblegraph/model/users"
	"babblegraph/wordsmith"

	"github.com/jmoiron/sqlx"
)

type podcastAccessor interface {
	LookupPodcastEpisodesForTopics(topics []content.TopicID) (map[content.TopicID][]podcasts.Episode, error)
}

type DefaultPodcastAccessor struct {
	languageCode              wordsmith.LanguageCode
	userNewsletterPreferences usernewsletterpreferences.UserNewsletterPreferences
	validSourceIDs            []content.SourceID
}

func GetDefaultPodcastAccessor(tx *sqlx.Tx, languageCode wordsmith.LanguageCode, userID users.UserID) (*DefaultPodcastAccessor, error) {
	validSourceIDs, err := content.LookupActiveSourceIDsByType(tx, content.SourceTypePodcast)
	if err != nil {
		return nil, err
	}
	userNewsletterPreferences, err := usernewsletterpreferences.GetUserNewsletterPrefrencesForLanguage(tx, userID, languageCode)
	if err != nil {
		return nil, err
	}
	return &DefaultPodcastAccessor{
		languageCode:              languageCode,
		validSourceIDs:            validSourceIDs,
		userNewsletterPreferences: *userNewsletterPreferences,
	}, nil
}

func (d *DefaultPodcastAccessor) LookupPodcastEpisodesForTopics(topicIDs []content.TopicID) (map[content.TopicID][]podcasts.Episode, error) {
	if !d.userNewsletterPreferences.PodcastPreferences.ArePodcastsEnabled {
		return nil, nil
	}
	out := make(map[content.TopicID][]podcasts.Episode)
	for _, t := range topicIDs {
		scoredEpisodes, err := podcasts.QueryEpisodes(d.languageCode, podcasts.QueryEpisodesInput{
			ValidSourceIDs:          d.validSourceIDs,
			TopicID:                 t,
			IncludeExplicitPodcasts: d.userNewsletterPreferences.PodcastPreferences.IncludeExplicitPodcasts,
			MinDurationNanoseconds:  d.userNewsletterPreferences.PodcastPreferences.MaximumDurationNanoseconds,
			MaxDurationNanoseconds:  d.userNewsletterPreferences.PodcastPreferences.MaximumDurationNanoseconds,
		})
		if err != nil {
			return nil, err
		}
		for _, ep := range scoredEpisodes {
			out[t] = append(out[t], ep.Episode)
		}
	}
	return out, nil
}
