package wordsmith

import (
	"fmt"

	"github.com/jmoiron/sqlx"
)

func GetWords(words []string, language LanguageCode) ([]Word, error) {
	var matches []dbWord
	if err := withTx(func(tx *sqlx.Tx) error {
		query, args, err := sqlx.In("SELECT * FROM words WHERE language=? AND word_text IN (?)", language.Str(), words)
		if err != nil {
			return err
		}
		sql := tx.Rebind(query)
		return tx.Select(&matches, sql, args...)
	}); err != nil {
		return nil, err
	}
	var out []Word
	for _, match := range matches {
		out = append(out, match.ToWord())
	}
	return out, nil
}

func LookupLemmas(lemmas []string, language LanguageCode) ([]Lemma, error) {
	var matches []dbLemma
	if err := withTx(func(tx *sqlx.Tx) error {
		query, args, err := sqlx.In("SELECT * FROM lemmas WHERE language=? AND lemma_text IN (?)", language.Str(), lemmas)
		if err != nil {
			return err
		}
		sql := tx.Rebind(query)
		return tx.Select(&matches, sql, args...)
	}); err != nil {
		return nil, err
	}
	var out []Lemma
	for _, match := range matches {
		out = append(out, match.ToNonDB())
	}
	return out, nil
}

const getAllLemmasForLanguageQuery = "SELECT * FROM lemmas WHERE language='%s'"

func GetAllLemmasForLanguage(language LanguageCode) ([]Lemma, error) {
	var matches []dbLemma
	if err := withTx(func(tx *sqlx.Tx) error {
		return tx.Select(&matches, fmt.Sprintf(getAllLemmasForLanguageQuery, language.Str()))
	}); err != nil {
		return nil, err
	}
	var out []Lemma
	for _, match := range matches {
		out = append(out, match.ToNonDB())
	}
	return out, nil
}
