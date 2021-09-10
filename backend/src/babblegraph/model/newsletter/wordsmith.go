package newsletter

import (
	"babblegraph/wordsmith"

	"github.com/jmoiron/sqlx"
)

type wordsmithAccessor interface {
	GetLemmaByID(lemmaID wordsmith.LemmaID) (*wordsmith.Lemma, error)
}

type DefaultWordsmithAccessor struct{}

func GetDefaultWordsmithAccessor() *DefaultWordsmithAccessor {
	return &DefaultWordsmithAccessor{}
}

func (d *DefaultWordsmithAccessor) GetLemmaByID(lemmaID wordsmith.LemmaID) (*wordsmith.Lemma, error) {
	var lemma *wordsmith.Lemma
	if err := wordsmith.WithWordsmithTx(func(tx *sqlx.Tx) error {
		var err error
		lemma, err = wordsmith.GetLemmaByID(tx, lemmaID)
		return err
	}); err != nil {
		return nil, err
	}
	return lemma, nil
}
