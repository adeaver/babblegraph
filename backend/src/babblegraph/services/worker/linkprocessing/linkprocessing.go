package linkprocessing

import (
	"babblegraph/model/content"
	"babblegraph/model/contenttopics"
	"babblegraph/model/documents"
	"babblegraph/model/domains"
	"babblegraph/model/links2"
	"babblegraph/model/urltopicmapping"
	"babblegraph/services/worker/indexing"
	"babblegraph/services/worker/ingesthtml"
	"babblegraph/services/worker/textprocessing"
	"babblegraph/util/async"
	"babblegraph/util/bufferedfetch"
	"babblegraph/util/ctx"
	"babblegraph/util/database"
	"babblegraph/util/opengraph"
	"babblegraph/util/ptr"
	"babblegraph/util/urlparser"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/jmoiron/sqlx"
)

const (
	defaultChunkSize     = 500
	defaultTimeUntilFree = 5 * time.Second
	defaultRefreshPeriod = 1 * time.Hour
)

type Domain struct {
	Domain string
	FreeAt time.Time
}

type LinkProcessor struct {
	mu             sync.Mutex
	OrderedDomains []Domain
	DomainSet      map[string]bool
}

func CreateLinkProcessor() (*LinkProcessor, error) {
	domains := domains.GetDomains()
	var orderedDomains []Domain
	domainHash := make(map[string]bool)
	for _, d := range domains {
		domainHash[d] = true
		orderedDomains = append(orderedDomains, Domain{
			Domain: d,
			FreeAt: time.Now(),
		})
		if err := bufferedfetch.Register(getBufferedFetchKeyForDomain(d), makeBufferedFetchForDomain(d)); err != nil {
			return nil, err
		}
	}
	return &LinkProcessor{
		DomainSet:      domainHash,
		OrderedDomains: orderedDomains,
	}, nil
}

func (l *LinkProcessor) resetDomains(c ctx.LogContext) error {
	domains := domains.GetDomains()
	var orderedDomains []Domain
	domainHash := make(map[string]bool)
	for _, d := range domains {
		domainHash[d] = true
		orderedDomains = append(orderedDomains, Domain{
			Domain: d,
			FreeAt: time.Now(),
		})
		if err := bufferedfetch.ForceRefill(c, getBufferedFetchKeyForDomain(d)); err != nil {
			return err
		}
	}
	l.DomainSet = domainHash
	l.OrderedDomains = orderedDomains
	return nil
}

func makeBufferedFetchForDomain(domain string) func() (interface{}, error) {
	return func() (interface{}, error) {
		var links []links2.Link
		log.Println(fmt.Sprintf("Fetching links for domain: %s", domain))
		if err := database.WithTx(func(tx *sqlx.Tx) error {
			var err error
			links, err = links2.LookupBulkUnfetchedLinksForDomain(tx, domain, defaultChunkSize)
			return err
		}); err != nil {
			return nil, err
		}
		return links, nil
	}
}

func getBufferedFetchKeyForDomain(domain string) string {
	return fmt.Sprintf("linkprocessing-%s", domain)
}

func (l *LinkProcessor) ProcessLinks(maxWorkers int) func(c async.Context) {
	return func(c async.Context) {
		c.Infof("Starting link processor")
		addURLs := make(chan []string)
		workerManagerErrs := make(chan error)
		urlManagerErrs := make(chan error)
		timer := time.NewTimer(defaultRefreshPeriod)
		async.WithContext(workerManagerErrs, "worker-manager", l.startWorkerManager(maxWorkers, addURLs)).Start()
		async.WithContext(urlManagerErrs, "url-manager", l.startURLManager(addURLs)).Start()
		for {
			select {
			case err := <-workerManagerErrs:
				c.Warnf("Error with worker manager %s", err.Error())
				async.WithContext(workerManagerErrs, "worker-manager", l.startWorkerManager(maxWorkers, addURLs)).Start()
			case err := <-urlManagerErrs:
				c.Warnf("Error with URL manager %s", err.Error())
				async.WithContext(urlManagerErrs, "url-manager", l.startURLManager(addURLs)).Start()
			case _ = <-timer.C:
				c.Infof("Refreshing link processor")
				if err := l.resetDomains(c); err != nil {
					c.Errorf("Error resetting domains: %s", err.Error())
				}
				timer = time.NewTimer(defaultRefreshPeriod)
			}
		}
	}
}

func (l *LinkProcessor) startWorkerManager(maxWorkers int, addURLs chan []string) func(c async.Context) {
	return func(c async.Context) {
		threadComplete := make(chan *links2.Link, 1)
		workerErrs := make(chan error)
		timer := time.NewTimer(1 * time.Second)
		var link *links2.Link
		var numWorkers int
		spinOffWorkerOrWait := func() (_shouldBreak bool) {
			link, duration, err := l.getLink(c)
			switch {
			case err != nil:
				c.Errorf("Error getting links: %s, will retry", err.Error())
				timer = time.NewTimer(2 * time.Minute)
				return true
			case duration != nil:
				c.Infof("Waiting")
				timer = time.NewTimer(*duration)
				return true
			case link != nil:
				numWorkers++
				async.WithContext(workerErrs, "ingest-worker", processSingleLink(threadComplete, addURLs, link)).Start()
			}
			return false
		}
		// Start initial workers
		for i := 0; i < maxWorkers; i++ {
			if shouldBreak := spinOffWorkerOrWait(); shouldBreak {
				break
			}
		}
		// At this point, there should either be:
		// maxWorkers # of workers, in which case we wait for errors or thread completes
		// or a timer, in which case, we wait for a timer
		var duration *time.Duration
		link, _, err := l.getLink(c)
		if err != nil {
			c.Errorf("Error getting links: %s", err.Error())
		}
		for {
			select {
			case _ = <-threadComplete:
				c.Infof("Thread is complete")
				numWorkers--
				if link != nil {
					async.WithContext(workerErrs, "ingest-worker", processSingleLink(threadComplete, addURLs, link)).Start()
					link = nil
				} else {
					link, duration, err = l.getLink(c)
					switch {
					case err != nil:
						c.Errorf("Error getting links: %s", err.Error())
						timer = time.NewTimer(2 * time.Minute)
					case duration != nil:
						timer = time.NewTimer(*duration)
					case link != nil:
						async.WithContext(workerErrs, "ingest-worker", processSingleLink(threadComplete, addURLs, link)).Start()
						link = nil
					default:
						panic("unreachable")
					}
				}
			case err := <-workerErrs:
				c.Infof("Thread encountered error %s", err.Error())
				numWorkers--
				if link != nil {
					async.WithContext(workerErrs, "ingest-worker", processSingleLink(threadComplete, addURLs, link)).Start()
					link = nil
				} else {
					link, duration, err = l.getLink(c)
					switch {
					case err != nil:
						c.Errorf("Error getting links: %s", err.Error())
						timer = time.NewTimer(2 * time.Minute)
					case duration != nil:
						timer = time.NewTimer(*duration)
					case link != nil:
						async.WithContext(workerErrs, "ingest-worker", processSingleLink(threadComplete, addURLs, link)).Start()
						link = nil
					default:
						panic("unreachable")
					}
				}
			case _ = <-timer.C:
				c.Infof("Worker manager timer has finished. Currently there are %d workers", numWorkers)
				switch {
				case numWorkers == maxWorkers && link != nil:
					c.Infof("All workers are busy, and link is non-nil, continuing...")
				case numWorkers == maxWorkers && link == nil:
					c.Infof("All workers are busy, but link needs replenshing")
					link, duration, err = l.getLink(c)
					if err != nil {
						c.Errorf("Error getting link: %s", err.Error())
						timer = time.NewTimer(2 * time.Minute)
					}
				case numWorkers < maxWorkers && link == nil:
					link, duration, err = l.getLink(c)
					switch {
					case err != nil:
						c.Errorf("Error getting links: %s", err.Error())
						timer = time.NewTimer(2 * time.Minute)
					case duration != nil:
						timer = time.NewTimer(*duration)
					case link != nil:
						async.WithContext(workerErrs, "ingest-worker", processSingleLink(threadComplete, addURLs, link)).Start()
					default:
						panic("unreachable")
					}
				case numWorkers < maxWorkers && link != nil:
					async.WithContext(workerErrs, "ingest-worker", processSingleLink(threadComplete, addURLs, link)).Start()
					link = nil
				}
			}
		}
	}
}

func (l *LinkProcessor) startURLManager(addURLs chan []string) func(c async.Context) {
	handleURLs := func(c async.Context, urls []string) {
		domainSet := make(map[string]bool)
		var parsedURLs []urlparser.ParsedURL
		var contentTopics [][]contenttopics.ContentTopic
		for _, u := range urls {
			if parsedURL := urlparser.ParseURL(u); parsedURL != nil && domains.IsURLAllowed(*parsedURL) {
				domainSet[parsedURL.Domain] = true
				domainMetadata, err := domains.GetDomainMetadata(parsedURL.Domain)
				if err != nil {
					c.Warnf("Got error getting metadata for domain %s on url %s: %s. Continuing...", parsedURL.Domain, u, err.Error())
					continue
				}
				parsedURLs = append(parsedURLs, *parsedURL)
				contentTopics = append(contentTopics, domainMetadata.Topics)
			}
		}
		if len(parsedURLs) == 0 {
			return
		}
		c.Infof("Acquiring lock for URLs")
		l.mu.Lock()
		defer l.mu.Unlock()
		defer func() {
			c.Infof("Releasing lock for URLs")
		}()
		for domain := range domainSet {
			if _, ok := l.DomainSet[domain]; !ok {
				l.DomainSet[domain] = true
				l.OrderedDomains = append(l.OrderedDomains, Domain{
					Domain: domain,
					FreeAt: time.Now().Add(defaultTimeUntilFree),
				})
			}
		}
		if err := database.WithTx(func(tx *sqlx.Tx) error {
			if err := links2.InsertLinks(tx, parsedURLs); err != nil {
				return err
			}
			return nil
		}); err != nil {
			c.Warnf("Error saving URLs: %s", err.Error())
		}
	}
	return func(c async.Context) {
		for {
			select {
			case urls := <-addURLs:
				handleURLs(c, urls)
			}
		}
	}
}

func processSingleLink(threadComplete chan *links2.Link, addURLs chan []string, link *links2.Link) func(c async.Context) {
	return func(c async.Context) {
		c.Infof("Starting new thread")
		u := link.URL
		domain := link.Domain
		if p := urlparser.ParseURL(u); p != nil && domains.IsSeedURL(*p) {
			c.Debugf("Received url %s, which is a seed url. Skipping...", u)
			threadComplete <- nil
			return
		}
		c.Infof("Processing URL %s with identifier %s", u, link.URLIdentifier)
		parsedHTMLPage, err := ingesthtml.ProcessURLDEPRECATED(u, domain)
		if err != nil {
			c.Infof("Got error ingesting html for url %s: %s. Continuing...", u, err.Error())
			threadComplete <- link
			return
		}
		domainMetadata, err := domains.GetDomainMetadata(domain)
		if err != nil {
			c.Warnf("Got error getting metadata for domain %s on url %s: %s. Continuing...", domain, u, err.Error())
			threadComplete <- link
			return
		}
		languageCode := domainMetadata.LanguageCode
		c.Infof("Attempting to add URLs")
		addURLs <- parsedHTMLPage.Links
		c.Infof("Successfully added URLs")
		c.Debugf("Processing text for url %s", u)
		var description *string
		if d, ok := parsedHTMLPage.Metadata[opengraph.DescriptionTag.Str()]; ok {
			description = ptr.String(d)
		}
		textMetadata, err := textprocessing.ProcessText(textprocessing.ProcessTextInput{
			BodyText:     parsedHTMLPage.BodyText,
			Description:  description,
			LanguageCode: languageCode,
		})
		if err != nil {
			c.Warnf("Got error processing text for url %s: %s. Continuing...", u, err.Error())
			threadComplete <- link
			return
		}
		var sourceID *content.SourceID
		var topicsForURL []contenttopics.ContentTopic
		var topicMappingIDs []content.TopicMappingID
		var topicIDs []content.TopicID
		if err := database.WithTx(func(tx *sqlx.Tx) error {
			var err error
			sourceID, err = link.GetSourceID(tx)
			topicsForURL, topicMappingIDs, err = urltopicmapping.GetTopicsAndMappingIDsForURL(tx, u)
			if err != nil {
				return err
			}
			var sourceSeedTopicMappings []content.SourceSeedTopicMappingID
			for _, topicMappingID := range topicMappingIDs {
				sourceSeedTopicMapping, sourceTopicMapping, err := topicMappingID.GetOriginID()
				switch {
				case err != nil:
					return err
				case sourceTopicMapping != nil:
					c.Warnf("Found source topic mapping from ID %s, which is unsupported", topicMappingID)
				case sourceSeedTopicMapping != nil:
					sourceSeedTopicMappings = append(sourceSeedTopicMappings, *sourceSeedTopicMapping)
				default:
					return fmt.Errorf("unreachable")
				}
			}
			if len(sourceSeedTopicMappings) > 0 {
				topicIDs, err = content.LookupTopicsForSourceSeedMappingIDs(tx, sourceSeedTopicMappings)
				if err != nil {
					return err
				}
			}
			return err
		}); err != nil {
			c.Warnf("Error getting topics for url %s: %s. Continuing...", u, err.Error())
			threadComplete <- link
			return
		}
		c.Debugf("Indexing text for URL %s", u)
		err = indexing.IndexDocument(c, indexing.IndexDocumentInput{
			ParsedHTMLPage:         *parsedHTMLPage,
			TextMetadata:           *textMetadata,
			LanguageCode:           languageCode,
			DocumentVersion:        documents.CurrentDocumentVersion,
			URL:                    urlparser.MustParseURL(u),
			SourceID:               sourceID,
			TopicsForURL:           topicsForURL,
			TopicIDs:               topicIDs,
			TopicMappingIDs:        topicMappingIDs,
			SeedJobIngestTimestamp: link.SeedJobIngestTimestamp,
		})
		if err != nil {
			c.Warnf("Got error indexing document for url %s: %s. Continuing...", u, err.Error())
			threadComplete <- link
			return
		}
		c.Infof("Thread is complete, sending terminate request")
		threadComplete <- nil
		c.Infof("Thread is exiting")
	}
}

func (l *LinkProcessor) getLink(c ctx.LogContext) (*links2.Link, *time.Duration, error) {
	c.Infof("Trying to get lock for links")
	l.mu.Lock()
	c.Infof("Acquired lock")
	var firstNonEmptyDomainIdx *int
	defer func() {
		if firstNonEmptyDomainIdx != nil {
			c.Infof("Removing all domains up to %d that are empty", *firstNonEmptyDomainIdx)
			l.OrderedDomains = append([]Domain{}, l.OrderedDomains[*firstNonEmptyDomainIdx:]...)
			c.Infof("Available domains are now of length %d", len(l.OrderedDomains))
		}
		c.Infof("Releasing lock")
		l.mu.Unlock()
	}()
	for idx, domain := range l.OrderedDomains {
		if domain.FreeAt.After(time.Now()) {
			c.Infof("Domain %s is not free, sending wait", domain.Domain)
			waitTime := domain.FreeAt.Sub(time.Now())
			return nil, &waitTime, nil
		}
		var link *links2.Link
		err := bufferedfetch.WithNextBufferedValue(getBufferedFetchKeyForDomain(domain.Domain), func(i interface{}) error {
			l, ok := i.(links2.Link)
			if !ok {
				return fmt.Errorf("error getting next value for domain %s: incorrect type in buffered fetch", domain.Domain)
			}
			link = &l
			return nil
		})
		switch {
		case err != nil:
			return nil, nil, err
		case link == nil:
			c.Infof("Domain %s has no links, skipping", domain.Domain)
			firstNonEmptyDomainIdx = ptr.Int(idx + 1)
		case link != nil:
			firstNonEmptyDomainIdx = ptr.Int(idx + 1)
			l.OrderedDomains = append(l.OrderedDomains, Domain{
				Domain: domain.Domain,
				FreeAt: time.Now().Add(defaultTimeUntilFree),
			})
			if err := database.WithTx(func(tx *sqlx.Tx) error {
				return links2.SetURLAsFetched(tx, link.URLIdentifier)
			}); err != nil {
				return nil, nil, err
			}
			return link, nil, nil
		}
	}
	c.Infof("No links available, sending wait")
	return nil, ptr.Duration(defaultRefreshPeriod), nil
}

func (l *LinkProcessor) AddURLs(urls []string, topics []contenttopics.ContentTopic) error {
	var parsedURLs []urlparser.ParsedURL
	// using a hash set because the domains are
	// likely to repeat
	domainSet := make(map[string]bool)
	for _, u := range urls {
		if p := urlparser.ParseURL(u); p != nil && domains.IsURLAllowed(*p) {
			domainSet[p.Domain] = true
			parsedURLs = append(parsedURLs, *p)
		}
	}
	l.mu.Lock()
	defer l.mu.Unlock()
	for domain := range domainSet {
		if _, ok := l.DomainSet[domain]; !ok {
			l.DomainSet[domain] = true
			l.OrderedDomains = append(l.OrderedDomains, Domain{
				Domain: domain,
				FreeAt: time.Now().Add(defaultTimeUntilFree),
			})
		}
	}
	if len(parsedURLs) == 0 {
		return nil
	}
	return database.WithTx(func(tx *sqlx.Tx) error {
		if err := links2.InsertLinks(tx, parsedURLs); err != nil {
			return err
		}
		return nil
	})
}

func (l *LinkProcessor) ReseedDomains() {
	l.mu.Lock()
	defer l.mu.Unlock()
	domains := domains.GetDomains()
	var orderedDomains []Domain
	domainHash := make(map[string]bool)
	for _, d := range domains {
		domainHash[d] = true
		orderedDomains = append(orderedDomains, Domain{
			Domain: d,
			FreeAt: time.Now(),
		})
	}
	l.DomainSet = domainHash
	l.OrderedDomains = orderedDomains
}
