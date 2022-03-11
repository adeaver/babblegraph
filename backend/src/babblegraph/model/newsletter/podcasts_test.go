package newsletter

import (
	"babblegraph/model/content"
	"babblegraph/model/email"
	"babblegraph/model/podcasts"
	"babblegraph/model/usernewsletterpreferences"
	"babblegraph/model/userpodcasts"
	"babblegraph/util/ctx"
	"babblegraph/util/ptr"
	"babblegraph/wordsmith"
	"time"
)

type testPodcastAccessor struct {
	languageCode              wordsmith.LanguageCode
	userNewsletterPreferences usernewsletterpreferences.UserNewsletterPreferences
	validSourceIDs            []content.SourceID

	podcastEpisodes []podcasts.Episode
}

func (t *testPodcastAccessor) LookupPodcastEpisodesForTopics(topics []content.TopicID) (map[content.TopicID][]podcasts.Episode, error) {
	// This is a hack
	c := ctx.GetDefaultLogContext()
	podcastPreferences := t.userNewsletterPreferences.PodcastPreferences
	episodesByTopic := make(map[content.TopicID][]podcasts.Episode)
	for _, ep := range t.podcastEpisodes {
		switch {
		case !podcastPreferences.IncludeExplicitPodcasts && ep.IsExplicit:
			c.Debugf("Filtering out podcast because of explicit tag")
		case podcastPreferences.MinimumDurationNanoseconds != nil && ep.DurationNanoseconds < *podcastPreferences.MinimumDurationNanoseconds:
			c.Debugf("Filtering out podcast because of it's too short")
		case podcastPreferences.MaximumDurationNanoseconds != nil && ep.DurationNanoseconds > *podcastPreferences.MaximumDurationNanoseconds:
			c.Debugf("Filtering out podcast because of it's too long")
		case !isSourceValid(ep.SourceID.Ptr(), t.validSourceIDs):
			c.Debugf("Filtering out podcast because the source is not valid")
		// no-op
		default:
			for _, t := range topics {
				if containsTopic(t, ep.TopicIDs) {
					episodesByTopic[t] = append(episodesByTopic[t], ep)
				}
			}
		}
	}
	return episodesByTopic, nil
}

func (t *testPodcastAccessor) GetPodcastMetadataForSourceID(sourceID content.SourceID) (*podcasts.PodcastMetadata, error) {
	return &podcasts.PodcastMetadata{
		ImageURL:  ptr.String("https://static.babblegraph.com"),
		ContentID: sourceID,
	}, nil
}

func (t *testPodcastAccessor) InsertUserPodcastAndGetID(emailRecordID email.ID, episode podcasts.Episode) (*userpodcasts.ID, error) {
	return userpodcasts.ID("test-podcast").Ptr(), nil
}

func getDefaultPodcast(topic content.TopicID) podcasts.Episode {
	return podcasts.Episode{
		ID:                  podcasts.EpisodeID("test-podcast"),
		Title:               "Test Podcast",
		DurationNanoseconds: time.Hour,
		IsExplicit:          true,
		LanguageCode:        wordsmith.LanguageCodeSpanish,
		TopicIDs:            []content.TopicID{topic},
		SourceID:            content.SourceID("test-source"),
	}
}
