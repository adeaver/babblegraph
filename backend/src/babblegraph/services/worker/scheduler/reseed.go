package scheduler

import (
	"babblegraph/model/content"
	"babblegraph/model/links2"
	"babblegraph/model/urltopicmapping"
	"babblegraph/services/worker/contentingestion/ingesthtml"
	"babblegraph/util/async"
	"babblegraph/util/ctx"
	"babblegraph/util/database"
	"babblegraph/util/urlparser"

	"github.com/jmoiron/sqlx"
)

type fetchNewLinksTask struct {
	SourceID     *content.SourceID
	SourceSeedID *content.SourceSeedID
	URL          string
}

func fetchNewLinksForSeedURLs(c async.Context) {
	c.Infof("Starting refetch of seed domains...")
	var htmlIngestSources []content.Source
	var sourceSeeds []content.SourceSeed
	if err := database.WithTx(func(tx *sqlx.Tx) error {
		var err error
		htmlIngestSources, err = content.LookupSourcesForIngestStrategy(tx, content.IngestStrategyWebsiteHTML1)
		if err != nil {
			return err
		}
		for _, s := range htmlIngestSources {
			seedsForSource, err := content.LookupActiveSourceSeedsForSource(tx, s.ID)
			if err != nil {
				return err
			}
			sourceSeeds = append(sourceSeeds, seedsForSource...)
		}
		return nil
	}); err != nil {
		c.Errorf("Error getting sources and seeds: %s", err.Error())
		return
	}
	tasksBySourceID := make(map[content.SourceID][]fetchNewLinksTask)
	for _, source := range htmlIngestSources {
		var newTasks []fetchNewLinksTask
		if source.ShouldUseURLAsSeedURL {
			newTasks = append(newTasks, fetchNewLinksTask{
				SourceID: source.ID.Ptr(),
				URL:      source.URL,
			})
		}
		tasksBySourceID[source.ID] = newTasks
	}
	for _, seed := range sourceSeeds {
		tasksBySourceID[seed.RootID] = append(tasksBySourceID[seed.RootID], fetchNewLinksTask{
			SourceSeedID: seed.ID.Ptr(),
			URL:          seed.URL,
		})
	}
	for len(tasksBySourceID) > 0 {
		for sourceID, tasks := range tasksBySourceID {
			c.Infof("Processing task for source ID %s", sourceID)
			task := tasks[0]
			if err := processRefetchTask(c, sourceID, task); err != nil {
				c.Errorf("Error processing task %+v: %s", task, err.Error())
			}
			nextTasks := append([]fetchNewLinksTask{}, tasks[1:]...)
			if len(nextTasks) == 0 {
				delete(tasksBySourceID, sourceID)
			} else {
				tasksBySourceID[sourceID] = nextTasks
			}
		}
	}
}

func processRefetchTask(c ctx.LogContext, sourceID content.SourceID, task fetchNewLinksTask) error {
	var source *content.Source
	var sourceFilter *content.SourceFilter
	if err := database.WithTx(func(tx *sqlx.Tx) error {
		var err error
		source, err = content.GetSource(tx, sourceID)
		if err != nil {
			return err
		}
		sourceFilter, err = content.LookupSourceFilterForSource(tx, sourceID)
		return err
	}); err != nil {
		return err
	}
	parsedHTMLPage, err := ingesthtml.ProcessURL(ingesthtml.ProcessURLInput{
		URL:          task.URL,
		Source:       *source,
		SourceFilter: sourceFilter,
	})
	if err != nil {
		return err
	}
	var parsedURLs []urlparser.ParsedURL
	for _, u := range parsedHTMLPage.Links {
		parsedURL := urlparser.ParseURL(u)
		if parsedURL == nil {
			continue
		}
		parsedURLs = append(parsedURLs, *parsedURL)
	}
	return database.WithTx(func(tx *sqlx.Tx) error {
		urlIdentifierHashSet := make(map[string]bool)
		var filteredURLs []links2.URLWithSourceMapping
		for _, u := range parsedURLs {
			if _, ok := urlIdentifierHashSet[u.URLIdentifier]; ok {
				continue
			}
			sourceID, _, err := content.LookupSourceIDForDomain(tx, u.Domain)
			switch {
			case err != nil:
				return err
			case sourceID == nil:
				// no-op
			default:
				urlIdentifierHashSet[u.URLIdentifier] = true
				filteredURLs = append(filteredURLs, links2.URLWithSourceMapping{
					URL:      u,
					SourceID: *sourceID,
				})
			}
		}
		c.Debugf("Got %d urls to insert", len(filteredURLs))
		if len(filteredURLs) == 0 {
			return nil
		}
		if err := links2.UpsertURLMappingsWithEmptyFetchStatus(tx, filteredURLs, true); err != nil {
			return err
		}
		if task.SourceSeedID == nil {
			return nil
		}
		topicMappingIDs, topicIDs, err := content.LookupTopicMappingIDForSourceSeedID(tx, *task.SourceSeedID)
		switch {
		case err != nil:
			return err
		case len(topicMappingIDs) == 0:
			c.Warnf("No topic mapping IDs for source seed: %s", *task.SourceSeedID)
			return nil
		}
		var topicMappingUnions []urltopicmapping.TopicMappingUnion
		for idx, topicMappingID := range topicMappingIDs {
			asContentTopic, err := content.GetContentTopicForTopicID(tx, topicIDs[idx])
			if err != nil {
				return err
			}
			topicMappingUnions = append(topicMappingUnions, urltopicmapping.TopicMappingUnion{
				Topic:          *asContentTopic,
				TopicMappingID: topicMappingID,
			})
		}
		for _, u := range filteredURLs {
			if err := urltopicmapping.ApplyContentTopicsToURL(tx, u.URL, topicMappingUnions); err != nil {
				return err
			}
		}
		return nil
	})
}
