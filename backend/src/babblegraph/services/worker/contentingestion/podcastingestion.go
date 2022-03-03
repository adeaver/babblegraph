package contentingestion

import (
	"babblegraph/model/content"
	"babblegraph/services/worker/contentingestion/ingestrss"
	"babblegraph/util/ctx"
)

func processPodcastRSS1SourceSeed(c ctx.LogContext, sourceSeed content.SourceSeed) error {
	if err := ingestrss.GetPodcastDataForRSSFeed(c, sourceSeed.URL); err != nil {
		return err
	}
	return nil
}
