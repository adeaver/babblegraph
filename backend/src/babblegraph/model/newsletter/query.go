package newsletter

import (
	"babblegraph/model/documents"
	"babblegraph/model/domains"
	"babblegraph/model/email"
	"babblegraph/model/routes"
	"babblegraph/model/useraccounts"
	"babblegraph/util/deref"
	"babblegraph/util/ptr"
	"fmt"
	"log"
	"math/rand"
	"time"
)

const (
	defaultNumberOfArticlesPerEmail = 12
	defaultNumberOfTopicsPerEmail   = 4

	minimumDaysSinceLastSpotlight = 3
)

func CreateNewsletter(wordsmithAccessor wordsmithAccessor, emailAccessor emailAccessor, userAccessor userPreferencesAccessor, docsAccessor documentAccessor) (*Newsletter, error) {
	emailRecordID := email.NewEmailRecordID()
	if err := emailAccessor.InsertEmailRecord(emailRecordID, userAccessor.getUserID()); err != nil {
		return nil, err
	}
	var numberOfDocumentsInNewsletter *int
	userSubscriptionLevel := userAccessor.getUserSubscriptionLevel()
	switch {
	case userSubscriptionLevel == nil:
		// no-op
	case *userSubscriptionLevel == useraccounts.SubscriptionLevelBetaPremium,
		*userSubscriptionLevel == useraccounts.SubscriptionLevelPremium:
		if userScheduleForDay := userAccessor.getUserScheduleForDay(); userScheduleForDay != nil {
			if !userScheduleForDay.IsActive {
				return nil, nil
			}
			numberOfDocumentsInNewsletter = ptr.Int(userScheduleForDay.NumberOfArticles)
		}

	default:
		return nil, fmt.Errorf("Unrecognized subscription level: %s", *userSubscriptionLevel)
	}
	categories, err := getDocumentCategories(getDocumentCategoriesInput{
		emailRecordID:                 emailRecordID,
		languageCode:                  userAccessor.getLanguageCode(),
		userAccessor:                  userAccessor,
		docsAccessor:                  docsAccessor,
		numberOfDocumentsInNewsletter: numberOfDocumentsInNewsletter,
	})
	if err != nil {
		return nil, err
	}
	var setTopicsLink *string
	switch {
	case len(userAccessor.getUserTopics()) > 0:
		// no-op
	case userAccessor.getDoesUserHaveAccount():
		setTopicsLink = ptr.String(routes.MakeLoginLinkWithContentTopicsRedirect())
	default:
		setTopicsLink, err = routes.MakeSetTopicsLink(userAccessor.getUserID())
		if err != nil {
			return nil, err
		}
	}
	var reinforcementLink *string
	if userAccessor.getDoesUserHaveAccount() {
		reinforcementLink = ptr.String(routes.MakeLoginLinkWithReinforcementRedirect())
	} else {
		reinforcementLink, err = routes.MakeWordReinforcementLink(userAccessor.getUserID())
		if err != nil {
			return nil, err
		}
	}
	spotlightRecord, err := getSpotlightLemmaForNewsletter(getSpotlightLemmaForNewsletterInput{
		emailRecordID:     emailRecordID,
		categories:        categories,
		userAccessor:      userAccessor,
		docsAccessor:      docsAccessor,
		wordsmithAccessor: wordsmithAccessor,
	})
	if err != nil {
		return nil, err
	}
	return &Newsletter{
		UserID:        userAccessor.getUserID(),
		EmailRecordID: emailRecordID,
		LanguageCode:  userAccessor.getLanguageCode(),
		Body: NewsletterBody{
			LemmaReinforcementSpotlight: spotlightRecord,
			Categories:                  categories,
			SetTopicsLink:               setTopicsLink,
			ReinforcementLink:           *reinforcementLink,
		},
	}, nil
}

func getAllowableDomains(accessor userPreferencesAccessor) ([]string, error) {
	currentUserDomainCounts := accessor.getUserDomainCounts()
	domainCountByDomain := make(map[string]int64)
	for _, domainCount := range currentUserDomainCounts {
		domainCountByDomain[domainCount.Domain] = domainCount.Count
	}
	var out []string
	for _, d := range domains.GetDomains() {
		countForDomain, ok := domainCountByDomain[d]
		if ok {
			metadata, err := domains.GetDomainMetadata(d)
			if err != nil {
				return nil, err
			}
			if metadata.NumberOfMonthlyFreeArticles != nil && countForDomain >= *metadata.NumberOfMonthlyFreeArticles {
				continue
			}
		}
		out = append(out, d)
	}
	return out, nil
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

func makeLinkFromDocument(emailRecordID email.ID, userAccessor userPreferencesAccessor, doc documents.Document) (*Link, error) {
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
	domain, err := domains.GetDomainMetadata(doc.Domain)
	if err != nil {
		log.Println(fmt.Sprintf("Error getting domain: %s", err.Error()))
		return nil, nil
	}
	userDocumentID, err := userAccessor.insertDocumentForUserAndReturnID(emailRecordID, doc)
	if err != nil {
		return nil, err
	}
	articleLink, err := routes.MakeArticleLink(*userDocumentID)
	if err != nil {
		log.Println(fmt.Sprintf("Error making article link: %s", err.Error()))
		return nil, nil
	}
	paywallReportLink, err := routes.MakePaywallReportLink(*userDocumentID)
	if err != nil {
		log.Println(fmt.Sprintf("Error making paywall report link: %s", err.Error()))
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
