package newsletter

import (
	"babblegraph/model/content"
	"babblegraph/model/documents"
	"babblegraph/model/email"
	"babblegraph/model/routes"
	"babblegraph/util/ctx"
	"babblegraph/util/deref"
	"babblegraph/util/ptr"
	"babblegraph/wordsmith"
	"sort"
	"strings"
	"time"
)

type getSpotlightLemmaForNewsletterInput struct {
	emailRecordID           email.ID
	categories              []Category
	documentIDsInNewsletter []documents.DocumentID
	userAccessor            userPreferencesAccessor
	docsAccessor            documentAccessor
	contentAccessor         contentAccessor
	wordsmithAccessor       wordsmithAccessor
	// TODO: remove this option at the end of the experiment
	excludeFocusContentIneligible bool
}

func getSpotlightLemmaForNewsletter(c ctx.LogContext, input getSpotlightLemmaForNewsletterInput) (*LemmaReinforcementSpotlight, error) {
	if newsletterPreferences := input.userAccessor.getUserNewsletterPreferences(); newsletterPreferences == nil || !newsletterPreferences.ShouldIncludeLemmaReinforcementSpotlight {
		return nil, nil
	}
	c.Infof("Getting spotlight")
	documentIDsToExclude := input.userAccessor.getSentDocumentIDs()
	for _, category := range input.categories {
		for _, l := range category.Links {
			documentIDsToExclude = append(documentIDsToExclude, l.DocumentID)
		}
	}
	for _, documentID := range input.documentIDsInNewsletter {
		documentIDsToExclude = append(documentIDsToExclude, documentID)
	}
	allowableSourceIDs := input.userAccessor.getAllowableSources()
	orderedListOfSpotlightRecords := getOrderedListOfPotentialSpotlightLemmas(input.userAccessor)
	c.Debugf("Ordered spotlight records %+v", orderedListOfSpotlightRecords)
	var preferencesLink string
	if input.userAccessor.getDoesUserHaveAccount() {
		preferencesLink = routes.MakeLoginLinkWithNewsletterPreferencesRedirect()
	} else {
		prefLink, err := routes.MakeNewsletterPreferencesLink(input.userAccessor.getUserID())
		if err != nil {
			return nil, err
		}
		preferencesLink = *prefLink
	}
	reinforcementSpotlight, err := lookupSpotlightForAllPotentialSpotlights(c, lookupSpotlightForAllPotentialSpotlightsInput{
		getSpotlightLemmaForNewsletterInput: input,
		documentIDsToExclude:                documentIDsToExclude,
		potentialSpotlights:                 orderedListOfSpotlightRecords,
		allowableSourceIDs:                  allowableSourceIDs,
		preferencesLink:                     preferencesLink,
		excludeFocusContentIneligible:       input.excludeFocusContentIneligible,
	})
	switch {
	case err != nil:
		return nil, err
	case reinforcementSpotlight != nil:
		return reinforcementSpotlight, nil
	}
	c.Infof("Trying older documents")
	// TODO: create metric here
	return lookupSpotlightForAllPotentialSpotlights(c, lookupSpotlightForAllPotentialSpotlightsInput{
		getSpotlightLemmaForNewsletterInput: input,
		documentIDsToExclude:                documentIDsToExclude,
		potentialSpotlights:                 orderedListOfSpotlightRecords,
		allowableSourceIDs:                  allowableSourceIDs,
		preferencesLink:                     preferencesLink,
		shouldSearchNonRecentDocuments:      true,
		excludeFocusContentIneligible:       input.excludeFocusContentIneligible,
	})
}

type lookupSpotlightForAllPotentialSpotlightsInput struct {
	getSpotlightLemmaForNewsletterInput
	documentIDsToExclude           []documents.DocumentID
	potentialSpotlights            []wordsmith.LemmaID
	shouldSearchNonRecentDocuments bool
	preferencesLink                string
	allowableSourceIDs             []content.SourceID
	excludeFocusContentIneligible  bool
}

func lookupSpotlightForAllPotentialSpotlights(c ctx.LogContext, input lookupSpotlightForAllPotentialSpotlightsInput) (*LemmaReinforcementSpotlight, error) {
	for _, potentialSpotlight := range input.potentialSpotlights {
		documents, err := input.docsAccessor.GetDocumentsForUserForLemma(c, getDocumentsForUserForLemmaInput{
			getDocumentsBaseInput: getDocumentsBaseInput{
				LanguageCode:        input.userAccessor.getLanguageCode(),
				ExcludedDocumentIDs: input.documentIDsToExclude,
				ValidSourceIDs:      input.allowableSourceIDs,
				MinimumReadingLevel: ptr.Int64(input.userAccessor.getReadingLevel().LowerBound),
				MaximumReadingLevel: ptr.Int64(input.userAccessor.getReadingLevel().UpperBound),
			},
			Lemma:           potentialSpotlight,
			Topics:          input.userAccessor.getUserTopics(),
			SearchNonRecent: input.shouldSearchNonRecentDocuments,
		})
		if err != nil {
			return nil, err
		}
		for _, d := range documents {
			if input.excludeFocusContentIneligible && !isDocumentFocusContentEligible(d.Document) {
				continue
			}
			description := deref.String(d.Document.LemmatizedDescription, "")
			for _, lemmaID := range strings.Split(description, " ") {
				if lemmaID == string(potentialSpotlight) {
					link, err := makeLinkFromDocument(c, makeLinkFromDocumentInput{
						emailRecordID:   input.emailRecordID,
						userAccessor:    input.userAccessor,
						contentAccessor: input.contentAccessor,
						document:        d.Document,
					})
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
						PreferencesLink: input.preferencesLink,
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
