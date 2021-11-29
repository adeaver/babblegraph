package newsletter

import (
	"babblegraph/model/documents"
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
		case doc.LanguageCode != input.LanguageCode,
			isIDExcluded(doc.ID, input.ExcludedDocumentIDs),
			!isDomainValid(doc.Domain, input.ValidDomains),
			input.MinimumReadingLevel != nil && *input.MinimumReadingLevel < doc.ReadabilityScore,
			input.MaximumReadingLevel != nil && *input.MaximumReadingLevel > doc.ReadabilityScore,
			doc.LemmatizedDescription == nil || !containsLemma(input.Lemma, *doc.LemmatizedDescription):
			continue
		}
		docs = append(docs, docWithScore)
	}
	return docs, nil
}
