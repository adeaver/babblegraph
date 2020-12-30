CREATE TABLE word_rankings(
    _id TEXT NOT NULL,
    language TEXT NOT NULL REFERENCES languages(code),
    corpus_id TEXT NOT NULL REFERENCES corpora(_id),

    word TEXT NOT NULL,
    ranking NUMERIC NOT NULL,
    count NUMERIC NOT NULL,

    PRIMARY KEY (_id)
);
