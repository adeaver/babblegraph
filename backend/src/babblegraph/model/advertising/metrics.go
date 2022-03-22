package advertising

import (
	"babblegraph/model/users"
	"babblegraph/util/cache"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
)

const (
	getNumberOfSendsQuery        = "SELECT advertisement_id, user_id, COUNT(*) count FROM advertising_user_advertisements WHERE advertisement_id IN (?) GROUP BY advertisement_id, user_id"
	getNumberOfOpenedEmailsQuery = `
        SELECT
            a.advertisement_id AS advertisement_id, e.user_id AS user_id, COUNT(*) count
        FROM
            advertising_user_advertisements a
        JOIN
            email_records e
        ON
            e._id = a.email_record_id
        WHERE
            e.first_opened_at IS NOT NULL AND
            a.advertisement_id IN (?)
        GROUP BY a.advertisement_id, e.user_id`
	getNumberOfClicksQuery = `
        SELECT
            a.advertisement_id AS advertisement_id, c.user_advertisement_id AS user_advertisement_id, COUNT(*) count
        FROM
            advertising_user_advertisements a
        JOIN
            advertising_advertisement_clicks c
        ON
            a._id = c.user_advertisement_id
        WHERE
            a.advertisement_id IN (?)
        GROUP BY a.advertisement_id, c.user_advertisement_id`
)

type CampaignMetrics struct {
	AdvertisementMetrics []AdvertisementMetrics `json:"advertisement_metrics"`
	LastRefreshedAt      time.Time              `json:"last_refereshed_at"`
}

type AdvertisementMetrics struct {
	AdvertisementID      AdvertisementID `json:"advertisement_id"`
	NumberOfSends        MetricNumber    `json:"number_of_sends"`
	NumberOfOpenedEmails MetricNumber    `json:"number_of_opened_emails"`
	NumberOfClicks       MetricNumber    `json:"number_of_clicks"`
}

type MetricNumber struct {
	Total  int64 `json:"total"`
	Unique int64 `json:"unique"`
}

func GetAdvertisementMetrics(tx *sqlx.Tx, campaignID CampaignID) (*CampaignMetrics, error) {
	var campaignMetrics CampaignMetrics
	if err := cache.WithCache(fmt.Sprintf("advertising_metrics_campaign_id_%s", campaignID), &campaignMetrics, 1*time.Second, func() (interface{}, error) {
		advertisements, err := GetAllAdvertisementsForCampaignID(tx, campaignID)
		if err != nil {
			return nil, err
		}
		now := time.Now()
		var advertisementIDs []AdvertisementID
		for _, a := range advertisements {
			advertisementIDs = append(advertisementIDs, a.ID)
		}
		if len(advertisementIDs) == 0 {
			return CampaignMetrics{
				LastRefreshedAt: now,
			}, nil
		}
		sendMetrics, err := getNumberOfSendsMetrics(tx, advertisementIDs)
		if err != nil {
			return nil, err
		}
		clickMetrics, err := getNumberOfClicksMetrics(tx, advertisementIDs)
		if err != nil {
			return nil, err
		}
		emailMetrics, err := getNumberOfEmailsMetrics(tx, advertisementIDs)
		if err != nil {
			return nil, err
		}
		var advertisementMetrics []AdvertisementMetrics
		for _, id := range advertisementIDs {
			advertisementMetrics = append(advertisementMetrics, AdvertisementMetrics{
				AdvertisementID:      id,
				NumberOfSends:        sendMetrics[id],
				NumberOfOpenedEmails: emailMetrics[id],
				NumberOfClicks:       clickMetrics[id],
			})
		}
		return CampaignMetrics{
			LastRefreshedAt:      now,
			AdvertisementMetrics: advertisementMetrics,
		}, nil
	}); err != nil {
		return nil, err
	}
	return &campaignMetrics, nil
}

func getNumberOfSendsMetrics(tx *sqlx.Tx, advertisementIDs []AdvertisementID) (map[AdvertisementID]MetricNumber, error) {
	var matches []struct {
		AdvertisementID AdvertisementID `db:"advertisement_id"`
		UserID          users.UserID    `db:"user_id"`
		Count           int64           `db:"count"`
	}
	query, args, err := sqlx.In(getNumberOfSendsQuery, advertisementIDs)
	if err != nil {
		return nil, err
	}
	sql := tx.Rebind(query)
	if err := tx.Select(&matches, sql, args...); err != nil {
		return nil, err
	}
	out := make(map[AdvertisementID]MetricNumber)
	for _, m := range matches {
		currentMetricNumber := out[m.AdvertisementID]
		currentMetricNumber.Total += m.Count
		currentMetricNumber.Unique += 1
		out[m.AdvertisementID] = currentMetricNumber
	}
	return out, nil
}

func getNumberOfEmailsMetrics(tx *sqlx.Tx, advertisementIDs []AdvertisementID) (map[AdvertisementID]MetricNumber, error) {
	var matches []struct {
		AdvertisementID AdvertisementID `db:"advertisement_id"`
		UserID          users.UserID    `db:"user_id"`
		Count           int64           `db:"count"`
	}
	query, args, err := sqlx.In(getNumberOfOpenedEmailsQuery, advertisementIDs)
	if err != nil {
		return nil, err
	}
	sql := tx.Rebind(query)
	if err := tx.Select(&matches, sql, args...); err != nil {
		return nil, err
	}
	out := make(map[AdvertisementID]MetricNumber)
	for _, m := range matches {
		currentMetricNumber := out[m.AdvertisementID]
		currentMetricNumber.Total += m.Count
		currentMetricNumber.Unique += 1
		out[m.AdvertisementID] = currentMetricNumber
	}
	return out, nil
}

func getNumberOfClicksMetrics(tx *sqlx.Tx, advertisementIDs []AdvertisementID) (map[AdvertisementID]MetricNumber, error) {
	var matches []struct {
		AdvertisementID AdvertisementID     `db:"advertisement_id"`
		UserID          UserAdvertisementID `db:"user_advertisement_id"`
		Count           int64               `db:"count"`
	}
	query, args, err := sqlx.In(getNumberOfClicksQuery, advertisementIDs)
	if err != nil {
		return nil, err
	}
	sql := tx.Rebind(query)
	if err := tx.Select(&matches, sql, args...); err != nil {
		return nil, err
	}
	out := make(map[AdvertisementID]MetricNumber)
	for _, m := range matches {
		currentMetricNumber := out[m.AdvertisementID]
		currentMetricNumber.Total += m.Count
		currentMetricNumber.Unique += 1
		out[m.AdvertisementID] = currentMetricNumber
	}
	return out, nil
}
