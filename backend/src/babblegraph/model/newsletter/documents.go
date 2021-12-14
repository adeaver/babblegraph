package newsletter

import (
	"babblegraph/model/contenttopics"
	"babblegraph/model/documents"
	"babblegraph/util/ctx"
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
	Lemma           wordsmith.LemmaID
	Topics          []contenttopics.ContentTopic
	SearchNonRecent bool
}

type documentAccessor interface {
	GetDocumentsForUser(c ctx.LogContext, input getDocumentsForUserInput) (*documentsOutput, error)
	GetDocumentsForUserForLemma(c ctx.LogContext, input getDocumentsForUserForLemmaInput) ([]documents.DocumentWithScore, error)
}

type documentsOutput struct {
	RecentDocuments    []documents.DocumentWithScore
	NonRecentDocuments []documents.DocumentWithScore
}

type DefaultDocumentsAccessor struct{}

func GetDefaultDocumentsAccessor() *DefaultDocumentsAccessor {
	return &DefaultDocumentsAccessor{}
}

func (d *DefaultDocumentsAccessor) GetDocumentsForUser(c ctx.LogContext, input getDocumentsForUserInput) (*documentsOutput, error) {
	dailyEmailDocQueryBuilder := documents.NewDailyEmailDocumentsQueryBuilder()
	dailyEmailDocQueryBuilder.ContainingLemmas(input.Lemmas)
	dailyEmailDocQueryBuilder.ForTopic(input.Topic)
	dailyEmailDocQueryBuilder.WithRecencyBias(documents.RecencyBiasMostRecent)
	recentDocuments, err := documents.ExecuteDocumentQuery(c, dailyEmailDocQueryBuilder, documents.ExecuteDocumentQueryInput{
		LanguageCode:        input.getDocumentsBaseInput.LanguageCode,
		ValidDomains:        input.getDocumentsBaseInput.ValidDomains,
		ExcludedDocumentIDs: input.getDocumentsBaseInput.ExcludedDocumentIDs,
		MinimumReadingLevel: input.getDocumentsBaseInput.MinimumReadingLevel,
		MaximumReadingLevel: input.getDocumentsBaseInput.MaximumReadingLevel,
	})
	if err != nil {
		return nil, err
	}
	dailyEmailDocQueryBuilder.WithRecencyBias(documents.RecencyBiasNotRecent)
	notRecentDocuments, err := documents.ExecuteDocumentQuery(c, dailyEmailDocQueryBuilder, documents.ExecuteDocumentQueryInput{
		LanguageCode:        input.getDocumentsBaseInput.LanguageCode,
		ValidDomains:        input.getDocumentsBaseInput.ValidDomains,
		ExcludedDocumentIDs: input.getDocumentsBaseInput.ExcludedDocumentIDs,
		MinimumReadingLevel: input.getDocumentsBaseInput.MinimumReadingLevel,
		MaximumReadingLevel: input.getDocumentsBaseInput.MaximumReadingLevel,
	})
	if err != nil {
		return nil, err
	}
	return &documentsOutput{
		RecentDocuments:    recentDocuments,
		NonRecentDocuments: notRecentDocuments,
	}, nil
}

func (d *DefaultDocumentsAccessor) GetDocumentsForUserForLemma(c ctx.LogContext, input getDocumentsForUserForLemmaInput) ([]documents.DocumentWithScore, error) {
	spotlightQueryBuilder := documents.NewLemmaSpotlightQueryBuilder(input.Lemma)
	spotlightQueryBuilder.AddTopics(input.Topics)
	recencyBias := documents.RecencyBiasMostRecent
	if input.SearchNonRecent {
		recencyBias = documents.RecencyBiasNotRecent
	}
	spotlightQueryBuilder.WithRecencyBias(recencyBias)
	return documents.ExecuteDocumentQuery(c, spotlightQueryBuilder, documents.ExecuteDocumentQueryInput{
		LanguageCode:        input.getDocumentsBaseInput.LanguageCode,
		ValidDomains:        input.getDocumentsBaseInput.ValidDomains,
		ExcludedDocumentIDs: input.getDocumentsBaseInput.ExcludedDocumentIDs,
		MinimumReadingLevel: input.getDocumentsBaseInput.MinimumReadingLevel,
		MaximumReadingLevel: input.getDocumentsBaseInput.MaximumReadingLevel,
	})
}
