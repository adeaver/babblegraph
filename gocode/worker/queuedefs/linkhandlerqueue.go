package queuedefs

import (
	"babblegraph/worker/linkhandler"
	"babblegraph/worker/links"
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/adeaver/babblegraph/lib/util/database"
	"github.com/adeaver/babblegraph/lib/util/queue"
	"github.com/jmoiron/sqlx"
)

const queueTopicNameLinkHandler queueTopicName = "link-handler-topic"

type linkHandlerQueue struct {
	mux            sync.Mutex
	orderedDomains []links.Domain
	linksByDomain  map[links.Domain][]links.Link
}

func (l *linkHandlerQueue) GetTopicName() string {
	return queueTopicNameLinkHandler.Str()
}

func initializeLinkHandlerQueue(errs chan error) (*linkHandlerQueue, error) {
	q := &linkHandlerQueue{}
	var unfetchedLinks map[links.Domain][]links.Link
	if err := database.WithTx(func(tx *sqlx.Tx) error {
		var err error
		unfetchedLinks, err = links.GetUnfetchedLinks(tx)
		return err
	}); err != nil {
		return nil, err
	}
	for d, _ := range unfetchedLinks {
		q.orderedDomains = append(q.orderedDomains, d)
	}
	q.linksByDomain = unfetchedLinks
	go startLinkHandler(q, errs)
	return q, nil
}

func startLinkHandler(q *linkHandlerQueue, errs chan error) {
	for {
		if err := q.sendNewLink(); err != nil {
			errs <- err
			return
		}
		time.Sleep(2200 * time.Millisecond)
	}
}

func (q *linkHandlerQueue) sendNewLink() error {
	q.mux.Lock()
	log.Println("unlocked to send links")
	defer q.mux.Unlock()
	if len(q.orderedDomains) == 0 {
		log.Println("no top domains")
		return nil
	}
	topDomain := q.orderedDomains[0]
	linksForDomain, ok := q.linksByDomain[topDomain]
	if !ok || len(linksForDomain) == 0 {
		return nil
	}
	topLink := linksForDomain[0]
	if err := publishMessageToFetchQueue(topLink.URL); err != nil {
		return err
	}
	log.Println(fmt.Sprintf("Sending url from domain %s: %s", topDomain.Str(), topLink.URL))
	if err := database.WithTx(func(tx *sqlx.Tx) error {
		return links.SetURLAsFetched(tx, topLink.URL)
	}); err != nil {
		return err
	}
	if len(linksForDomain) == 1 {
		q.orderedDomains = q.orderedDomains[1:]
		delete(q.linksByDomain, topDomain)
	} else {
		q.orderedDomains = append(q.orderedDomains[1:], topDomain)
		q.linksByDomain[topDomain] = linksForDomain[1:]
	}
	return nil
}

type linkHandlerQueueMessage struct {
	Links []string `json:"links"`
}

func (l *linkHandlerQueue) ProcessMessage(tx *sqlx.Tx, msg queue.Message) error {
	var m linkHandlerQueueMessage
	if err := json.Unmarshal([]byte(msg.MessageBody), &m); err != nil {
		log.Println(fmt.Sprintf("Error unmarshalling message for fetch queue: %s... marking complete", err.Error()))
		return nil
	}
	filteredLinks, err := linkhandler.FilterLinksAndInsert(tx, m.Links)
	if err != nil {
		return err
	}
	l.mux.Lock()
	defer l.mux.Unlock()
	for _, li := range filteredLinks {
		if linksForDomain, present := l.linksByDomain[li.Domain]; !present {
			l.linksByDomain[li.Domain] = []links.Link{li}
			l.orderedDomains = append(l.orderedDomains, li.Domain)
		} else {
			l.linksByDomain[li.Domain] = append(linksForDomain, li)
		}
	}
	return nil
}

func publishMessageToLinkHandlerQueue(urls []string) error {
	return queue.PublishMessageToQueueByName(queueTopicNameLinkHandler.Str(), linkHandlerQueueMessage{
		Links: urls,
	})
}
