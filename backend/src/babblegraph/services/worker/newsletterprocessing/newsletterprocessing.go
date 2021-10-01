package newsletterprocessing

import (
	"babblegraph/model/newslettersendrequests"
	"babblegraph/model/users"
	"babblegraph/util/database"
	"babblegraph/wordsmith"
	"fmt"
	"log"
	"sort"
	"sync"
	"time"

	"github.com/jmoiron/sqlx"
)

const syncInterval = 3 * time.Hour

type NewsletterProcessor struct {
	mu                           sync.Mutex
	timeOfLastSync               time.Time
	orderedSendRequestsToPreload []newslettersendrequests.NewsletterSendRequest
	orderedSendRequestsToFulfill []newslettersendrequests.NewsletterSendRequest
}

func CreateNewsletterProcessor() (*NewsletterProcessor, error) {
	n := &NewsletterProcessor{}
	n.mu.Lock()
	defer n.mu.Unlock()
	n.syncSendRequests()
	return n, nil
}

func (n *NewsletterProcessor) GetNextSendRequestToPreload() (*newslettersendrequests.NewsletterSendRequest, error) {
	n.mu.Lock()
	defer n.mu.Unlock()
	if time.Now().After(n.timeOfLastSync.Add(syncInterval)) {
		if err := n.syncSendRequests(); err != nil {
			return nil, err
		}
	}
	if len(n.orderedSendRequestsToPreload) == 0 {
		return nil, nil
	}
	nextSendRequestToPreload := n.orderedSendRequestsToPreload[0]
	n.orderedSendRequestsToPreload = append([]newslettersendrequests.NewsletterSendRequest{}, n.orderedSendRequestsToPreload[1:]...)
	return &nextSendRequestToPreload, nil
}

func (n *NewsletterProcessor) GetNextSendRequestToFulfill() (*newslettersendrequests.NewsletterSendRequest, error) {
	n.mu.Lock()
	defer n.mu.Unlock()
	if time.Now().After(n.timeOfLastSync.Add(syncInterval)) {
		if err := n.syncSendRequests(); err != nil {
			return nil, err
		}
	}
	if len(n.orderedSendRequestsToFulfill) == 0 {
		return nil, nil
	}
	nextSendRequestToFulfill := n.orderedSendRequestsToFulfill[0]
	if time.Now().Before(nextSendRequestToFulfill.DateOfSend) {
		return nil, nil
	}
	n.orderedSendRequestsToFulfill = append([]newslettersendrequests.NewsletterSendRequest{}, n.orderedSendRequestsToFulfill[1:]...)
	return &nextSendRequestToFulfill, nil
}

// Acquire lock before calling this function
func (n *NewsletterProcessor) syncSendRequests() error {
	log.Println(fmt.Sprintf("Syncing send requests"))
	toPreload, toFulfill, err := getSendRequestsByType()
	if err != nil {
		return err
	}
	log.Println(fmt.Sprintf("Got %d send requests to preload and %d to fulfill", len(toPreload), len(toFulfill)))
	n.orderedSendRequestsToPreload = toPreload
	n.orderedSendRequestsToFulfill = toFulfill
	sort.Slice(n.orderedSendRequestsToFulfill, func(i, j int) bool {
		return n.orderedSendRequestsToFulfill[i].DateOfSend.Before(n.orderedSendRequestsToFulfill[j].DateOfSend)
	})
	n.timeOfLastSync = time.Now()
	log.Println(fmt.Sprintf("Sync complete"))
	return nil
}

func getSendRequestsByType() (_toPreload, _toFulfill []newslettersendrequests.NewsletterSendRequest, _err error) {
	var toPreload, toFulfill []newslettersendrequests.NewsletterSendRequest
	today := time.Now()
	sendRequestsForToday, err := getSendRequestsForDay(today)
	if err != nil {
		return nil, nil, err
	}
	for _, s := range sendRequestsForToday {
		switch s.PayloadStatus {
		case newslettersendrequests.PayloadStatusNeedsPreload:
			toPreload = append(toPreload, s)
		case newslettersendrequests.PayloadStatusPayloadReady:
			toFulfill = append(toFulfill, s)
		case newslettersendrequests.PayloadStatusNoSendRequested,
			newslettersendrequests.PayloadStatusSent,
			newslettersendrequests.PayloadStatusDeleted:
			// no-op
		default:
			log.Println(fmt.Sprintf("Unrecgonized payload status: %s", s.PayloadStatus))
		}
	}
	tomorrow := today.Add(24 * time.Hour)
	sendRequestsForTomorrow, err := getSendRequestsForDay(tomorrow)
	if err != nil {
		return nil, nil, err
	}
	for _, s := range sendRequestsForTomorrow {
		switch s.PayloadStatus {
		case newslettersendrequests.PayloadStatusNeedsPreload:
			toPreload = append(toPreload, s)
		case newslettersendrequests.PayloadStatusPayloadReady,
			newslettersendrequests.PayloadStatusNoSendRequested,
			newslettersendrequests.PayloadStatusSent,
			newslettersendrequests.PayloadStatusDeleted:
			// no-op
		default:
			log.Println(fmt.Sprintf("Unrecgonized payload status: %s", s.PayloadStatus))
		}
	}
	return toPreload, toFulfill, nil
}

func getSendRequestsForDay(t time.Time) ([]newslettersendrequests.NewsletterSendRequest, error) {
	var sendRequestsForDay []newslettersendrequests.NewsletterSendRequest
	if err := database.WithTx(func(tx *sqlx.Tx) error {
		activeUsers, err := users.GetAllActiveUsers(tx)
		if err != nil {
			return err
		}
		var userIDs []users.UserID
		for _, u := range activeUsers {
			userIDs = append(userIDs, u.ID)
		}
		sendRequestsForDay, err = newslettersendrequests.GetOrCreateSendRequestsForUsersForDay(tx, userIDs, wordsmith.LanguageCodeSpanish, t)
		return err
	}); err != nil {
		return nil, err
	}
	return sendRequestsForDay, nil
}
