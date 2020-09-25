package documents

type DocumentID string

type dbDocument struct {
	ID       DocumentID     `db:"_id"`
	URL      string         `db:"url"`
	Language *string        `db:"language"`
	Metadata []MetadataPair `db:"metadata"`
}

type MetadataPair struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type Document struct {
	ID       DocumentID     `db:"_id"`
	URL      string         `db:"url"`
	Language *string        `db:"language"`
	Metadata []MetadataPair `db:"metadata"`
}
