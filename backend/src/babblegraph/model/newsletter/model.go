package newsletter

import (
	"babblegraph/model/users"
	"babblegraph/wordsmith"
)

type Newsletter struct {
	UserID       users.UserID           `json:"user_id"`
	LanguageCode wordsmith.LanguageCode `json:"language_code"`
	Body         NewsletterBody         `json:"body"`
}

type NewsletterBody struct {
	LemmaReinforcementSpotlight *LemmaReinforcementSpotlight `json:"lemma_reinforcement_spotlight,omitempty"`
	Categories                  []Category                   `json:"category"`
	SetTopicsLink               *string                      `json:"set_topics_link,omitempty"`
	ReinforcementLink           string                       `json:"reinforcement_link"`
}

type LemmaReinforcementSpotlight struct {
	LemmaText       string `json:"lemma_text"`
	Document        Link   `json:"document"`
	PreferencesLink string `json:"preferences_link"`
}

type Category struct {
	Name  string `json:"name"`
	Links []Link `json:"links"`
}

type Link struct {
	ImageURL         *string `json:"image_url,omitempty"`
	Title            *string `json:"title,omitempty"`
	Description      *string `json:"description,omitempty"`
	URL              string  `json:"url"`
	PaywallReportURL string  `json:"paywall_report_url"`
	Domain           *Domain `json:"domain"`
}

type Domain struct {
	FlagAsset  string `json:"flag_asset"`
	DomainName string `json:"domain_name"`
}
