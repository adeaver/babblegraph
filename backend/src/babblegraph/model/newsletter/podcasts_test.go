package newsletter

import (
	"babblegraph/model/content"
	"babblegraph/model/podcasts"
	"babblegraph/model/usernewsletterpreferences"
	"babblegraph/wordsmith"
)

type testPodcastAccessor struct {
	languageCode              wordsmith.LanguageCode
	userNewsletterPreferences usernewsletterpreferences.UserNewsletterPreferences
	validSourceIDs            []content.SourceID

	podcastEpisodes []podcasts.Episode
}

func (t *testPodcastAccessor) LookupPodcastEpisodesForTopics(topics []content.TopicID) (map[content.TopicID][]podcasts.Episode, error) {
	podcastPreferences := t.userNewsletterPreferences.PodcastPreferences
	episodesByTopic := make(map[content.TopicID][]podcasts.Episode)
	for _, ep := range t.podcastEpisodes {
		switch {
		case !podcastPreferences.IncludeExplicitPodcasts && ep.IsExplicit,
			podcastPreferences.MinimumDurationNanoseconds != nil && ep.DurationNanoseconds < *podcastPreferences.MinimumDurationNanoseconds,
			podcastPreferences.MaximumDurationNanoseconds != nil && ep.DurationNanoseconds > *podcastPreferences.MaximumDurationNanoseconds,
			!isSourceValid(ep.SourceID.Ptr(), t.validSourceIDs):
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
