package linkprocessing

import (
	"babblegraph/model/links2"
	"babblegraph/util/database"
	"babblegraph/util/urlparser"
	"sync"
	"time"

	"github.com/jmoiron/sqlx"
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
	var domains []string
	if err := database.WithTx(func(tx *sqlx.Tx) error {
		var err error
		domains, err = links2.GetDomainsWithUnfetchedLinks(tx)
		return err
	}); err != nil {
		return nil, err
	}
	var orderedDomains []Domain
	domainHash := make(map[string]bool)
	for _, d := range domains {
		domainHash[d] = true
		orderedDomains = append(orderedDomains, Domain{
			Domain: d,
			FreeAt: time.Now(),
		})
	}
	return &LinkProcessor{
		DomainSet:      domainHash,
		OrderedDomains: orderedDomains,
	}
}

func (l *LinkProcessor) GetLink() (*links2.Link, *time.Duration, error) {
	l.mu.Lock()
	defer l.mu.Unlock()
	firstDomain := l.OrderedDomains[0]
	if firstDomain.FreeAt.After(time.Now()) {
		waitTime := firstDomain.FreeAt.Sub(time.Now())
		return nil, &waitTime, nil
	}
	l.OrderedDomains = append(l.OrderedDomains[1:], Domain{
		Domain: firstDomain.Domain,
		FreeAt: time.Now().Add(15 * time.Second),
	})
	var link *links2.Link
	if err := database.WithTx(func(tx *sqlx.Tx) error {
		var err error
		link, err = links2.GetUnfetchedLinkForDomain(tx, firstDomain.Domain)
		if err != nil {
			return err
		}
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
		if p := urlparser.ParseURL(u); p != nil {
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
