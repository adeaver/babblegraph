package newsletter

import (
	"babblegraph/model/email"
	"babblegraph/model/routes"
	"babblegraph/model/useraccounts"
	"babblegraph/util/deref"
	"babblegraph/util/ptr"
	"babblegraph/wordsmith"
	"sort"
	"strings"
	"time"
)

type getSpotlightLemmaForNewsletterInput struct {
	emailRecordID     email.ID
	categories        []Category
	userAccessor      userPreferencesAccessor
	docsAccessor      documentAccessor
	wordsmithAccessor wordsmithAccessor
}

func getSpotlightLemmaForNewsletter(input getSpotlightLemmaForNewsletterInput) (*LemmaReinforcementSpotlight, error) {
	userSubscriptionLevel := input.userAccessor.getUserSubscriptionLevel()
	switch {
	case userSubscriptionLevel == nil:
		return nil, nil
	case *userSubscriptionLevel == useraccounts.SubscriptionLevelBetaPremium,
		*userSubscriptionLevel == useraccounts.SubscriptionLevelPremium:
		if newsletterPreferences := input.userAccessor.getUserNewsletterPreferences(); newsletterPreferences == nil || !newsletterPreferences.ShouldIncludeLemmaReinforcementSpotlight {
			return nil, nil
		}
	default:
		panic("unreachable")
	}
	documentIDsToExclude := input.userAccessor.getSentDocumentIDs()
	for _, category := range input.categories {
		for _, l := range category.Links {
			documentIDsToExclude = append(documentIDsToExclude, l.DocumentID)
		}
	}
	allowableDomains, err := getAllowableDomains(input.userAccessor)
	if err != nil {
		return nil, err
	}
	orderedListOfSpotlightRecords := getOrderedListOfPotentialSpotlightLemmas(input.userAccessor)
	for _, potentialSpotlight := range orderedListOfSpotlightRecords {
		documents, err := input.docsAccessor.GetDocumentsForUserForLemma(getDocumentsForUserForLemmaInput{
			getDocumentsBaseInput: getDocumentsBaseInput{
				LanguageCode:        input.userAccessor.getLanguageCode(),
				ExcludedDocumentIDs: documentIDsToExclude,
				ValidDomains:        allowableDomains,
				MinimumReadingLevel: ptr.Int64(input.userAccessor.getReadingLevel().LowerBound),
				MaximumReadingLevel: ptr.Int64(input.userAccessor.getReadingLevel().UpperBound),
			},
			Lemma:  potentialSpotlight,
			Topics: input.userAccessor.getUserTopics(),
		})
		if err != nil {
			return nil, err
		}
		for _, d := range documents {
			description := deref.String(d.Document.LemmatizedDescription, "")
			for _, lemmaID := range strings.Split(description, " ") {
				if lemmaID == string(potentialSpotlight) {
					link, err := makeLinkFromDocument(input.emailRecordID, input.userAccessor, d.Document)
					switch {
					case err != nil:
						return nil, err
					case link == nil:
						continue
					}
					if err := input.userAccessor.insertSpotlightReinforcementRecord(potentialSpotlight); err != nil {
						return nil, err
					}
					lemma, err := input.wordsmithAccessor.GetLemmaByID(potentialSpotlight)
					if err != nil {
						return nil, err
					}
					if err := input.userAccessor.insertSpotlightReinforcementRecord(potentialSpotlight); err != nil {
						return nil, err
					}
					return &LemmaReinforcementSpotlight{
						LemmaText:       lemma.LemmaText,
						Document:        *link,
						PreferencesLink: routes.MakeLoginLinkWithNewsletterPreferencesRedirect(),
					}, nil
				}
			}
		}
	}
	return nil, nil
}

func getOrderedListOfPotentialSpotlightLemmas(userAccessor userPreferencesAccessor) []wordsmith.LemmaID {
	lemmaReinforcementSpotlightRecords := userAccessor.getSpotlightRecordsOrderedBySentOn()
	lemmaSpotlightRecordSentOnTimeByID := make(map[wordsmith.LemmaID]time.Time)
	for _, spotlightRecord := range lemmaReinforcementSpotlightRecords {
		lemmaSpotlightRecordSentOnTimeByID[spotlightRecord.LemmaID] = spotlightRecord.LastSentOn
	}
	now := time.Now()
	var lemmasNotSent, sentLemmas []wordsmith.LemmaID
	for _, lemmaID := range userAccessor.getTrackingLemmas() {
		lastSent, ok := lemmaSpotlightRecordSentOnTimeByID[lemmaID]
		if !ok {
			lemmasNotSent = append(lemmasNotSent, lemmaID)
		} else {
			if lastSent.Add(minimumDaysSinceLastSpotlight * 24 * time.Hour).Before(now) {
				sentLemmas = append(sentLemmas, lemmaID)
			}
		}
	}
	sort.Slice(sentLemmas, func(i, j int) bool {
		iSentOn, _ := lemmaSpotlightRecordSentOnTimeByID[sentLemmas[i]]
		jSentOn, _ := lemmaSpotlightRecordSentOnTimeByID[sentLemmas[j]]
		return iSentOn.Before(jSentOn)
	})
	return append(lemmasNotSent, sentLemmas...)
}
