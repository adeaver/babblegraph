package newsletter

import (
	"babblegraph/model/content"
	"fmt"

	"github.com/jmoiron/sqlx"
)

type contentAccessor interface {
	GetSourceByID(sourceID content.SourceID) (*content.Source, error)
}

type DefaultContentAccessor struct {
	allowableSourcesByID map[content.SourceID]content.Source
}

func GetDefaultContentAccessor(tx *sqlx.Tx) (*DefaultContentAccessor, error) {
	sources, err := content.GetAllowableSources(tx)
	if err != nil {
		return nil, err
	}
	sourcesByID := make(map[content.SourceID]content.Source)
	for _, source := range sources {
		sourcesByID[source.ID] = source
	}
	return &DefaultContentAccessor{
		allowableSourcesByID: sourcesByID,
	}, nil
}

func (d *DefaultContentAccessor) GetSourceByID(sourceID content.SourceID) (*content.Source, error) {
	source, ok := d.allowableSourcesByID[sourceID]
	if !ok {
		return nil, fmt.Errorf("Source %s not found", sourceID)
	}
	return &source, nil
}
