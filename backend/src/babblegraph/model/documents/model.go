package documents

import (
	"babblegraph/wordsmith"
	"time"
)

type Version int64

const (
	Version1 Version = 1
	Version2 Version = 2
	Version3 Version = 3

	CurrentDocumentVersion Version = Version2
)

func (v Version) Ptr() *Version {
	return &v
}

type Type string

const (
	TypeArticle Type = "article"
)

func (t Type) Ptr() *Type {
	return &t
}

func (t Type) Str() string {
	return string(t)
}

type Metadata struct {
	Title              *string    `json:"title,omitempty"`
	Image              *string    `json:"image,omitempty"`
	URL                *string    `json:"url,omitempty"`
	Description        *string    `json:"description,omitempty"`
	PublicationTimeUTC *time.Time `json:"publication_time_utc,omitempty"`
}

type DocumentID string

type Document struct {
	ID           DocumentID             `json:"id"`
	Version      Version                `json:"version"`
	URL          string                 `json:"url"`
	LanguageCode wordsmith.LanguageCode `json:"language_code"`
	DocumentType Type                   `json:"document_type"`
	Metadata     Metadata               `json:"metadata"`
	Domain       string                 `json:"domain"`

	ReadabilityScore  int64              `json:"readability_score"`
	WordStatsVersion1 *WordStatsVersion1 `json:"word_stats_version_1,omitempty"`
	Labels            []string           `json:"labels,omitempty"`
}

type WordStatsVersion1 struct {
	AverageWordRanking  int64 `json:"average_word_ranking"`
	MedianWordRanking   int64 `json:"median_word_ranking"`
	TotalNumberOfWords  int64 `json:"total_number_of_words"` // Includes repeats
	NumberOfUniqueWords int64 `json:"number_of_unique_words"`

	LeastFrequentWordRanking         int64         `json:"least_frequent_word_ranking"`
	LeastFrequentWordExclusion       WordExclusion `json:"least_frequent_word_exclusion"`
	SecondLeastFrequentWordExclusion WordExclusion `json:"second_least_frequent_word_exclusion"`
	ThirdLeastFrequentWordExclusion  WordExclusion `json:"third_least_frequent_word_exclusion"`
}

type WordExclusion struct {
	LeastFrequentRankingWithoutWord int64  `json:"least_frequent_ranking_without_word"`
	WordText                        string `json:"word_text"`
}
