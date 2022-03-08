package newsletter

import (
	"babblegraph/model/documents"
	"babblegraph/util/ctx"
	"time"
)

type testDocsAccessor struct {
	documents []documents.DocumentWithScore
}

func (t *testDocsAccessor) GetDocumentsForUser(c ctx.LogContext, input getDocumentsForUserInput) (*documentsOutput, error) {
	var recentDocuments, nonRecentDocuments []documents.DocumentWithScore
	for _, docWithScore := range t.documents {
		doc := docWithScore.Document
		switch {
		case doc.LanguageCode != input.LanguageCode,
			isIDExcluded(doc.ID, input.ExcludedDocumentIDs),
			!isSourceValid(doc.SourceID, input.ValidSourceIDs),
			input.MinimumReadingLevel != nil && *input.MinimumReadingLevel > doc.ReadabilityScore,
			input.MaximumReadingLevel != nil && *input.MaximumReadingLevel < doc.ReadabilityScore,
			input.Topic != nil && !containsTopic(*input.Topic, doc.TopicIDs):
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

func (t *testDocsAccessor) GetDocumentsForUserForLemma(c ctx.LogContext, input getDocumentsForUserForLemmaInput) ([]documents.DocumentWithScore, error) {
	var docs []documents.DocumentWithScore
	for _, docWithScore := range t.documents {
		doc := docWithScore.Document
		switch {
		case doc.LanguageCode != input.LanguageCode:
			c.Debugf("Language code does not match, %s", doc.LanguageCode)
		case isIDExcluded(doc.ID, input.ExcludedDocumentIDs):
			c.Debugf("ID does not match: %s", doc.ID)
		case !isSourceValid(doc.SourceID, input.ValidSourceIDs):
			c.Debugf("Domain not valid: %s", doc.Domain)
		case input.MinimumReadingLevel != nil && *input.MinimumReadingLevel > doc.ReadabilityScore:
			c.Debugf("Reading score too high: %+v", doc.ReadabilityScore)
		case input.MaximumReadingLevel != nil && *input.MaximumReadingLevel < doc.ReadabilityScore:
			c.Debugf("Reading score too low: %+v", doc.ReadabilityScore)
		case doc.LemmatizedDescription == nil || !containsLemma(input.Lemma, *doc.LemmatizedDescription):
			c.Debugf("No description: %+v", doc.LemmatizedDescription)
		default:
			docs = append(docs, docWithScore)
		}
	}
	return docs, nil
}
