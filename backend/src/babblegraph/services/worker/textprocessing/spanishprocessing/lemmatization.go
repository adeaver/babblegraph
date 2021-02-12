package spanishprocessing

import (
	"babblegraph/util/math/decimal"
	"babblegraph/util/ptr"
	"babblegraph/util/text"
	"babblegraph/wordsmith"

	"github.com/jmoiron/sqlx"
)

func LemmatizeText(t string) ([]*wordsmith.LemmaID, error) {
	tokens := text.Tokenize(t)
	wordsByText, err := getWordsByText(tokens)
	if err != nil {
		return nil, err
	}
	return convertTokensToLemmas(tokens, wordsByText)
}

func getWordsByText(tokens []string) (map[string][]wordsmith.Word, error) {
	var words []wordsmith.Word
	if err := wordsmith.WithWordsmithTx(func(tx *sqlx.Tx) error {
		var err error
		words, err = wordsmith.GetWordsByText(tx, wordsmith.SpanishUPCWikiCorpus, tokens)
		return err
	}); err != nil {
		return nil, err
	}
	out := make(map[string][]wordsmith.Word)
	for _, w := range words {
		arr, _ := out[w.WordText]
		out[w.WordText] = append(arr, w)
	}
	return out, nil
}

// This function will return a parallel list of tokens -> lemma ID
// a nil entry means that we don't know what the lemma is
func convertTokenToLemmas(tokens []string, wordsByText map[string][]wordsmith.Word) ([]*wordsmith.LemmaID, error) {
	var out []*wordsmith.LemmaID
	for idx, token := range tokens {
		// Grab all the words that map to this particular token
		wordsForToken, _ := wordsByText[token]
		switch {
		case len(wordsForToken) == 0:
			// There are no known wordsmith words that map
			// to this token, so add nil to our output list
			out = append(out, nil)
		case len(wordsForToken) == 1:
			// One to one mapping
			out = append(out, wordsForToken[0].LemmaID)
		case len(wordsForToken) >= 2:
			var priorWord, nextWord *string
			switch {
			case idx == 0:
				nextWord = ptr.String(tokens[idx+1])
			case idx == len(tokens)-1:
				priorWord = ptr.String(tokens[idx-1])
			case idx > 0 && idx < len(tokens)-1:
				nextWord = ptr.String(tokens[idx+1])
				priorWord = ptr.String(tokens[idx-1])
			default:
				panic("unreachable")
			}
			var bigramCountsEndingInToken, bigramCountsStartingInToken []wordsmith.WordBigramCount
			if err := wordsmith.WithWordsmithTx(func(tx *sqlx.Tx) error {
				var err error
				if priorWord != nil {
					bigramCountsEndingInToken, err = wordsmith.GetWordBigramCountsByWordText(tx, wordsmith.SpanishUPCWikiCorpus, *priorWord, token)
					if err != nil {
						return err
					}
				}
				if nextWord != nil {
					bigramCountsStartingInToken, err = wordsmith.GetWordBigramCountsByWordText(tx, wordsmith.SpanishUPCWikiCorpus, token, *nextWord)
					if err != nil {
						return err
					}
				}
				return err
			}); err != nil {
				return nil, err
			}
			bestWordChoice := pickBestWordUsingBigrams(pickBestWordUsingBigramInput{
				wordChoices:                 wordsForToken,
				bigramCountsEndingInToken:   bigramCountsEndingInToken,
				bigramCountsStartingInToken: bigramCountsStartingInToken,
			})
			out = append(out, &bestWordChoice.LemmaID)
		default:
			panic("unreachable")
		}
	}
	return out, nil
}

type pickBestWordUsingBigramInput struct {
	wordChoices                 []wordsmith.Word
	bigramCountsEndingInToken   []wordsmith.WordBigramCount
	bigramCountsStartingInToken []wordsmith.WordBigramCount
}

func pickBestWordUsingBigrams(input pickBestWordUsingBigramInput) wordsmith.Word {
	type wordChoice struct {
		word        wordsmith.Word
		probability decimal.Number
	}
	var currentBestChoice *wordChoice
	for _, word := range input.wordChoices {
		probabilityOfEndingInToken := calculateProbabilityOfEndingInToken(word, input.bigramCountsEndingInToken)
		probabilityOfStartingWithToken := calculateProbabilityOfStartingWithToken(word, input.bigramCountsStartingInToken)
		probabilityOfWord := probabilityOfEndingInToken.Multiply(probabilityOfStartingWithToken)
		if currentBestChoice == nil || currentBestChoice.probability.LessThan(probabilityOfWord) {
			currentBestChoice = &wordChoice{
				word:        word,
				probability: probabilityOfWord,
			}
		}
	}
	if currentBestChoice == nil {
		panic("there should be at least one word")
	}
	return currentBestChoice.word
}

func calculateProbabilityOfEndingInToken(word wordsmith.Word, bigramCountsEndingInToken []wordsmith.WordBigramCount) decimal.Number {
	return calculateBigramProbability(calculateBigramProbabilityInput{
		word:         word,
		bigramCounts: bigramCountsEndingInToken,
		isCurrentWord: func(word wordsmith.Word, bigramCount wordsmith.WordBigramCount) bool {
			return bigramCount.SecondWord.LemmaID == word.LemmaID
		},
	})
}

func calculateProbabilityOfStartingWithToken(word wordsmith.Word, bigramCountsStartingWithToken []wordsmith.WordBigramCount) decimal.Number {
	return calculateBigramProbability(calculateBigramProbabilityInput{
		word:         word,
		bigramCounts: bigramCountsStartingWithToken,
		isCurrentWord: func(word wordsmith.Word, bigramCount wordsmith.WordBigramCount) bool {
			return bigramCount.FirstWord.LemmaID == word.LemmaID
		},
	})
}

type calculateBigramProbabilityInput struct {
	word          wordsmith.Word
	bigramCounts  []wordsmith.WordBigramCount
	isCurrentWord func(wordsmith.Word, wordsmith.WordBigramCount) bool
}

func calculateBigramProbability(input calculateBigramProbabilityInput) decimal.Number {
	totalCountOfBigrams := decimal.FromInt64(0)
	totalCountForCurrentWord := decimal.FromInt64(1)
	for _, bigramCount := range input.bigramCounts {
		if input.isCurrentWord(input.word, bigramCount) {
			totalCountForCurrentWord.Add(decimal.FromInt64(bigramCount.Count))
		}
		totalCountOfBigrams := totalCountOfBigrams.Add(decimal.FromInt64(bigramCount.Count))
	}
	if totalCountOfBigrams.EqualTo(decimal.FromInt64(0)) {
		return decimal.FromInt64(1)
	}
	return totalCountForCurrentWord.Divide(totalCountOfBigrams)
}
