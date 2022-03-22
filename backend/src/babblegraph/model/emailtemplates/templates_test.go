package emailtemplates

import (
	"babblegraph/model/documents"
	"babblegraph/model/email"
	"babblegraph/model/newsletter"
	"babblegraph/model/users"
	"babblegraph/util/ptr"
	"strings"
	"testing"
)

func TestCreateNewsletterTemplate(t *testing.T) {
	userAccessor := &testUserAccessor{
		userHasAccount: false,
		userID:         users.UserID("12345"),
	}
	testNewsletter := newsletter.NewsletterBody{
		SetTopicsLink:     ptr.String("babblegraph.com/topics"),
		ReinforcementLink: "babblegraph.com/reinforce",
	}
	html, err := MakeNewsletterHTML(MakeNewsletterHTMLInput{
		EmailRecordID: email.NewEmailRecordID(),
		UserAccessor:  userAccessor,
		Body:          testNewsletter,
	})
	if err != nil {
		t.Fatalf(err.Error())
	}
	if len(*html) == 0 {
		t.Fatalf("Got empty body")
	}
	if !strings.Contains(*html, "babblegraph.com/topics") {
		t.Errorf("Expected topics link")
	}
	if !strings.Contains(*html, "babblegraph.com/reinforce") {
		t.Errorf("Expected reinforcement link")
	}
}

func TestCreateNewsletterWithLinksTemplate(t *testing.T) {
	userAccessor := &testUserAccessor{
		userHasAccount: false,
		userID:         users.UserID("12345"),
	}
	testNewsletter := newsletter.NewsletterBody{
		SetTopicsLink:     ptr.String("babblegraph.com/topics"),
		ReinforcementLink: "babblegraph.com/reinforce",
		PreferencesLink:   ptr.String("babblegraph.com/preferences"),
		Categories: []newsletter.Category{
			{
				Name: ptr.String("Test Category"),
				Links: []newsletter.Link{
					{
						DocumentID:       documents.DocumentID("test"),
						ImageURL:         ptr.String("babblegraph.com"),
						Title:            ptr.String("Test Link"),
						Description:      ptr.String("Test Description"),
						URL:              "babblegraph.com",
						PaywallReportURL: "babblegraph.com",
						Domain: &newsletter.Domain{
							FlagAsset: "babblegraph.com",
							Name:      "Babblegraph",
						},
					},
				},
				PodcastLinks: []newsletter.PodcastLink{
					{
						PodcastName:        "Test",
						WebsiteURL:         "babblegraph.com",
						PodcastImageURL:    ptr.String("babblegraph.com"),
						EpisodeTitle:       "Test Episode",
						EpisodeDescription: "This episode is a test",
						ListenURL:          "babblegraph.com",
					},
				},
			},
		},
	}
	html, err := MakeNewsletterHTML(MakeNewsletterHTMLInput{
		EmailRecordID: email.NewEmailRecordID(),
		UserAccessor:  userAccessor,
		Body:          testNewsletter,
	})
	if err != nil {
		t.Fatalf(err.Error())
	}
	if len(*html) == 0 {
		t.Fatalf("Got empty body")
	}
	if !strings.Contains(*html, "babblegraph.com/topics") {
		t.Errorf("Expected topics link")
	}
	if !strings.Contains(*html, "babblegraph.com/reinforce") {
		t.Errorf("Expected reinforcement link")
	}
	if !strings.Contains(*html, "Test Link") {
		t.Errorf("Expected title")
	}
	if !strings.Contains(*html, "Test Description") {
		t.Errorf("Expected description")
	}
	if !strings.Contains(*html, "Test Episode") {
		t.Errorf("Expected episode title")
	}
	if !strings.Contains(*html, "This episode is a test") {
		t.Errorf("Expected episode description")
	}
}

func TestCreateNewsletterNoPodcastLinksTemplate(t *testing.T) {
	userAccessor := &testUserAccessor{
		userHasAccount: false,
		userID:         users.UserID("12345"),
	}
	testNewsletter := newsletter.NewsletterBody{
		SetTopicsLink:     ptr.String("babblegraph.com/topics"),
		ReinforcementLink: "babblegraph.com/reinforce",
		PreferencesLink:   ptr.String("babblegraph.com/preferences"),
		Categories: []newsletter.Category{
			{
				Name: ptr.String("Test Category"),
				Links: []newsletter.Link{
					{
						DocumentID:       documents.DocumentID("test"),
						ImageURL:         ptr.String("babblegraph.com"),
						Title:            ptr.String("Test Link"),
						Description:      ptr.String("Test Description"),
						URL:              "babblegraph.com",
						PaywallReportURL: "babblegraph.com",
						Domain: &newsletter.Domain{
							FlagAsset: "babblegraph.com",
							Name:      "Babblegraph",
						},
					},
				},
				PodcastLinks: []newsletter.PodcastLink{},
			},
		},
	}
	html, err := MakeNewsletterHTML(MakeNewsletterHTMLInput{
		EmailRecordID: email.NewEmailRecordID(),
		UserAccessor:  userAccessor,
		Body:          testNewsletter,
	})
	if err != nil {
		t.Fatalf(err.Error())
	}
	if len(*html) == 0 {
		t.Fatalf("Got empty body")
	}
	if !strings.Contains(*html, "babblegraph.com/topics") {
		t.Errorf("Expected topics link")
	}
	if !strings.Contains(*html, "babblegraph.com/reinforce") {
		t.Errorf("Expected reinforcement link")
	}
	if !strings.Contains(*html, "Test Link") {
		t.Errorf("Expected title")
	}
	if !strings.Contains(*html, "Test Description") {
		t.Errorf("Expected description")
	}
}

func TestCreateNewsletterWithAdsTemplate(t *testing.T) {
	userAccessor := &testUserAccessor{
		userHasAccount: false,
		userID:         users.UserID("12345"),
	}
	testNewsletter := newsletter.NewsletterBody{
		SetTopicsLink:     ptr.String("babblegraph.com/topics"),
		ReinforcementLink: "babblegraph.com/reinforce",
		PreferencesLink:   ptr.String("babblegraph.com/preferences"),
		Categories: []newsletter.Category{
			{
				Name: ptr.String("Test Category"),
				Links: []newsletter.Link{
					{
						DocumentID:       documents.DocumentID("test"),
						ImageURL:         ptr.String("babblegraph.com"),
						Title:            ptr.String("Test Link"),
						Description:      ptr.String("Test Description"),
						URL:              "babblegraph.com",
						PaywallReportURL: "babblegraph.com",
						Domain: &newsletter.Domain{
							FlagAsset: "babblegraph.com",
							Name:      "Babblegraph",
						},
					},
				},
				PodcastLinks: []newsletter.PodcastLink{},
			},
		},
		Advertisement: &newsletter.NewsletterAdvertisement{
			Title:                   "This is a title",
			URL:                     "babblegraph.com/this-is-a-test-advertisement",
			ImageURL:                "babblegraph.com/advertisement.jpeg",
			Description:             "This is an ad",
			PremiumLink:             "babblegraph.com/premium-link",
			AdvertisementPolicyLink: "babblegraph.com/advertising-policy-link",
		},
	}
	html, err := MakeNewsletterHTML(MakeNewsletterHTMLInput{
		EmailRecordID: email.NewEmailRecordID(),
		UserAccessor:  userAccessor,
		Body:          testNewsletter,
	})
	if err != nil {
		t.Fatalf(err.Error())
	}
	if len(*html) == 0 {
		t.Fatalf("Got empty body")
	}
	if !strings.Contains(*html, "babblegraph.com/topics") {
		t.Errorf("Expected topics link")
	}
	if !strings.Contains(*html, "babblegraph.com/reinforce") {
		t.Errorf("Expected reinforcement link")
	}
	if !strings.Contains(*html, "babblegraph.com/this-is-a-test-advertisement") {
		t.Errorf("Expected advertising link")
	}
	if !strings.Contains(*html, "This is an ad") {
		t.Errorf("Expected description")
	}
}
