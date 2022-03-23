package newsletter

import (
	"babblegraph/model/content"
	"babblegraph/model/documents"
	"babblegraph/model/email"
	"babblegraph/model/users"
	"babblegraph/wordsmith"
)

type Newsletter struct {
	UserID        users.UserID           `json:"user_id"`
	EmailRecordID email.ID               `json:"email_record_id"`
	LanguageCode  wordsmith.LanguageCode `json:"language_code"`
	Body          NewsletterBody         `json:"body"`
}

type NewsletterBody struct {
	LemmaReinforcementSpotlight *LemmaReinforcementSpotlight `json:"lemma_reinforcement_spotlight,omitempty"`
	Categories                  []Category                   `json:"categories"`
	SetTopicsLink               *string                      `json:"set_topics_link,omitempty"`
	ReinforcementLink           string                       `json:"reinforcement_link"`
	PreferencesLink             *string                      `json:"preferences_link,omitempty"`
	Advertisement               *NewsletterAdvertisement     `json:"advertisement,omitempty"`
}

type LemmaReinforcementSpotlight struct {
	LemmaText       string `json:"lemma_text"`
	Document        Link   `json:"document"`
	PreferencesLink string `json:"preferences_link"`
}

type Category struct {
	// This is not needed on deserialization
	topicID *content.TopicID

	Name         *string       `json:"name,omitempty"`
	Links        []Link        `json:"links"`
	PodcastLinks []PodcastLink `json:"podcast_links,omitempty"`
}

type PodcastLink struct {
	PodcastName        string  `json:"podcast_name"`
	WebsiteURL         string  `json:"website_url"`
	PodcastImageURL    *string `json:"podcast_image_url"`
	EpisodeTitle       string  `json:"episode_title"`
	EpisodeDescription string  `json:"episode_description"`
	ListenURL          string  `json:"listen_url"`
}

type Link struct {
	DocumentID       documents.DocumentID `json:"document_id"`
	ImageURL         *string              `json:"image_url,omitempty"`
	Title            *string              `json:"title,omitempty"`
	Description      *string              `json:"description,omitempty"`
	URL              string               `json:"url"`
	PaywallReportURL string               `json:"paywall_report_url"`
	Domain           *Domain              `json:"domain"`
}

type Domain struct {
	FlagAsset string `json:"flag_asset"`
	Name      string `json:"name"`
}

type NewsletterAdvertisement struct {
	Title                       string                       `json:"title"`
	Description                 string                       `json:"description"`
	ImageURL                    string                       `json:"image_url"`
	URL                         string                       `json:"url"`
	AdditionalAdvertisementLink *AdditionalAdvertisementLink `json:"additional_link,omitempty"`
	PremiumLink                 string                       `json:"premium_link"`
	AdvertisementPolicyLink     string                       `json:"advertisement_policy_link"`
}

type AdditionalAdvertisementLink struct {
	URL      string  `json:"url"`
	LinkText *string `json:"link_text"`
}
