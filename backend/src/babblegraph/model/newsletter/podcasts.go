package newsletter

import (
	"babblegraph/model/content"
	"babblegraph/model/email"
	"babblegraph/model/podcasts"
	"babblegraph/model/usernewsletterpreferences"
	"babblegraph/model/userpodcasts"
	"babblegraph/model/users"
	"babblegraph/wordsmith"

	"github.com/jmoiron/sqlx"
)

const (
	maxPodcastsPerTopic = 1
	// I didn't really know what to name this, but this is an arbitrary
	// value representing the number of articles that must exist in a category
	// before we replace article links with podcast links.
	// <= this number, podcast links will be appended
	// > this number, len(podcastLinks) links will be removed from a category
	podcastArticleRemovalBreakpoint = 2
)

type podcastAccessor interface {
	LookupPodcastEpisodesForTopics(topics []content.TopicID) (map[content.TopicID][]podcasts.Episode, error)
	GetPodcastMetadataForSourceID(sourceID content.SourceID) (*podcasts.PodcastMetadata, error)
	InsertUserPodcastAndGetID(emailRecordID email.ID, episode podcasts.Episode) (*userpodcasts.ID, error)
}

type DefaultPodcastAccessor struct {
	tx *sqlx.Tx

	seenPodcastIDs            []podcasts.EpisodeID
	languageCode              wordsmith.LanguageCode
	userID                    users.UserID
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
	userPodcasts, err := userpodcasts.GetByUserID(tx, userID)
	if err != nil {
		return nil, err
	}
	var seenPodcastIDs []podcasts.EpisodeID
	for _, u := range userPodcasts {
		seenPodcastIDs = append(seenPodcastIDs, u.EpisodeID)
	}
	return &DefaultPodcastAccessor{
		tx: tx,

		userID:                    userID,
		seenPodcastIDs:            seenPodcastIDs,
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
			SeenPodcastIDs:          d.seenPodcastIDs,
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

func (d *DefaultPodcastAccessor) GetPodcastMetadataForSourceID(sourceID content.SourceID) (*podcasts.PodcastMetadata, error) {
	return podcasts.GetPodcastMetadataForSourceID(d.tx, sourceID)
}

func (d *DefaultPodcastAccessor) InsertUserPodcastAndGetID(emailRecordID email.ID, episode podcasts.Episode) (*userpodcasts.ID, error) {
	return userpodcasts.InsertUserPodcastAndReturnID(d.tx, episode.ID, d.userID, episode.SourceID, emailRecordID)
}
