package utm

import "time"

type dbUTMPageHit struct {
	ID         ID          `db:"_id"`
	Source     *Source     `db:"source"`
	Medium     *Medium     `db:"medium"`
	CampaignID *CampaignID `db:"campaign_id"`
	URLPath    *string     `db:"url_path"`
	TrackingID TrackingID  `db:"tracking_id"`
	CreatedAt  time.Time   `db:"created_at"`
}

type ID string

type Source string

const (
	SourceGoogle Source = "google"
)

func (s Source) Ptr() *Source {
	return &s
}

func sourceFromStrOrNil(s *string) *Source {
	if s == nil {
		return nil
	}
	return Source(*s).Ptr()
}

type Medium string

const (
	MediumClick Medium = "click"
)

func (m Medium) Ptr() *Medium {
	return &m
}

func mediumFromStrOrNil(s *string) *Medium {
	if s == nil {
		return nil
	}
	return Medium(*s).Ptr()
}

type CampaignID string

const (
	CampaignIDGoogleAds1 CampaignID = "googleads1"
)

func (c CampaignID) Ptr() *CampaignID {
	return &c
}

func campaignIDFromStrOrNil(s *string) *CampaignID {
	if s == nil {
		return nil
	}
	return CampaignID(*s).Ptr()
}

type TrackingID string

func (t TrackingID) Ptr() *TrackingID {
	return &t
}

func (t *TrackingID) Str() string {
	return string(*t)
}

type Event struct {
	ID         ID         `db:"_id"`
	CreatedAt  time.Time  `db:"created_at"`
	TrackingID TrackingID `db:"tracking_id"`
	EventType  EventType  `db:"event_type"`
}

type EventType string

const (
	EventTypeSignup EventType = "signup"
)

func (e EventType) Str() string {
	return string(e)
}
