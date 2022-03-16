package advertising

import (
	"babblegraph/model/content"
	"babblegraph/model/email"
	"babblegraph/model/users"
	"babblegraph/util/ctx"
	"babblegraph/wordsmith"
	"time"

	"github.com/jmoiron/sqlx"
)

const (
	// Roughly one week between ads
	minimumTimeBetweenAds = 7*24*time.Hour - 3*time.Hour

	// Roughly three months
	minimumTimeBetweenAdsSameCampaign = 90 * 24 * time.Hour

	getUserAdvertisementsByDateQuery = "SELECT * FROM advertising_user_advertisements WHERE user_id = $1 ORDER BY created_at DESC"
	insertUserAdvertisementQuery     = "INSERT INTO advertising_user_advertisements (user_id, advertisement_id, campaign_id, email_record_id) VALUES ($1, $2, $3, $4) RETURNING _id"

	getAllInactiveCampaignIDsQuery = `
        SELECT
            _id
        FROM
            advertising_campaigns
        WHERE
            is_active = FALSE OR
            vendor_id IN (
                SELECT
                    _id
                FROM
                    advertising_vendors
                WHERE
                    is_active = FALSE
            ) OR
            source_id IN (
                SELECT
                    _id
                FROM
                    advertising_sources
                WHERE
                    is_active = FALSE
            )
        `
	getAllCampaignIDsWithNoTopicMappingQuery = `
        SELECT
            DISTINCT(campaign_id) AS _id
        FROM
            advertising_campaign_topic_mappings
        WHERE
            (
                is_active = TRUE AND
                topic_id = $1
            ) = FALSE
        `
	getAllCampaignIDsNotEligibleForAllUsersQuery = `
        SELECT
            _id
        FROM
            advertising_campaigns
        WHERE
            should_apply_to_all_users = FALSE AND
            is_active = TRUE
    `
	lookupAdvertisementQuery = `
        SELECT
            *
        FROM
            advertising_advertisements
        WHERE
            is_active=TRUE AND
            campaign_id NOT IN (?)
    `
)

type UserAdvertisementEligibility struct {
	IsUserEligibleForAdvertisement bool
	IneligibleCampaignIDs          []CampaignID
}

func GetUserAdvertisementEligibility(tx *sqlx.Tx, userID users.UserID) (*UserAdvertisementEligibility, error) {
	var matches []dbUserAdvertisement
	if err := tx.Select(&matches, getUserAdvertisementsByDateQuery, userID); err != nil {
		return nil, err
	}
	isEligible, ineligibleCampaigns := determineEligiblityFromUserAdvertisements(matches)
	return &UserAdvertisementEligibility{
		IsUserEligibleForAdvertisement: isEligible,
		IneligibleCampaignIDs:          ineligibleCampaigns,
	}, nil
}

// Separate function for testing
func determineEligiblityFromUserAdvertisements(matches []dbUserAdvertisement) (_isEligible bool, _ineligibleCampaigns []CampaignID) {
	ineligibleCampaignIDs := make(map[CampaignID]bool)
	for _, m := range matches {
		if time.Now().Before(m.CreatedAt.Add(minimumTimeBetweenAds)) {
			// User is ineligible for any campaign
			return false, nil
		}
		if time.Now().Before(m.CreatedAt.Add(minimumTimeBetweenAdsSameCampaign)) {
			ineligibleCampaignIDs[m.CampaignID] = true
		}
	}
	var out []CampaignID
	for campaignID := range ineligibleCampaignIDs {
		out = append(out, campaignID)
	}
	return true, out
}

func QueryAdvertisementsForUser(c ctx.LogContext, tx *sqlx.Tx, userID users.UserID, topic *content.TopicID, languageCode wordsmith.LanguageCode, ineligibleCampaignIDs []CampaignID) (*Advertisement, error) {
	allIneligibleCampaignIDs, err := getFullListOfIneligibleCampaignIDs(c, tx, topic, ineligibleCampaignIDs)
	if err != nil {
		return nil, err
	}
	c.Debugf("Got %d ineligible campaign ids", len(allIneligibleCampaignIDs))
	query, args, err := sqlx.In(lookupAdvertisementQuery, allIneligibleCampaignIDs)
	if err != nil {
		return nil, err
	}
	sql := tx.Rebind(query)
	var matches []dbAdvertisement
	if err := tx.Select(&matches, sql, args...); err != nil {
		return nil, err
	}
	for _, m := range matches {
		if m.LanguageCode == languageCode {
			// TODO(here): use experiment package to determine eligbility
			out := m.ToNonDB()
			return &out, nil
		}
	}
	c.Debugf("Did not find suitable ad in %+v", matches)
	return nil, nil
}

func getFullListOfIneligibleCampaignIDs(c ctx.LogContext, tx *sqlx.Tx, topic *content.TopicID, ineligibleCampaignIDs []CampaignID) ([]CampaignID, error) {
	type queryMatch struct {
		CampaignID CampaignID `db:"_id"`
	}
	ineligibleCampaignIDHashSet := make(map[CampaignID]bool)
	// Collect all the campaigns for which this user is not eligible
	for _, c := range ineligibleCampaignIDs {
		ineligibleCampaignIDHashSet[c] = true
	}
	// Collect all campaigns that are inactive or have inactive vendors/sources
	var inactiveCampaigns []queryMatch
	if err := tx.Select(&inactiveCampaigns, getAllInactiveCampaignIDsQuery); err != nil {
		return nil, err
	}
	for _, c := range inactiveCampaigns {
		ineligibleCampaignIDHashSet[c.CampaignID] = true
	}
	if topic != nil {
		// Collect all campaigns that don't apply to
		var noTopicCampaigns []queryMatch
		if err := tx.Select(&noTopicCampaigns, getAllCampaignIDsWithNoTopicMappingQuery, *topic); err != nil {
			return nil, err
		}
		for _, c := range noTopicCampaigns {
			ineligibleCampaignIDHashSet[c.CampaignID] = true
		}
	} else {
		// Collect all campaigns that don't apply to all users
		var applyToAllCampaigns []queryMatch
		if err := tx.Select(&applyToAllCampaigns, getAllCampaignIDsNotEligibleForAllUsersQuery); err != nil {
			return nil, err
		}
		for _, c := range applyToAllCampaigns {
			ineligibleCampaignIDHashSet[c.CampaignID] = true
		}
	}
	var out []CampaignID
	for campaignID := range ineligibleCampaignIDHashSet {
		out = append(out, campaignID)
	}
	return out, nil
}

func InsertUserAdvertisementAndGetID(tx *sqlx.Tx, userID users.UserID, ad Advertisement, emailRecordID email.ID) (*UserAdvertisementID, error) {
	rows, err := tx.Queryx(insertUserAdvertisementQuery, userID, ad.ID, ad.CampaignID, emailRecordID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var userAdvertisementID UserAdvertisementID
	for rows.Next() {
		if err := rows.Scan(&userAdvertisementID); err != nil {
			return nil, err
		}
	}
	return &userAdvertisementID, nil
}
