package newsletter

import (
	"babblegraph/model/documents"
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
	ContentAccessor   contentAccessor
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

type makeLinkFromDocumentInput struct {
	emailRecordID   email.ID
	userAccessor    userPreferencesAccessor
	contentAccessor contentAccessor
	document        documents.Document
}

func makeLinkFromDocument(c ctx.LogContext, input makeLinkFromDocumentInput) (*Link, error) {
	var title, imageURL, description *string
	if isNotEmpty(input.document.Metadata.Title) {
		title = input.document.Metadata.Title
	}
	if isNotEmpty(input.document.Metadata.Image) {
		imageURL = input.document.Metadata.Image
	}
	if isNotEmpty(input.document.Metadata.Description) {
		description = input.document.Metadata.Description
	}
	source, err := input.contentAccessor.GetSourceByID(*input.document.SourceID)
	if err != nil {
		c.Errorf("Error getting source: %s", err.Error())
		return nil, nil
	}
	userDocumentID, err := input.userAccessor.insertDocumentForUserAndReturnID(input.emailRecordID, input.document)
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
		DocumentID:       input.document.ID,
		URL:              *articleLink,
		PaywallReportURL: *paywallReportLink,
		ImageURL:         imageURL,
		Title:            title,
		Description:      description,
		Domain: &Domain{
			Name:      string(source.URL),
			FlagAsset: routes.GetFlagAssetForCountryCode(source.Country),
		},
	}, nil
}

func isNotEmpty(s *string) bool {
	return len(deref.String(s, "")) > 0
}
