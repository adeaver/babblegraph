package contentingestion

import (
	"babblegraph/model/content"
	"babblegraph/model/links2"
	"babblegraph/util/async"
	"babblegraph/util/bufferedfetch"
	"babblegraph/util/ctx"
	"babblegraph/util/database"
	"babblegraph/util/ptr"
	"fmt"
	"sync"
	"time"

	"github.com/jmoiron/sqlx"
)

const (
	defaultChunkSize = 300
)

type ingestConfig struct {
	maxWorkers           int
	defaultTimeUntilFree time.Duration
	defaultRefreshPeriod time.Duration
}

var validIngestStrategies = map[content.IngestStrategy]ingestConfig{
	content.IngestStrategyWebsiteHTML1: {
		maxWorkers:           5,
		defaultTimeUntilFree: 10 * time.Second,
		defaultRefreshPeriod: 1 * time.Hour,
	},
	content.IngestStrategyPodcastRSS1: {
		maxWorkers:           2,
		defaultTimeUntilFree: 10 * time.Second,
		defaultRefreshPeriod: 24 * time.Hour,
	},
}

type ingestionSource struct {
	sourceID content.SourceID
	freeAt   time.Time
}

type ingestor struct {
	ingestionType  content.IngestStrategy
	mu             sync.Mutex
	orderedSources []ingestionSource
	sourceSet      map[content.SourceID]bool
}

func (i *ingestor) initialize(c ctx.LogContext) error {
	config, ok := validIngestStrategies[i.ingestionType]
	if !ok {
		return fmt.Errorf("Invalid ingestion type %s", i.ingestionType)
	}
	var sources []content.Source
	if err := database.WithTx(func(tx *sqlx.Tx) error {
		var err error
		sources, err = content.LookupSourcesForIngestStrategy(tx, i.ingestionType)
		return err
	}); err != nil {
		return err
	}
	var orderedSources []ingestionSource
	sourceSet := make(map[content.SourceID]bool)
	for _, s := range sources {
		orderedSources = append(orderedSources, ingestionSource{
			sourceID: s.ID,
			freeAt:   time.Now().Add(config.defaultTimeUntilFree),
		})
		sourceSet[s.ID] = true
		if err := i.registerBufferedFetchForSource(c, s.ID); err != nil {
			return err
		}
	}
	i.orderedSources = orderedSources
	i.sourceSet = sourceSet
	return nil
}

func (i *ingestor) registerBufferedFetchForSource(c ctx.LogContext, sourceID content.SourceID) error {
	bufferedFetchKey := i.getBufferedKeyFetchForSourceID(sourceID)
	switch i.ingestionType {
	case content.IngestStrategyWebsiteHTML1:
		if err := bufferedfetch.Register(bufferedFetchKey, func() (interface{}, error) {
			var links []links2.Link
			if err := database.WithTx(func(tx *sqlx.Tx) error {
				var err error
				links, err = links2.LookupBulkUnfetchedLinksForSourceID(tx, sourceID, defaultChunkSize)
				return err
			}); err != nil {
				return nil, err
			}
			return links, nil
		}); err != nil {
			return err
		}
		return bufferedfetch.ForceRefill(c, bufferedFetchKey)
	case content.IngestStrategyPodcastRSS1:
		// This is definitely not necessary, but it cleans up the code a lot and doesn't hurt
		if err := bufferedfetch.Register(bufferedFetchKey, func() (interface{}, error) {
			var sourceSeeds []content.SourceSeed
			if err := database.WithTx(func(tx *sqlx.Tx) error {
				var err error
				sourceSeeds, err = content.LookupActiveSourceSeedsForSource(tx, sourceID)
				return err
			}); err != nil {
				return nil, err
			}
			return sourceSeeds, nil
		}); err != nil {
			return err
		}
		return bufferedfetch.ForceRefill(c, bufferedFetchKey)
	default:
		return fmt.Errorf("Unsupported ingestion type %s", i.ingestionType)
	}
}

func (i *ingestor) getBufferedKeyFetchForSourceID(sourceID content.SourceID) string {
	return fmt.Sprintf("%s-%s", i.ingestionType, sourceID)
}

func (i *ingestor) processSources() func(c async.Context) {
	config, ok := validIngestStrategies[i.ingestionType]
	if !ok {
		panic(fmt.Sprintf("Invalid ingestion type %s", i.ingestionType))
	}
	return func(c async.Context) {
		c.Infof("Starting ingestor of type %s", i.ingestionType)
		workerManagerErrs := make(chan error)
		timer := time.NewTimer(config.defaultRefreshPeriod)
		async.WithContext(workerManagerErrs, fmt.Sprintf("%s-worker-manager", i.ingestionType), i.startWorkerManager(config.maxWorkers)).Start()
		for {
			select {
			case err := <-workerManagerErrs:
				c.Warnf("Error with %s worker manager: %s", i.ingestionType, err.Error())
				async.WithContext(workerManagerErrs, fmt.Sprintf("%s-worker-manager", i.ingestionType), i.startWorkerManager(config.maxWorkers)).Start()
			case _ = <-timer.C:
				c.Infof("Refreshing %s ingestor", i.ingestionType)
				if err := i.initialize(c); err != nil {
					c.Errorf("Error refreshing %s ingestor: %s", i.ingestionType, err.Error())
				}
				timer = time.NewTimer(config.defaultRefreshPeriod)
			}
		}
	}
}

func (i *ingestor) startWorkerManager(maxWorkers int) func(c async.Context) {
	return func(c async.Context) {
		threadComplete := make(chan interface{}, 1)
		workerErrs := make(chan error)
		timer := time.NewTimer(1 * time.Second)
		var numWorkers int
		spinOffWorkerOrWait := func() (_shouldBreak bool) {
			task, waitDuration, err := i.getTask(c)
			switch {
			case err != nil:
				c.Errorf("Error getting task for %s ingestor: %s. Will retry", i.ingestionType, err.Error())
				waitDuration = ptr.Duration(time.Duration(2 * time.Minute))
			case waitDuration != nil:
				c.Infof("Waiting")
			case task != nil:
				numWorkers++
				async.WithContext(workerErrs, fmt.Sprintf("%s-worker", i.ingestionType), i.processTask(workerErrs, threadComplete, task)).Start()
				return false
			default:
				// This should probably not happen
				c.Warnf("Got null task for ingestion type %s, but no wait time or error either", i.ingestionType)
				waitDuration = ptr.Duration(time.Duration(10 * time.Second))
			}
			timer = time.NewTimer(*waitDuration)
			return true
		}
		// Initialize workers
		for i := 0; i < maxWorkers; i++ {
			if shouldBreak := spinOffWorkerOrWait(); shouldBreak {
				break
			}
		}
		task, waitDuration, err := i.getTask(c)
		if err != nil {
			c.Errorf("Error getting tasks for ingestion type %s: %s", i.ingestionType, err.Error())
		}
		for {
			select {
			case _ = <-threadComplete:
				c.Infof("Thread is complete")
				numWorkers--
			case err = <-workerErrs:
				c.Infof("Thread for ingestion type %s encountered error: %s", i.ingestionType, err.Error())
				numWorkers--
			case _ = <-timer.C:
				c.Infof("Worker manager for %s ingestor timer has finished. Currently there are %d workers", i.ingestionType, numWorkers)
				switch {
				case numWorkers == maxWorkers && task != nil:
					c.Infof("All workers are currently busy and there is work to do, continuing...")
					continue
				case numWorkers == maxWorkers && task == nil:
					c.Infof("All workers are busy, but link needs replenshing")
					task, waitDuration, err = i.getTask(c)
					if err != nil {
						c.Errorf("Error getting link: %s", err.Error())
						waitDuration = ptr.Duration(time.Duration(2 * time.Minute))
					}
					if waitDuration != nil {
						timer = time.NewTimer(*waitDuration)
					}
					continue
				case numWorkers < maxWorkers && task == nil:
					c.Infof("Timer is complete")
				}
			}
			if task != nil {
				async.WithContext(workerErrs, fmt.Sprintf("%s-worker", i.ingestionType), i.processTask(workerErrs, threadComplete, task)).Start()
				task = nil
			} else {
				task, waitDuration, err = i.getTask(c)
				switch {
				case err != nil:
					c.Errorf("Error getting task for %s ingestor: %s. Will retry", i.ingestionType, err.Error())
					waitDuration = ptr.Duration(time.Duration(2 * time.Minute))
				case waitDuration != nil:
					c.Infof("Waiting")
				case task != nil:
					numWorkers++
					async.WithContext(workerErrs, fmt.Sprintf("%s-worker", i.ingestionType), i.processTask(workerErrs, threadComplete, task)).Start()
					task = nil
				default:
					// This should probably not happen
					c.Warnf("Got null task for ingestion type %s, but no wait time or error either", i.ingestionType)
					waitDuration = ptr.Duration(time.Duration(10 * time.Second))
				}
				if waitDuration != nil {
					timer = time.NewTimer(*waitDuration)
				}
			}
		}
	}
}

func (i *ingestor) getTask(c ctx.LogContext) (_task interface{}, _waitPeriod *time.Duration, _err error) {
	return nil, nil, nil
}

func (i *ingestor) processTask(workerErrs chan error, threadComplete chan interface{}, task interface{}) func(c async.Context) {
	return func(c async.Context) {
		var err error
		defer func() {
			if err != nil {
				workerErrs <- err
			} else {
				threadComplete <- task
			}
		}()
		switch i.ingestionType {
		case content.IngestStrategyWebsiteHTML1:
			link, ok := task.(links2.Link)
			if !ok {
				err = fmt.Errorf("Expected the task to be of type Link, but was not")
			} else {
				err = processWebsiteHTML1Link(c, link)
			}
		case content.IngestStrategyPodcastRSS1:
			sourceSeed, ok := task.(content.SourceSeed)
			if !ok {
				err = fmt.Errorf("Expected the task to be of type Source Seed, but was not")
			} else {
				err = processPodcastRSS1SourceSeed(c, sourceSeed)
			}
		default:
			err = fmt.Errorf("Ingest strategy %s is unsupported", i.ingestionType)
		}
	}
}
