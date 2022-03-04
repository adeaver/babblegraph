package tasks

import (
	"babblegraph/model/documents"
	"babblegraph/model/podcasts"
	"babblegraph/util/ctx"
)

func CreateElasticIndexes() error {
	podcasts.CreatePodcastIndexes(ctx.GetDefaultLogContext())
	return documents.CreateDocumentIndex()
}
