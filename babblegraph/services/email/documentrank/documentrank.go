package documentrank

import (
	"babblegraph/model/documents"
	"babblegraph/model/index"
	"babblegraph/services/email/labels"
	"babblegraph/services/email/wordrank"
	"babblegraph/util/math/decimal"
	"babblegraph/wordsmith"
	"fmt"
	"log"
	"sort"

	"github.com/jmoiron/sqlx"
)

type GetDocumentsRankedByLabelInput struct {
	RankedWords      map[wordsmith.LemmaID]wordrank.RankedWord
	LabelSearchTerms map[labels.LabelName][]wordsmith.LemmaID
	DocumentCount    int64
}

type RankedDocument struct {
	DocumentID documents.DocumentID
	Score      decimal.Number
	Rank       int
}

func GetDocumentsRankedByLabel(tx *sqlx.Tx, input GetDocumentsRankedByLabelInput) (map[labels.LabelName][]RankedDocument, error) {
	out := make(map[labels.LabelName][]RankedDocument)
	for labelName, searchTermLemmaIDs := range input.LabelSearchTerms {
		documentIDsForLabel, err := getRelevantDocumentIDForLabel(tx, searchTermLemmaIDs)
		if err != nil {
			return nil, err
		}
		termEntriesForRelevantDocuments, err := getAllTermEntriesForRelevantDocuments(tx, documentIDsForLabel)
		if err != nil {
			return nil, err
		}
		searchTermLemmaIDsMap := make(map[wordsmith.LemmaID]bool)
		for _, lemmaID := range searchTermLemmaIDs {
			searchTermLemmaIDsMap[lemmaID] = true
		}
		inverseDocumentFrequencyForLabel := calculateInverseDocumentFrequencyForLabel(input.RankedWords, searchTermLemmaIDs, input.DocumentCount)
		scoredDocuments := scoreAllDocumentsByRelevance(scoreAllDocumentsByRelevanceInput{
			DocumentsWithTermEntries:                 termEntriesForRelevantDocuments,
			LabelSearchTermsHashset:                  searchTermLemmaIDsMap,
			LabelSearchTermsInverseDocumentFrequency: inverseDocumentFrequencyForLabel,
		})
		sort.Sort(byScore(scoredDocuments))
		var rankedDocuments []RankedDocument
		for i, doc := range scoredDocuments {
			rankedDocuments = append(rankedDocuments, RankedDocument{
				DocumentID: doc.DocumentID,
				Score:      doc.Score,
				Rank:       i + 1,
			})
		}
		out[labelName] = rankedDocuments
	}
	return out, nil
}

func calculateInverseDocumentFrequencyForLabel(rankedWords map[wordsmith.LemmaID]wordrank.RankedWord, labelSearchTerms []wordsmith.LemmaID, documentCount int64) decimal.Number {
	var relevantTermDocumentCount decimal.Number
	for _, lemmaID := range labelSearchTerms {
		relevantTermDocumentCount = relevantTermDocumentCount.Add(decimal.FromInt64(rankedWords[lemmaID].DocumentFrequency))
	}
	inverseDocumentFrequency := decimal.FromInt64(documentCount).Divide(relevantTermDocumentCount)
	return inverseDocumentFrequency
}

func getRelevantDocumentIDForLabel(tx *sqlx.Tx, labelSearchTerms []wordsmith.LemmaID) ([]documents.DocumentID, error) {
	termEntries, err := index.GetTermEntriesForTerms(tx, labelSearchTerms)
	if err != nil {
		return nil, err
	}
	var out []documents.DocumentID
	for _, entry := range termEntries {
		out = append(out, entry.DocumentID)
	}
	return out, nil
}

func getAllTermEntriesForRelevantDocuments(tx *sqlx.Tx, documentIDs []documents.DocumentID) (map[documents.DocumentID][]index.DocumentTermEntry, error) {
	termEntries, err := index.GetTermEntriesForDocuments(tx, documentIDs)
	if err != nil {
		return nil, err
	}
	out := make(map[documents.DocumentID][]index.DocumentTermEntry)
	for _, entry := range termEntries {
		entries, ok := out[entry.DocumentID]
		if !ok {
			entries = make([]index.DocumentTermEntry, 0)
		}
		entries = append(entries, entry)
		out[entry.DocumentID] = entries
	}
	return out, nil
}

type scoreAllDocumentsByRelevanceInput struct {
	DocumentsWithTermEntries                 map[documents.DocumentID][]index.DocumentTermEntry
	LabelSearchTermsHashset                  map[wordsmith.LemmaID]bool
	LabelSearchTermsInverseDocumentFrequency decimal.Number
}

type aggregateDocument struct {
	DocumentID             documents.DocumentID
	TotalWordCount         int64
	LabelSearchTermEntries map[wordsmith.LemmaID]index.DocumentTermEntry
}

func scoreAllDocumentsByRelevance(input scoreAllDocumentsByRelevanceInput) []documentWithScore {
	var docs []aggregateDocument
	for documentID, termEntries := range input.DocumentsWithTermEntries {
		var totalWordCount int64
		searchTermEntries := make(map[wordsmith.LemmaID]index.DocumentTermEntry)
		for _, entry := range termEntries {
			totalWordCount += entry.Count
			if _, ok := input.LabelSearchTermsHashset[entry.TermID]; ok {
				searchTermEntries[entry.TermID] = entry
			}
		}
		docs = append(docs, aggregateDocument{
			DocumentID:             documentID,
			TotalWordCount:         totalWordCount,
			LabelSearchTermEntries: searchTermEntries,
		})
	}
	var out []documentWithScore
	for _, doc := range docs {
		out = append(out, scoreDocument(doc, input.LabelSearchTermsInverseDocumentFrequency))
	}
	return out
}

type documentWithScore struct {
	DocumentID documents.DocumentID
	Score      decimal.Number
}

func scoreDocument(doc aggregateDocument, inverseDocumentFrequency decimal.Number) documentWithScore {
	var relevantTermCount decimal.Number
	for _, termEntry := range doc.LabelSearchTermEntries {
		relevantTermCount = relevantTermCount.Add(decimal.FromInt64(termEntry.Count))
	}
	relevantTermFrequency := relevantTermCount.Divide(decimal.FromInt64(doc.TotalWordCount))
	log.Println(fmt.Sprintf("Term: %f", relevantTermFrequency.ToFloat64()))
	score := relevantTermFrequency.Multiply(inverseDocumentFrequency)
	return documentWithScore{
		DocumentID: doc.DocumentID,
		Score:      score,
	}
}

type byScore []documentWithScore

func (s byScore) Len() int           { return len(s) }
func (s byScore) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s byScore) Less(i, j int) bool { return s[i].Score.GreaterThan(s[j].Score) }
