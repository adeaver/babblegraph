package wordsmith

type LemmaID string

type Lemma struct {
	ID             LemmaID
	Language       LanguageCode
	CorpusID       CorpusID
	PartOfSpeechID PartOfSpeechID
	LemmaText      string
}

type dbLemma struct {
	ID             LemmaID        `db:"_id"`
	Language       LanguageCode   `db:"language"`
	CorpusID       CorpusID       `db:"corpus_id"`
	PartOfSpeechID PartOfSpeechID `db:"part_of_speech_id"`
	LemmaText      string         `db:"lemma_text"`
}

func (d dbLemma) ToNonDB() Lemma {
	return Lemma{
		ID:             d.ID,
		Language:       d.Language,
		CorpusID:       d.CorpusID,
		PartOfSpeechID: d.PartOfSpeechID,
		LemmaText:      d.LemmaText,
	}
}
