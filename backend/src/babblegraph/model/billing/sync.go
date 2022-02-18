package billing

import "github.com/jmoiron/sqlx"

const (

	// TODO: add hold_until to these queries if need be.
	getAllPremiumNewsletterSyncRequestQuery = "SELECT * FROM billing_premium_newsletter_sync_request"
	// This model might need to change, but the current idea is that this should act as more of queue.
	// Where each subscription can only have whatever the latest update type is. As of right now, the only update type is that it makes a switch to active
	insertPremiumNewsletterSyncRequestQuery    = "INSERT INTO billing_premium_newsletter_sync_request VALUES (premium_newsletter_subscription_id, update_type) VALUES ($1, $2) ON CONFLICT DO NOTHING"
	deletePremiumNewsletterSyncRequestQuery    = "DELETE FROM billing_premium_newsletter_sync_request WHERE premium_newsletter_subscription_id = $1"
	incrementPremiumNewsletterSyncRequestQuery = "UPDATE billing_premium_newsletter_sync_request SET attempt_number = attempt_number + 1 WHERE premium_newsletter_subscription_id = $1"
)

func InsertPremiumNewsletterSyncRequest(tx *sqlx.Tx, id PremiumNewsletterSubscriptionID, updateType PremiumNewsletterSubscriptionUpdateType) error {
	if _, err := tx.Exec(insertPremiumNewsletterSyncRequestQuery, id, updateType); err != nil {
		return err
	}
	return nil
}

func GetPremiumNewsletterSyncRequests(tx *sqlx.Tx) (map[PremiumNewsletterSubscriptionID]PremiumNewsletterSubscriptionUpdateType, error) {
	var matches []dbPremiumNewsletterSubscriptionSyncRequest
	if err := tx.Select(&matches, getAllPremiumNewsletterSyncRequestQuery); err != nil {
		return nil, err
	}
	out := make(map[PremiumNewsletterSubscriptionID]PremiumNewsletterSubscriptionUpdateType)
	for _, m := range matches {
		out[m.PremiumNewsletterSubscriptionID] = m.UpdateType
	}
	return out, nil
}

func MarkPremiumNewsletterSyncRequestDone(tx *sqlx.Tx, id PremiumNewsletterSubscriptionID) error {
	if _, err := tx.Exec(deletePremiumNewsletterSyncRequestQuery, id); err != nil {
		return err
	}
	return nil
}

func MarkPremiumNewsletterSyncRequestForRetry(tx *sqlx.Tx, id PremiumNewsletterSubscriptionID) error {
	if _, err := tx.Exec(incrementPremiumNewsletterSyncRequestQuery, id); err != nil {
		return err
	}
	return nil
}
