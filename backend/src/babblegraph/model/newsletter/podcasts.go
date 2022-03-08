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
	LookupPodcastEpisodeForNewsletter(topics []content.TopicID) (*podcasts.Episode, *content.TopicID, error)
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

func (d *DefaultPodcastAccessor) LookupPodcastEpisodesForNewsletter(topicIDs []content.TopicID) (*podcasts.Episode, *content.TopicID, error) {
	if !d.userNewsletterPreferences.PodcastPreferences.ArePodcastsEnabled {
		return nil, nil, nil
	}
	var topicID *content.TopicID
	var currentEpisode *podcasts.ScoredEpisode
	for _, t := range topicIDs {
		t := t
		scoredEpisodes, err := podcasts.QueryEpisodes(d.languageCode, podcasts.QueryEpisodesInput{
			ValidSourceIDs:          d.validSourceIDs,
			TopicID:                 t,
			IncludeExplicitPodcasts: d.userNewsletterPreferences.PodcastPreferences.IncludeExplicitPodcasts,
			MinDurationNanoseconds:  d.userNewsletterPreferences.PodcastPreferences.MaximumDurationNanoseconds,
			MaxDurationNanoseconds:  d.userNewsletterPreferences.PodcastPreferences.MaximumDurationNanoseconds,
		})
		if err != nil {
			return nil, nil, err
		}
		for _, ep := range scoredEpisodes {
			ep := ep
			if currentEpisode == nil || ep.Score.GreaterThan(currentEpisode.Score) {
				currentEpisode = &ep
				topicID = t.Ptr()
			}
		}
	}
	return &currentEpisode.Episode, topicID, nil
}
