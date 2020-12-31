CREATE TABLE IF NOT EXISTS lemmas(
    _id TEXT NOT NULL,
    language TEXT NOT NULL REFERENCES languages(code),
    corpus_id TEXT NOT NULL REFERENCES corpora(_id),


    lemma_text TEXT NOT NULL,
    part_of_speech_id TEXT NOT NULL REFERENCES parts_of_speech(_id),

    PRIMARY KEY (_id)
);
