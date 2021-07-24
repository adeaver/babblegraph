package tasks

import "babblegraph/model/documents"

func CreateElasticIndexes() error {
	return documents.CreateDocumentIndex()
}
