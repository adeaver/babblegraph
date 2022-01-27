package tasks

import (
	"babblegraph/model/content"
	"babblegraph/model/contenttopics"
	"babblegraph/model/domains"
	"babblegraph/util/ctx"
	"babblegraph/util/database"
	"babblegraph/util/ptr"
	"babblegraph/util/urlparser"
	"babblegraph/wordsmith"

	"github.com/jmoiron/sqlx"
)

func SyncHardcodedContent() error {
	c := ctx.GetDefaultLogContext()
	return database.WithTx(func(tx *sqlx.Tx) error {
		topicsToTopicID := make(map[contenttopics.ContentTopic]content.TopicID)
		// Add all content topics
		c.Infof("Adding all content topics")
		for _, t := range getAllContentTopics() {
			c.Infof("Adding topic: %s", t)
			topicID, err := insertContentTopic(tx, t)
			if err != nil {
				return err
			}
			topicsToTopicID[t] = *topicID
		}
		domainsToSourceID := make(map[string]content.SourceID)
		// Add all sources
		c.Infof("Adding all sources")
		ds := domains.GetDomains()
		for _, d := range ds {
			c.Infof("Adding source: %s", d)
			sourceID, err := insertSource(tx, d)
			if err != nil {
				return err
			}
			domainsToSourceID[d] = *sourceID
		}
		// Add seed URLs
		c.Infof("Adding seed urls")
		for u, topics := range domains.GetSeedURLs() {
			parsedURL := urlparser.ParseURL(u)
			if parsedURL == nil {
				c.Infof("Seed URL %s did not parse", u)
				continue
			}
			sourceID, ok := domainsToSourceID[parsedURL.Domain]
			if !ok {
				c.Infof("Seed URL %s did not correspond to a source", u)
				continue
			}
			if parsedURL.Domain == u {
				c.Infof("Seed URL %s did not correspond is a source, continuing", u)
				continue
			}
			sourceSeedID, err := content.AddSourceSeed(tx, sourceID, *parsedURL, true)
			if err != nil {
				return err
			}
			for _, t := range topics {
				topicID, ok := topicsToTopicID[t]
				if !ok {
					c.Infof("Seed URL %s had topic %s that did not correspond to a topic", u, t)
					continue
				}
				if err := content.UpsertSourceSeedMapping(tx, *sourceSeedID, topicID, true); err != nil {
					return err
				}
			}
		}
		// TODO: update user link clicks and user content topic mappings
		return nil
	})
}

func insertSource(tx *sqlx.Tx, domain string) (*content.SourceID, error) {
	domainMetadata, err := domains.GetDomainMetadata(domain)
	if err != nil {
		return nil, err
	}
	sourceID, err := content.InsertSource(tx, content.InsertSourceInput{
		Title:                 domain,
		LanguageCode:          domainMetadata.LanguageCode,
		URL:                   domain,
		Type:                  content.SourceTypeNewsWebsite,
		IngestStrategy:        content.IngestStrategyWebsiteHTML1,
		Country:               domainMetadata.Country,
		ShouldUseURLAsSeedURL: true,
		IsActive:              true,
		MonthlyAccessLimit:    domainMetadata.NumberOfMonthlyFreeArticles,
	})
	if err != nil {
		return nil, err
	}
	if domainMetadata.PaywallValidation != nil {
		_, err = content.UpsertSourceFilterForSource(tx, *sourceID, content.UpsertSourceFilterForSourceInput{
			IsActive:            true,
			PaywallClasses:      domainMetadata.PaywallValidation.PaywallClasses,
			PaywallIDs:          domainMetadata.PaywallValidation.PaywallIDs,
			UseLDJSONValidation: ptr.Bool(domainMetadata.PaywallValidation.UseLDJSONValidation != nil),
		})
		if err != nil {
			return nil, err
		}
	}
	return sourceID, nil
}

func insertContentTopic(tx *sqlx.Tx, t contenttopics.ContentTopic) (*content.TopicID, error) {
	topicID, err := content.AddTopic(tx, t.Str(), true)
	if err != nil {
		return nil, err
	}
	displayName, err := contenttopics.ContentTopicNameToDisplayName(t)
	if err != nil {
		return nil, err
	}
	_, err = content.AddTopicDisplayName(tx, *topicID, wordsmith.LanguageCodeSpanish, displayName.Str(), true)
	return topicID, err
}

func getAllContentTopics() []contenttopics.ContentTopic {
	return []contenttopics.ContentTopic{
		contenttopics.ContentTopicArt,
		contenttopics.ContentTopicAstronomy,
		contenttopics.ContentTopicArchitecture,
		contenttopics.ContentTopicAutomotive,
		contenttopics.ContentTopicBusiness,
		contenttopics.ContentTopicCelebrityNews,
		contenttopics.ContentTopicCooking,
		contenttopics.ContentTopicCulture,
		contenttopics.ContentTopicCurrentEventsArgentina,
		contenttopics.ContentTopicCurrentEventsChile,
		contenttopics.ContentTopicCurrentEventsColombia,
		contenttopics.ContentTopicCurrentEventsCostaRica,
		contenttopics.ContentTopicCurrentEventsElSalvador,
		contenttopics.ContentTopicCurrentEventsGuatemala,
		contenttopics.ContentTopicCurrentEventsHonduras,
		contenttopics.ContentTopicCurrentEventsMexico,
		contenttopics.ContentTopicCurrentEventsNicaragua,
		contenttopics.ContentTopicCurrentEventsPanama,
		contenttopics.ContentTopicCurrentEventsParaguay,
		contenttopics.ContentTopicCurrentEventsPeru,
		contenttopics.ContentTopicCurrentEventsSpain,
		contenttopics.ContentTopicCurrentEventsUnitedStates,
		contenttopics.ContentTopicCurrentEventsVenezuela,
		contenttopics.ContentTopicCurrentEventsUruguay,
		contenttopics.ContentTopicCurrentEventsPuertoRico,
		contenttopics.ContentTopicEconomy,
		contenttopics.ContentTopicEntertainment,
		contenttopics.ContentTopicEnvironment,
		contenttopics.ContentTopicFashion,
		contenttopics.ContentTopicFilm,
		contenttopics.ContentTopicFinance,
		contenttopics.ContentTopicHealth,
		contenttopics.ContentTopicHome,
		contenttopics.ContentTopicLifestyle,
		contenttopics.ContentTopicLiterature,
		contenttopics.ContentTopicMusic,
		contenttopics.ContentTopicOpinion,
		contenttopics.ContentTopicPolitics,
		contenttopics.ContentTopicScience,
		contenttopics.ContentTopicSports,
		contenttopics.ContentTopicTechnology,
		contenttopics.ContentTopicTheater,
		contenttopics.ContentTopicTravel,
		contenttopics.ContentTopicVideoGames,
		contenttopics.ContentTopicWorldNews,
	}
}
