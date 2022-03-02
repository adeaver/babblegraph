package content

import (
	"babblegraph/util/ctx"
	"babblegraph/util/ptr"
	"babblegraph/util/urlparser"
	"sync"
	"time"

	"github.com/jmoiron/sqlx"
)

const defaultSyncPeriod = 2 * time.Hour

type sourceUnion struct {
	source     Source
	sourceSeed SourceSeed
}

var (
	allowableSourceMap map[SourceID]Source
	seedURLHashSet     map[string]sourceUnion

	initializerMutex                sync.Mutex
	timeSinceLastInitializationSync *time.Time
)

func initializeInMemoryCache(tx *sqlx.Tx) error {
	initializerMutex.Lock()
	defer initializerMutex.Unlock()
	if len(seedURLHashSet) != 0 && len(allowableSourceMap) != 0 {
		return nil
	}
	sources, err := GetAllSources(tx)
	if err != nil {
		return err
	}
	for _, s := range sources {
		if s.IsActive {
			allowableSourceMap[s.ID] = s
		}
		seedURLHashSet[s.URL] = sourceUnion{
			source: s,
		}
		sourceSeeds, err := GetAllSourceSeedsForSource(tx, s.ID)
		if err != nil {
			return err
		}
		for _, seed := range sourceSeeds {
			seedURLHashSet[seed.URL] = sourceUnion{
				sourceSeed: seed,
			}
		}
	}
	timeSinceLastInitializationSync = ptr.Time(time.Now())
	return nil
}

func needsResync() bool {
	return len(seedURLHashSet) == 0 || len(allowableSourceMap) == 0 || timeSinceLastInitializationSync == nil || time.Now().After(timeSinceLastInitializationSync.Add(defaultSyncPeriod))
}

func IsParsedURLASeedURL(tx *sqlx.Tx, u urlparser.ParsedURL) (bool, error) {
	if needsResync() {
		if err := initializeInMemoryCache(tx); err != nil {
			return false, err
		}
	}
	_, ok := seedURLHashSet[u.URL]
	return ok, nil
}

func LookupActiveSourceForSourceID(c ctx.LogContext, tx *sqlx.Tx, sourceID SourceID) (*Source, error) {
	if needsResync() {
		if err := initializeInMemoryCache(tx); err != nil {
			return nil, err
		}
	}
	source, ok := allowableSourceMap[sourceID]
	if ok {
		return &source, nil
	}
	c.Infof("Source %s not found, querying db", sourceID)
	lookupSource, err := GetSource(tx, sourceID)
	switch {
	case err != nil:
		return nil, err
	case !lookupSource.IsActive:
		return nil, nil
	default:
		return lookupSource, nil
	}
}
