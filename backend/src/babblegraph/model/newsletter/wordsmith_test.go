package newsletter

import "babblegraph/wordsmith"

type testWordsmithAccessor struct {
	lemmasByID map[wordsmith.LemmaID]wordsmith.Lemma
}

func (t *testWordsmithAccessor) GetLemmaByID(lemmaID wordsmith.LemmaID) (*wordsmith.Lemma, error) {
	if lemma, ok := t.lemmasByID[lemmaID]; ok {
		return &lemma, nil
	}
	return nil, nil
}
