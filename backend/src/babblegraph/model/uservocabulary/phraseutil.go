package uservocabulary

import (
	"babblegraph/wordsmith"
	"strings"

	"github.com/jmoiron/sqlx"
)

func GetLemmaIDPhrasesForPhrase(phrase string) ([][]wordsmith.LemmaID, error) {
	var phrases [][]wordsmith.LemmaID = nil
	phraseWords := strings.Split(phrase, " ")
	if err := wordsmith.WithWordsmithTx(func(tx *sqlx.Tx) error {
		words, err := wordsmith.GetWordsByText(tx, wordsmith.SpanishUPCWikiCorpus, phraseWords)
		if err != nil {
			return err
		}
		lemmaIDsByText := make(map[string][]wordsmith.LemmaID)
		var lemmaIDs []wordsmith.LemmaID
		for _, w := range words {
			lemmaIDs = append(lemmaIDs, w.LemmaID)
			lemmaIDsByText[strings.ToLower(w.WordText)] = append(lemmaIDsByText[strings.ToLower(w.WordText)], w.LemmaID)
		}
		var lemmaIDsByWord [][]wordsmith.LemmaID
		for _, w := range phraseWords {
			lemmaIDs, ok := lemmaIDsByText[strings.ToLower(w)]
			if !ok {
				return nil
			}
			lemmaIDsByWord = append(lemmaIDsByWord, lemmaIDs)
		}
		phrases = makeLemmaIDPhrases(lemmaIDsByWord, []wordsmith.LemmaID{})
		return nil
	}); err != nil {
		return nil, err
	}
	return phrases, nil
}

func makeLemmaIDPhrases(lemmaIDsByWord [][]wordsmith.LemmaID, currentPhrase []wordsmith.LemmaID) [][]wordsmith.LemmaID {
	if len(lemmaIDsByWord) == 0 {
		return [][]wordsmith.LemmaID{currentPhrase}
	}
	var out [][]wordsmith.LemmaID
	for _, lemmaID := range lemmaIDsByWord[0] {
		out = append(out, makeLemmaIDPhrases(lemmaIDsByWord[1:], append(currentPhrase, lemmaID))...)
	}
	return out
}
