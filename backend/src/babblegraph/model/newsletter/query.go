package newsletter

import (
	"babblegraph/model/documents"
	"babblegraph/model/domains"
	"babblegraph/model/routes"
	"babblegraph/model/useraccounts"
	"babblegraph/util/deref"
	"babblegraph/wordsmith"
	"fmt"
	"math/rand"
	"time"
)

const (
	defaultNumberOfArticlesPerEmail = 12
	defaultNumberOfTopicsPerEmail   = 4

	minimumDaysSinceLastSpotlight = 3
)

func CreateNewsletter(languageCode wordsmith.LanguageCode, userAccessor userPreferencesAccessor, docsAccessor documentAccessor) (*Newsletter, error) {
	userSubscriptionLevel := userAccessor.getUserSubscriptionLevel()
	switch {
	case userSubscriptionLevel == nil:
		// no-op
	case *userSubscriptionLevel == useraccounts.SubscriptionLevelBetaPremium,
		*userSubscriptionLevel == useraccounts.SubscriptionLevelPremium:
		if userScheduleForDay := userAccessor.getUserScheduleForDay(); userScheduleForDay != nil && !userScheduleForDay.IsActive {
			return nil, nil
		}
	default:
		return nil, fmt.Errorf("Unrecognized subscription level: %s", *userSubscriptionLevel)
	}
	categories, err := getDocumentCategories(languageCode, userAccessor, docsAccessor)
	if err != nil {
		return nil, err
	}
	return &Newsletter{
		// UserID       users.UserID           `json:"user_id"`
		LanguageCode: languageCode,
		Body: NewsletterBody{
			// LemmaReinforcementSpotlight *LemmaReinforcementSpotlight `json:"lemma_reinforcement_spotlight,omitempty"`
			Categories:        categories,
			SetTopicsLink:     nil,
			ReinforcementLink: "",
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

func makeLinkFromDocument(userAccessor userPreferencesAccessor, doc documents.Document) (*Link, error) {
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
		return nil, err
	}
	// TODO: In order to avoid duplicate documents between emails
	// we'll need to do the following:
	// 1) Create EmailRecordID and EmailRecord when creating Newsletter
	// 2) Modify email record to have a null sent_at date for pending emails
	// 3) Add a method onto userAccessor to insert document for user and return ID
	return &Link{
		// URL              string  `json:"url"`
		// PaywallReportURL string  `json:"paywall_report_url"`
		ImageURL:    imageURL,
		Title:       title,
		Description: description,
		Domain: &Domain{
			Name:      string(domain.Domain),
			FlagAsset: routes.GetFlagAssetForCountryCode(domain.Country),
		},
	}, nil
}

func isNotEmpty(s *string) bool {
	return len(deref.String(s, "")) > 0
}
