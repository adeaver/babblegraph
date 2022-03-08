package newsletter

import (
	"babblegraph/model/documents"
	"babblegraph/model/domains"
	"babblegraph/model/email"
	"babblegraph/model/routes"
	"babblegraph/model/useraccounts"
	"babblegraph/util/ctx"
	"babblegraph/util/deref"
	"babblegraph/util/ptr"
	"fmt"
	"math/rand"
	"time"
)

const (
	DefaultNumberOfArticlesPerEmail = 12
	defaultNumberOfTopicsPerEmail   = 4

	minimumDaysSinceLastSpotlight = 3
)

type CreateNewsletterInput struct {
	WordsmithAccessor wordsmithAccessor
	EmailAccessor     emailAccessor
	UserAccessor      userPreferencesAccessor
	DocsAccessor      documentAccessor
	PodcastAcccessor  podcastAccessor
}

func CreateNewsletter(c ctx.LogContext, input CreateNewsletterInput) (*Newsletter, error) {
	emailRecordID := email.NewEmailRecordID()
	if err := input.EmailAccessor.InsertEmailRecord(emailRecordID, input.UserAccessor.getUserID()); err != nil {
		return nil, err
	}
	var numberOfDocumentsInNewsletter *int
	userSubscriptionLevel := input.UserAccessor.getUserSubscriptionLevel()
	switch {
	case userSubscriptionLevel == nil:
		// no-op
	case *userSubscriptionLevel == useraccounts.SubscriptionLevelBetaPremium,
		*userSubscriptionLevel == useraccounts.SubscriptionLevelPremium:
		if !input.UserAccessor.getUserNewsletterSchedule().IsSendRequested() {
			return nil, nil
		}
		numberOfDocumentsInNewsletter = ptr.Int(input.UserAccessor.getUserNewsletterSchedule().GetNumberOfDocuments())

	default:
		return nil, fmt.Errorf("Unrecognized subscription level: %s", *userSubscriptionLevel)
	}
	categories, err := getDocumentCategories(c, getDocumentCategoriesInput{
		emailRecordID:                 emailRecordID,
		languageCode:                  input.UserAccessor.getLanguageCode(),
		userAccessor:                  input.UserAccessor,
		docsAccessor:                  input.DocsAccessor,
		numberOfDocumentsInNewsletter: numberOfDocumentsInNewsletter,
	})
	if err != nil {
		return nil, err
	}
	var setTopicsLink *string
	switch {
	case len(input.UserAccessor.getUserTopics()) > 0:
		// no-op
	case input.UserAccessor.getDoesUserHaveAccount():
		setTopicsLink = ptr.String(routes.MakeLoginLinkWithContentTopicsRedirect())
	default:
		setTopicsLink, err = routes.MakeSetTopicsLink(input.UserAccessor.getUserID())
		if err != nil {
			return nil, err
		}
	}
	var reinforcementLink *string
	if input.UserAccessor.getDoesUserHaveAccount() {
		reinforcementLink = ptr.String(routes.MakeLoginLinkWithReinforcementRedirect())
	} else {
		reinforcementLink, err = routes.MakeWordReinforcementLink(input.UserAccessor.getUserID())
		if err != nil {
			return nil, err
		}
	}
	spotlightRecord, err := getSpotlightLemmaForNewsletter(c, getSpotlightLemmaForNewsletterInput{
		emailRecordID:     emailRecordID,
		categories:        categories,
		userAccessor:      input.UserAccessor,
		docsAccessor:      input.DocsAccessor,
		wordsmithAccessor: input.WordsmithAccessor,
	})
	if err != nil {
		return nil, err
	}
	return &Newsletter{
		UserID:        input.UserAccessor.getUserID(),
		EmailRecordID: emailRecordID,
		LanguageCode:  input.UserAccessor.getLanguageCode(),
		Body: NewsletterBody{
			LemmaReinforcementSpotlight: spotlightRecord,
			Categories:                  categories,
			SetTopicsLink:               setTopicsLink,
			ReinforcementLink:           *reinforcementLink,
		},
	}, nil
}

func pickUpToNRandomIndices(listLength, pickN int) []int {
	var availableIndices, out []int
	for i := 0; i < listLength; i++ {
		availableIndices = append(availableIndices, i)
	}
	if listLength <= pickN {
		return availableIndices
	}
	generator := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < pickN; i++ {
		idx := generator.Intn(int(len(availableIndices)))
		out = append(out, availableIndices[idx])
		availableIndices = append(availableIndices[:idx], availableIndices[idx+1:]...)
	}
	return out
}

func makeLinkFromDocument(c ctx.LogContext, emailRecordID email.ID, userAccessor userPreferencesAccessor, doc documents.Document) (*Link, error) {
	var title, imageURL, description *string
	if isNotEmpty(doc.Metadata.Title) {
		title = doc.Metadata.Title
	}
	if isNotEmpty(doc.Metadata.Image) {
		imageURL = doc.Metadata.Image
	}
	if isNotEmpty(doc.Metadata.Description) {
		description = doc.Metadata.Description
	}
	// TODO: make content accessor
	domain, err := domains.GetDomainMetadata(doc.Domain)
	if err != nil {
		c.Errorf("Error getting domain: %s", err.Error())
		return nil, nil
	}
	userDocumentID, err := userAccessor.insertDocumentForUserAndReturnID(emailRecordID, doc)
	if err != nil {
		return nil, err
	}
	articleLink, err := routes.MakeArticleLink(*userDocumentID)
	if err != nil {
		c.Errorf("Error making article link: %s", err.Error())
		return nil, nil
	}
	paywallReportLink, err := routes.MakePaywallReportLink(*userDocumentID)
	if err != nil {
		c.Errorf("Error making paywall report link: %s", err.Error())
		return nil, nil
	}
	return &Link{
		DocumentID:       doc.ID,
		URL:              *articleLink,
		PaywallReportURL: *paywallReportLink,
		ImageURL:         imageURL,
		Title:            title,
		Description:      description,
		Domain: &Domain{
			Name:      string(domain.Domain),
			FlagAsset: routes.GetFlagAssetForCountryCode(domain.Country),
		},
	}, nil
}

func isNotEmpty(s *string) bool {
	return len(deref.String(s, "")) > 0
}
