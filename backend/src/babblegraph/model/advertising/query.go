package advertising

import (
	"babblegraph/model/content"
	"babblegraph/model/email"
	"babblegraph/model/users"
	"babblegraph/util/ctx"
	"babblegraph/util/math/int2"
	"babblegraph/wordsmith"
	"math/rand"
	"time"

	"github.com/jmoiron/sqlx"
)

const (
	// Roughly three days between ads
	minimumTimeBetweenAds = 3*24*time.Hour - 3*time.Hour

	// Roughly one month
	minimumTimeBetweenAdsSameCampaign = 30 * 24 * time.Hour

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

	getAllCampaignIDsWithForTopicMappingQuery = `
        SELECT
            _id
        FROM
            advertising_campaigns
        WHERE
            _id IN
            (
                SELECT
                    campaign_id
                FROM
                    advertising_campaign_topic_mappings
                WHERE
                    is_active = TRUE AND
                    topic_id = $1
            )
        `
	getAllCampaignIDsEligibleForAllUsersQuery = `
        SELECT
            _id
        FROM
            advertising_campaigns
        WHERE
            should_apply_to_all_users = TRUE AND
            is_active = TRUE
    `
	lookupAdvertisementQuery = `
        SELECT
            *
        FROM
            advertising_advertisements
        WHERE
            is_active=TRUE AND
            campaign_id IN (?)
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
	allIneligibleCampaigns, err := getFullListOfIneligibleCampaignIDs(c, tx, topic, ineligibleCampaignIDs)
	if err != nil {
		return nil, err
	}
	c.Debugf("Got %d ineligible campaign ids", len(allIneligibleCampaigns))
	var matches []struct {
		CampaignID CampaignID `db:"_id"`
	}
	switch {
	case topic != nil:
		if err := tx.Select(&matches, getAllCampaignIDsWithForTopicMappingQuery, *topic); err != nil {
			return nil, err
		}
	case topic == nil:
		if err := tx.Select(&matches, getAllCampaignIDsEligibleForAllUsersQuery); err != nil {
			return nil, err
		}
	default:
		panic("unreachable")
	}
	if len(matches) == 0 {
		c.Debugf("No potentially valid advertisements found")
		return nil, nil
	}
	var validCampaignIDs []CampaignID
	for _, m := range matches {
		if _, ok := allIneligibleCampaigns[m.CampaignID]; ok {
			c.Debugf("Campaign ID %s is ineligible for user %s", m.CampaignID, userID)
			continue
		}
		validCampaignIDs = append(validCampaignIDs, m.CampaignID)
	}
	if len(validCampaignIDs) == 0 {
		c.Debugf("Filtered all valid advertisements")
		return nil, nil
	}
	var userAdvertisements []dbUserAdvertisement
	if err := tx.Select(&userAdvertisements, getUserAdvertisementsByDateQuery, userID); err != nil {
		return nil, err
	}
	query, args, err := sqlx.In(lookupAdvertisementQuery, validCampaignIDs)
	if err != nil {
		return nil, err
	}
	sql := tx.Rebind(query)
	var advertisementMatches []dbAdvertisement
	if err := tx.Select(&advertisementMatches, sql, args...); err != nil {
		return nil, err
	}
	return joinPossibleAdvertisements(c, userAdvertisements, advertisementMatches)
}

const defaultWeight int = 10

func joinPossibleAdvertisements(c ctx.LogContext, userAdvertisements []dbUserAdvertisement, advertisementMatches []dbAdvertisement) (*Advertisement, error) {
	// Current idea of heuristic:
	// All campaigns that don't have an advertisement get a weight of 10
	// Campaigns are then weighted as follows: min(10, weeks since ad last ran - 4)
	if len(advertisementMatches) == 0 {
		return nil, nil
	}
	weightsByCampaign := calculateWeightsForSeenAds(userAdvertisements)
	adsByCampaignID := make(map[CampaignID][]dbAdvertisement)
	for _, match := range advertisementMatches {
		adsByCampaignID[match.CampaignID] = append(adsByCampaignID[match.CampaignID], match)
		if _, ok := weightsByCampaign[match.CampaignID]; !ok {
			weightsByCampaign[match.CampaignID] = defaultWeight
		}
	}
	weightedCampaignIDs := getWeightedCampaignIDs(weightsByCampaign)
	rand.Seed(time.Now().UnixNano())
	idx := rand.Intn(len(weightedCampaignIDs))
	campaignID := weightedCampaignIDs[idx]
	advertisements, ok := adsByCampaignID[campaignID]
	if !ok {
		c.Warnf("Campaign ID %s selected, but there were no matching ads", campaignID)
		return nil, nil
	}
	adIdx := rand.Intn(len(advertisements))
	ad := advertisements[adIdx].ToNonDB()
	return &ad, nil
}

func calculateWeightsForSeenAds(userAdvertisements []dbUserAdvertisement) map[CampaignID]int {
	weightsByCampaign := make(map[CampaignID]int)
	for _, userAdvertisement := range userAdvertisements {
		weight, ok := weightsByCampaign[userAdvertisement.CampaignID]
		if !ok {
			weight = defaultWeight
		}
		weeksSinceAdAppeared := int(time.Now().Sub(userAdvertisement.CreatedAt) / (7 * 24 * time.Hour))
		weightsByCampaign[userAdvertisement.CampaignID] = int2.MustMinInt(weight, int2.MustMaxInt(0, weeksSinceAdAppeared-4))
	}
	return weightsByCampaign
}

func getWeightedCampaignIDs(weightsByCampaign map[CampaignID]int) []CampaignID {
	var weightedCampaignIDs []CampaignID
	for campaignID, weight := range weightsByCampaign {
		for i := 0; i < weight; i++ {
			weightedCampaignIDs = append(weightedCampaignIDs, campaignID)
		}
	}
	return weightedCampaignIDs
}

func getFullListOfIneligibleCampaignIDs(c ctx.LogContext, tx *sqlx.Tx, topic *content.TopicID, ineligibleCampaignIDs []CampaignID) (map[CampaignID]bool, error) {
	type queryMatch struct {
		CampaignID CampaignID `db:"_id"`
	}
	ineligibleCampaignIDHashSet := make(map[CampaignID]bool)
	// Collect all the campaigns for which this user is not eligible
	for _, campaignID := range ineligibleCampaignIDs {
		c.Debugf("Campaign ID %s was ineligible for user", campaignID)
		ineligibleCampaignIDHashSet[campaignID] = true
	}
	// Collect all campaigns that are inactive or have inactive vendors/sources
	var inactiveCampaigns []queryMatch
	if err := tx.Select(&inactiveCampaigns, getAllInactiveCampaignIDsQuery); err != nil {
		return nil, err
	}
	for _, campaign := range inactiveCampaigns {
		c.Debugf("Campaign %s is inactive", campaign.CampaignID)
		ineligibleCampaignIDHashSet[campaign.CampaignID] = true
	}
	return ineligibleCampaignIDHashSet, nil
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
