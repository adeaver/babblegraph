CREATE TABLE IF NOT EXISTS word_bigram_counts(
    _id TEXT NOT NULL,
    language TEXT NOT NULL REFERENCES languages(code),
    corpus_id TEXT NOT NULL REFERENCES corpora(_id),
    first_word_text TEXT NOT NULL,
    first_word_lemma_id TEXT NOT NULL REFERENCES lemmas(_id) ON DELETE CASCADE,
    second_word_text TEXT NOT NULL,
    second_word_lemma_id TEXT NOT NULL REFERENCES lemmas(_id) ON DELETE CASCADE,
    count NUMERIC NOT NULL,

    PRIMARY KEY (_id)
);

CREATE INDEX IF NOT EXISTS word_bigram_words_text_idx ON word_bigram_counts(first_word_text, second_word_text);
