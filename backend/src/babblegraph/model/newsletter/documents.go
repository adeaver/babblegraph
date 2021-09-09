package newsletter

import (
	"babblegraph/model/contenttopics"
	"babblegraph/model/documents"
	"babblegraph/wordsmith"
)

type getDocumentsBaseInput struct {
	LanguageCode        wordsmith.LanguageCode
	ExcludedDocumentIDs []documents.DocumentID
	ValidDomains        []string
	MinimumReadingLevel *int64
	MaximumReadingLevel *int64
}

type getDocumentsForUserInput struct {
	getDocumentsBaseInput
	Topic  *contenttopics.ContentTopic
	Lemmas []wordsmith.LemmaID
}

type getDocumentsForUserForLemmaInput struct {
	getDocumentsBaseInput
	Lemma  wordsmith.LemmaID
	Topics []contenttopics.ContentTopic
}

type documentAccessor interface {
	GetDocumentsForUser(input getDocumentsForUserInput) ([]documents.DocumentWithScore, error)
	GetDocumentsForUserForLemma(input getDocumentsForUserForLemmaInput) ([]documents.DocumentWithScore, error)
}

type DefaultDocumentsAccessor struct{}

func GetDefaultDocumentsAccessor() *DefaultDocumentsAccessor {
	return &DefaultDocumentsAccessor{}
}

func (d *DefaultDocumentsAccessor) GetDocumentsForUser(input getDocumentsForUserInput) ([]documents.DocumentWithScore, error) {
	dailyEmailDocQueryBuilder := documents.NewDailyEmailDocumentsQueryBuilder()
	dailyEmailDocQueryBuilder.ContainingLemmas(input.Lemmas)
	dailyEmailDocQueryBuilder.ForTopic(input.Topic)
	return documents.ExecuteDocumentQuery(dailyEmailDocQueryBuilder, documents.ExecuteDocumentQueryInput{
		LanguageCode:        input.getDocumentsBaseInput.LanguageCode,
		ValidDomains:        input.getDocumentsBaseInput.ValidDomains,
		ExcludedDocumentIDs: input.getDocumentsBaseInput.ExcludedDocumentIDs,
		MinimumReadingLevel: input.getDocumentsBaseInput.MinimumReadingLevel,
		MaximumReadingLevel: input.getDocumentsBaseInput.MaximumReadingLevel,
	})
}

func (d *DefaultDocumentsAccessor) GetDocumentsForUserForLemma(input getDocumentsForUserForLemmaInput) ([]documents.DocumentWithScore, error) {
	spotlightQueryBuilder := documents.NewLemmaSpotlightQueryBuilder(input.Lemma)
	spotlightQueryBuilder.AddTopics(input.Topics)
	return documents.ExecuteDocumentQuery(spotlightQueryBuilder, documents.ExecuteDocumentQueryInput{
		LanguageCode:        input.getDocumentsBaseInput.LanguageCode,
		ValidDomains:        input.getDocumentsBaseInput.ValidDomains,
		ExcludedDocumentIDs: input.getDocumentsBaseInput.ExcludedDocumentIDs,
		MinimumReadingLevel: input.getDocumentsBaseInput.MinimumReadingLevel,
		MaximumReadingLevel: input.getDocumentsBaseInput.MaximumReadingLevel,
	})
}
