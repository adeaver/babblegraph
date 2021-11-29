package newsletter

import (
	"babblegraph/model/documents"
	"fmt"
	"log"
	"time"
)

type testDocsAccessor struct {
	documents []documents.DocumentWithScore
}

func (t *testDocsAccessor) GetDocumentsForUser(input getDocumentsForUserInput) (*documentsOutput, error) {
	var recentDocuments, nonRecentDocuments []documents.DocumentWithScore
	for _, docWithScore := range t.documents {
		doc := docWithScore.Document
		switch {
		case doc.LanguageCode != input.LanguageCode,
			isIDExcluded(doc.ID, input.ExcludedDocumentIDs),
			!isDomainValid(doc.Domain, input.ValidDomains),
			input.MinimumReadingLevel != nil && *input.MinimumReadingLevel > doc.ReadabilityScore,
			input.MaximumReadingLevel != nil && *input.MaximumReadingLevel < doc.ReadabilityScore,
			input.Topic != nil && !containsTopic(*input.Topic, doc.Topics):
			// no-op
		default:
			recencyBoundary := time.Now().Add(documents.RecencyBiasBoundary).Unix()
			switch {
			case docWithScore.Document.SeedJobIngestTimestamp == nil,
				*docWithScore.Document.SeedJobIngestTimestamp < recencyBoundary:
				nonRecentDocuments = append(nonRecentDocuments, docWithScore)
			case *docWithScore.Document.SeedJobIngestTimestamp >= recencyBoundary:
				recentDocuments = append(recentDocuments, docWithScore)
			}
		}
	}
	return &documentsOutput{
		RecentDocuments:    recentDocuments,
		NonRecentDocuments: nonRecentDocuments,
	}, nil
}

func (t *testDocsAccessor) GetDocumentsForUserForLemma(input getDocumentsForUserForLemmaInput) ([]documents.DocumentWithScore, error) {
	var docs []documents.DocumentWithScore
	for _, docWithScore := range t.documents {
		doc := docWithScore.Document
		switch {
		case doc.LanguageCode != input.LanguageCode:
			log.Println(fmt.Sprintf("Language code does not match, %s", doc.LanguageCode))
		case isIDExcluded(doc.ID, input.ExcludedDocumentIDs):
			log.Println(fmt.Sprintf("ID does not match: %s", doc.ID))
		case !isDomainValid(doc.Domain, input.ValidDomains):
			log.Println(fmt.Sprintf("Domain not valid: %s", doc.Domain))
		case input.MinimumReadingLevel != nil && *input.MinimumReadingLevel < doc.ReadabilityScore:
			log.Println(fmt.Sprintf("Reading score too high: %s", doc.ReadabilityScore))
		case input.MaximumReadingLevel != nil && *input.MaximumReadingLevel > doc.ReadabilityScore:
			log.Println(fmt.Sprintf("Reading score too low: %s", doc.ReadabilityScore))
		case doc.LemmatizedDescription == nil || !containsLemma(input.Lemma, *doc.LemmatizedDescription):
			log.Println(fmt.Sprintf("No description: %+v", doc.LemmatizedDescription))
		default:
			docs = append(docs, docWithScore)
		}
	}
	return docs, nil
}
