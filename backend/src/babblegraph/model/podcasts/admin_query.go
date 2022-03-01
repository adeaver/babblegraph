package podcasts

import (
	"babblegraph/model/content"
	"babblegraph/util/geo"
	"babblegraph/util/urlparser"
	"babblegraph/wordsmith"
	"fmt"

	"github.com/jmoiron/sqlx"
)

const (
	addPodcastMetadataQuery = "INSERT INTO content_podcast_metadata (content_id) VALUES ($1)"
)

type AddPodcastInput struct {
	CountryCode  geo.CountryCode
	LanguageCode wordsmith.LanguageCode
	WebsiteURL   string
	Title        string
	RSSFeedURL   string
	TopicIDs     []content.TopicID
}

func AddPodcast(tx *sqlx.Tx, input AddPodcastInput) error {
	parsedWebsite := urlparser.ParseURL(input.WebsiteURL)
	if parsedWebsite == nil {
		return fmt.Errorf("Website was incorrect")
	}
	websiteURL := fmt.Sprintf("%s/%s", parsedWebsite.Domain, parsedWebsite.Path)
	if parsedWebsite.Protocol != nil {
		websiteURL = fmt.Sprintf("%s//%s", *parsedWebsite.Protocol, websiteURL)
	}
	sourceID, err := content.InsertSource(tx, content.InsertSourceInput{
		Title:                 input.Title,
		LanguageCode:          input.LanguageCode,
		URL:                   websiteURL,
		Country:               input.CountryCode,
		ShouldUseURLAsSeedURL: false,
		IsActive:              true,
		MonthlyAccessLimit:    nil,
		Type:                  content.SourceTypePodcast,
		IngestStrategy:        content.IngestStrategyPodcastRSS1,
	})
	if err != nil {
		return err
	}
	if _, err := tx.Exec(addPodcastMetadataQuery, *sourceID); err != nil {
		return err
	}
	parsedURL := urlparser.ParseURL(input.RSSFeedURL)
	if parsedURL == nil {
		return fmt.Errorf("URL %s does not parse", input.RSSFeedURL)
	}
	sourceSeedID, err := content.AddSourceSeed(tx, *sourceID, *parsedURL, true)
	if err != nil {
		return err
	}
	for _, topicID := range input.TopicIDs {
		if err := content.UpsertSourceSeedMapping(tx, *sourceSeedID, topicID, true); err != nil {
			return err
		}
	}
	return nil
}
