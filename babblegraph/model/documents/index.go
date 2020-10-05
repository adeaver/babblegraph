package documents

import (
	"babblegraph/util/elastic"
	"fmt"

	"github.com/google/uuid"
)

const documentIndexName string = "web_documents"

type documentIndex struct{}

func (d documentIndex) GetName() string {
	return documentIndexName
}

func (d documentIndex) ValidateDocument(document interface{}) error {
	if _, ok := document.(Document); !ok {
		return fmt.Errorf("could not validate interface %+v, to be of type web_document", document)
	}
	return nil
}

func (d documentIndex) GenerateIDForDocument(document interface{}) (*string, error) {
	uuid, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}
	docID := fmt.Sprintf("web_doc-%s", uuid.String())
	return &docID, nil
}

func CreateDocumentIndex() error {
	return elastic.CreateIndex(documentIndex{}, nil)
}

func AssignIDAndIndexDocument(document *Document) (*DocumentID, error) {
	index := documentIndex{}
	idStr, err := index.GenerateIDForDocument(document)
	if err != nil {
		return nil, err
	}
	docID := DocumentID(*idStr)
	document.ID = docID
	if err := elastic.IndexDocument(index, *document); err != nil {
		return nil, err
	}
	return &docID, nil
}
