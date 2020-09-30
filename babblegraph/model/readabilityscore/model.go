package readabilityscore

import (
	"babblegraph/model/documents"
	"babblegraph/util/math/decimal"
)

type DocumentReadabilityScoreID string

type DocumentReadabilityScore struct {
	ID              DocumentReadabilityScoreID
	DocumentID      documents.DocumentID
	ReadbilityScore decimal.Number
}

type dbDocumentReadabilityScore struct {
	ID              DocumentReadabilityScoreID `db:"id"`
	DocumentID      documents.DocumentID       `db:"document_id"`
	ReadbilityScore float64                    `db:"readbility_score"`
}

func (d dbDocumentReadabilityScore) ToNonDB() DocumentReadabilityScore {
	return DocumentReadabilityScore{
		ID:              d.ID,
		DocumentID:      d.DocumentID,
		ReadbilityScore: decimal.FromFloat64(d.ReadbilityScore),
	}
}
