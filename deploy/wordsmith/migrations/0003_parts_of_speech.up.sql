CREATE TABLE IF NOT EXISTS parts_of_speech(
    _id TEXT NOT NULL,
    language TEXT NOT NULL REFERENCES languages(code),
    corpus_id TEXT NOT NULL REFERENCES corpora(_id),

    code TEXT NOT NULL,

    PRIMARY KEY (_id)
);
