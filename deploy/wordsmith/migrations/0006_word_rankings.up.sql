CREATE TABLE IF NOT EXISTS word_rankings(
    _id TEXT NOT NULL,
    language TEXT NOT NULL REFERENCES languages(code),
    corpus_id TEXT NOT NULL REFERENCES corpora(_id),

    word TEXT NOT NULL,
    ranking NUMERIC NOT NULL,
    count NUMERIC NOT NULL,

    PRIMARY KEY (_id)
);

CREATE INDEX IF NOT EXISTS word_rankings_word_idx ON word_rankings(word);
CREATE INDEX IF NOT EXISTS word_rankings_seq_idx ON word_rankings(ranking);
