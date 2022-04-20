package newsletterprocessing

import (
	"babblegraph/config"
	"babblegraph/model/newslettersendrequests"
	"babblegraph/model/useraccounts"
	"babblegraph/model/users"
	"babblegraph/util/ctx"
	"babblegraph/util/database"
	"babblegraph/wordsmith"
	"sort"
	"sync"
	"time"

	"github.com/jmoiron/sqlx"
)

type NewsletterProcessor struct {
	mu                           sync.Mutex
	timeOfLastSync               time.Time
	orderedSendRequestsToPreload []newslettersendrequests.NewsletterSendRequest
	orderedSendRequestsToFulfill []newslettersendrequests.NewsletterSendRequest
}

func CreateNewsletterProcessor(c ctx.LogContext) (*NewsletterProcessor, error) {
	n := &NewsletterProcessor{}
	n.mu.Lock()
	defer n.mu.Unlock()
	n.syncSendRequests(c)
	return n, nil
}

func (n *NewsletterProcessor) GetNextSendRequestToPreload(c ctx.LogContext) (*newslettersendrequests.NewsletterSendRequest, error) {
	n.mu.Lock()
	defer n.mu.Unlock()
	if time.Now().After(n.timeOfLastSync.Add(config.NewsletterSendRequestSyncInterval)) {
		if err := n.syncSendRequests(c); err != nil {
			return nil, err
		}
	}
	if len(n.orderedSendRequestsToPreload) == 0 {
		return nil, nil
	}
	nextSendRequestToPreload := n.orderedSendRequestsToPreload[0]
	if time.Now().Before(nextSendRequestToPreload.DateOfSend.Add(-24 * time.Hour)) {
		return nil, nil
	}
	n.orderedSendRequestsToPreload = append([]newslettersendrequests.NewsletterSendRequest{}, n.orderedSendRequestsToPreload[1:]...)
	return &nextSendRequestToPreload, nil
}

func (n *NewsletterProcessor) GetNextSendRequestToFulfill(c ctx.LogContext) (*newslettersendrequests.NewsletterSendRequest, error) {
	n.mu.Lock()
	defer n.mu.Unlock()
	if time.Now().After(n.timeOfLastSync.Add(config.NewsletterSendRequestSyncInterval)) {
		if err := n.syncSendRequests(c); err != nil {
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
func (n *NewsletterProcessor) syncSendRequests(c ctx.LogContext) error {
	c.Infof("Syncing send requests")
	toPreload, toFulfill, err := getSendRequestsByType(c)
	if err != nil {
		return err
	}
	c.Infof("Got %d send requests to preload and %d to fulfill", len(toPreload), len(toFulfill))
	n.orderedSendRequestsToPreload = toPreload
	sort.Slice(n.orderedSendRequestsToPreload, func(i, j int) bool {
		return n.orderedSendRequestsToPreload[i].DateOfSend.Before(n.orderedSendRequestsToPreload[j].DateOfSend)
	})
	n.orderedSendRequestsToFulfill = toFulfill
	sort.Slice(n.orderedSendRequestsToFulfill, func(i, j int) bool {
		return n.orderedSendRequestsToFulfill[i].DateOfSend.Before(n.orderedSendRequestsToFulfill[j].DateOfSend)
	})
	n.timeOfLastSync = time.Now()
	c.Infof("Sync complete")
	return nil
}

func getSendRequestsByType(c ctx.LogContext) (_toPreload, _toFulfill []newslettersendrequests.NewsletterSendRequest, _err error) {
	var toPreload, toFulfill []newslettersendrequests.NewsletterSendRequest
	today := time.Now()
	sendRequestsForToday, err := getSendRequestsForDay(c, today)
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
			c.Infof("Unrecgonized payload status: %s", s.PayloadStatus)
		}
	}
	tomorrow := today.Add(24 * time.Hour)
	sendRequestsForTomorrow, err := getSendRequestsForDay(c, tomorrow)
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
			c.Infof("Unrecgonized payload status: %s", s.PayloadStatus)
		}
	}
	return toPreload, toFulfill, nil
}

func getSendRequestsForDay(c ctx.LogContext, t time.Time) ([]newslettersendrequests.NewsletterSendRequest, error) {
	var sendRequestsForDay []newslettersendrequests.NewsletterSendRequest
	if err := database.WithTx(func(tx *sqlx.Tx) error {
		activeUsers, err := users.GetAllActiveUsers(tx)
		if err != nil {
			return err
		}
		var userIDs []users.UserID
		for _, u := range activeUsers {
			userSubscription, err := useraccounts.LookupSubscriptionLevelForUser(tx, u.ID)
			switch {
			case err != nil:
				return err
			case userSubscription == nil:
				// no-op
			default:
				userIDs = append(userIDs, u.ID)
			}
		}
		sendRequestsForDay, err = newslettersendrequests.GetOrCreateSendRequestsForUsersForDay(c, tx, userIDs, wordsmith.LanguageCodeSpanish, t)
		return err
	}); err != nil {
		return nil, err
	}
	return sendRequestsForDay, nil
}
