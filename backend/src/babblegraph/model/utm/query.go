package utm

import (
	"babblegraph/util/ptr"
	"babblegraph/util/random"
	"net/http"

	"github.com/jmoiron/sqlx"
)

type Parameters struct {
	Source     *Source
	Medium     *Medium
	CampaignID *CampaignID
	URLPath    *string
}

const (
	sourceKey     = "utm_source"
	campaignIDKey = "utm_campaign"
	mediumKey     = "utm_medium"

	trackingIDLength = 40

	UTMTrackingIDCookieName = "uttrid"
)

func GetParametersForRequest(r *http.Request) *Parameters {
	var source, medium, campaignID *string
	urlParams := r.URL.Query()
	if s := urlParams.Get(sourceKey); len(s) > 0 {
		source = ptr.String(s)
	}
	if c := urlParams.Get(campaignIDKey); len(c) > 0 {
		campaignID = ptr.String(c)
	}
	if m := urlParams.Get(mediumKey); len(m) > 0 {
		medium = ptr.String(m)
	}
	if source == nil && medium == nil && campaignID == nil {
		return nil
	}
	return &Parameters{
		Source:     sourceFromStrOrNil(source),
		CampaignID: campaignIDFromStrOrNil(campaignID),
		Medium:     mediumFromStrOrNil(medium),
		URLPath:    ptr.String(r.URL.Path),
	}
}

func GetTrackingIDForRequest(r *http.Request) (*TrackingID, error) {
	for _, c := range r.Cookies() {
		if c.Name == UTMTrackingIDCookieName {
			trackingID := c.Value
			return TrackingID(trackingID).Ptr(), nil
		}
	}
	newTrackingID, err := random.MakeRandomString(trackingIDLength)
	if err != nil {
		return nil, err
	}
	return TrackingID(*newTrackingID).Ptr(), nil
}

const insertUTMPageHitQuery = "INSERT INTO utm_page_hits (source, medium, campaign_id, url_path, tracking_id) VALUES ($1, $2, $3, $4, $5)"

func RegisterUTMPageHit(tx *sqlx.Tx, trackingID TrackingID, params Parameters) error {
	if _, err := tx.Exec(insertUTMPageHitQuery, params.Source, params.Medium, params.CampaignID, params.URLPath, trackingID); err != nil {
		return err
	}
	return nil
}

const registerEventQuery = "INSERT INTO utm_events (event_type, tracking_id) VALUES ($1, $2)"

func RegisterEvent(tx *sqlx.Tx, eventName string, trackingID string) error {
	if _, err := tx.Exec(registerEventQuery, eventName, trackingID); err != nil {
		return err
	}
	return nil
}
