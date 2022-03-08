package contentingestion

import (
	"babblegraph/model/content"
	"babblegraph/model/podcasts"
	"babblegraph/services/worker/contentingestion/ingestrss"
	"babblegraph/util/ctx"
	"babblegraph/util/database"
	"babblegraph/util/ptr"
	"strconv"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
)

func processPodcastRSS1SourceSeed(c ctx.LogContext, sourceSeed content.SourceSeed) error {
	channel, err := ingestrss.GetPodcastDataForRSSFeed(c, sourceSeed.URL)
	if err != nil {
		return err
	}
	var topicIDs []content.TopicID
	var source *content.Source
	if err := database.WithTx(func(tx *sqlx.Tx) error {
		var err error
		_, topicIDs, err = content.LookupTopicMappingIDForSourceSeedID(tx, sourceSeed.ID)
		if err != nil {
			return err
		}
		source, err = content.GetSource(tx, sourceSeed.RootID)
		if err != nil {
			return err
		}
		if err := content.UpdateSource(tx, sourceSeed.RootID, content.UpdateSourceInput{
			LanguageCode:          source.LanguageCode,
			Title:                 channel.Title,
			URL:                   source.URL,
			Type:                  source.Type,
			IngestStrategy:        source.IngestStrategy,
			IsActive:              source.IsActive,
			MonthlyAccessLimit:    source.MonthlyAccessLimit,
			Country:               source.Country,
			ShouldUseURLAsSeedURL: source.ShouldUseURLAsSeedURL,
		}); err != nil {
			return err
		}
		return podcasts.UpsertPodcastMetadata(tx, sourceSeed.RootID, channel.Image.URL)
	}); err != nil {
		return err
	}
	for _, episode := range channel.Episodes {
		toIndex, err := convertIngestEpisodeToModelEpisode(c, episode)
		if err != nil {
			c.Warnf("Error converting episode with GUID %s for source id %s: %s", episode.ID, sourceSeed.RootID, err.Error())
			continue
		}
		toIndex.Version = podcasts.CurrentVersion
		toIndex.SourceID = source.ID
		toIndex.LanguageCode = source.LanguageCode
		toIndex.TopicIDs = topicIDs
		podcastEpisodeID, err := podcasts.AssignIDAndIndexPodcastEpisode(c, *toIndex)
		if err != nil {
			c.Errorf("Error indexing podcast for source %s with GUID %s: %s", source.ID, toIndex.GUID)
			continue
		}
		c.Infof("Indexed podcast episode with ID %s", *podcastEpisodeID)
	}
	return nil
}

func convertIngestEpisodeToModelEpisode(c ctx.LogContext, in ingestrss.PodcastEpisode) (*podcasts.IndexPodcastEpisodeInput, error) {
	publicationDate, err := time.Parse("Mon, 2 January 2006 15:04:05 MST", in.PublicationDate)
	if err != nil {
		c.Infof("Got error parsing time, trying different format")
		var fErr error
		publicationDate, fErr = time.Parse("Mon, 2 January 2006 15:04:05 -0700", in.PublicationDate)
		if fErr != nil {
			c.Errorf("Could not parse time. Original error %s", err.Error())
			return nil, fErr
		}
	}
	episodeDuration, err := parseDuration(in.Duration)
	if err != nil {
		return nil, err
	}
	return &podcasts.IndexPodcastEpisodeInput{
		Title:               in.Title,
		Description:         in.Description,
		PublicationDate:     publicationDate,
		EpisodeType:         in.EpisodeType.Str(),
		DurationNanoseconds: *episodeDuration,
		IsExplicit:          in.IsExplicit.ToBool(c),
		GUID:                in.ID,
		AudioFile: podcasts.AudioFile{
			URL:  in.AudioData.URL,
			Type: in.AudioData.Type,
		},
	}, nil
}

func parseDuration(durationStr string) (*time.Duration, error) {
	if strings.Contains(durationStr, ":") {
		parts := strings.Split(durationStr, ":")
		validDurations := []time.Duration{time.Second, time.Minute, time.Hour}
		var out time.Duration
		for idx := 0; idx < len(parts); idx++ {
			currentPart := parts[len(parts)-(idx+1)]
			magnitude, err := strconv.Atoi(currentPart)
			if err != nil {
				return nil, err
			}
			out = out + time.Duration(magnitude)*validDurations[idx]
		}
		return ptr.Duration(out), nil
	}
	numberOfSeconds, err := strconv.Atoi(durationStr)
	if err != nil {
		return nil, err
	}
	return ptr.Duration(time.Duration(numberOfSeconds) * time.Second), nil
}
