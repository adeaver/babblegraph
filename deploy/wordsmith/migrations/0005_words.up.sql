CREATE TABLE IF NOT EXISTS words(
    _id TEXT NOT NULL,
    language TEXT NOT NULL REFERENCES languages(code),
    corpus_id TEXT NOT NULL REFERENCES corpora(_id),

    part_of_speech_id TEXT NOT NULL REFERENCES parts_of_speech(_id),
    lemma_id TEXT NOT NULL REFERENCES lemmas(_id) ON DELETE CASCADE,
    word_text TEXT NOT NULL,

    PRIMARY KEY (_id)
);
