package linkprocessing

import (
	"babblegraph/model/links2"
	"babblegraph/services/worker/domains"
	"babblegraph/util/bufferedfetch"
	"babblegraph/util/database"
	"babblegraph/util/urlparser"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/jmoiron/sqlx"
)

const defaultChunkSize = 2000

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
	domains := domains.GetSeedURLs()
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

func (l *LinkProcessor) GetLink() (*links2.Link, *time.Duration, error) {
	l.mu.Lock()
	defer l.mu.Unlock()
	firstDomain := l.OrderedDomains[0]
	if firstDomain.FreeAt.After(time.Now()) {
		waitTime := firstDomain.FreeAt.Sub(time.Now())
		return nil, &waitTime, nil
	}
	var shouldKeepDomain bool
	defer func() {
		if !shouldKeepDomain {
			delete(l.DomainSet, firstDomain.Domain)
		}
	}()
	l.OrderedDomains = append([]Domain{}, l.OrderedDomains[1:]...)
	var link *links2.Link
	if err := bufferedfetch.WithNextBufferedValue(getBufferedFetchKeyForDomain(firstDomain.Domain), func(i interface{}) error {
		l, ok := i.(links2.Link)
		if !ok {
			return fmt.Errorf("error getting next value for domain %s: incorrect type in buffered fetch", firstDomain.Domain)
		}
		link = &l
		return nil
	}); err != nil {
		return nil, nil, err
	}
	if link == nil {
		return nil, nil, nil
	}
	shouldKeepDomain = true
	l.OrderedDomains = append(l.OrderedDomains, Domain{
		Domain: firstDomain.Domain,
		FreeAt: time.Now().Add(15 * time.Second),
	})
	if err := database.WithTx(func(tx *sqlx.Tx) error {
		return links2.SetURLAsFetched(tx, link.URLIdentifier)
	}); err != nil {
		return nil, nil, err
	}
	return link, nil, nil
}

func (l *LinkProcessor) AddURLs(urls []string) error {
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
	for domain, _ := range domainSet {
		if _, ok := l.DomainSet[domain]; !ok {
			l.DomainSet[domain] = true
			l.OrderedDomains = append(l.OrderedDomains, Domain{
				Domain: domain,
				FreeAt: time.Now().Add(15 * time.Second),
			})
		}
	}
	return database.WithTx(func(tx *sqlx.Tx) error {
		return links2.InsertLinks(tx, parsedURLs)
	})
}
