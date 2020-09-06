CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE lemmas(
    _id uuid DEFAULT uuid_generate_v4 (),
    lemma TEXT NOT NULL,
    part_of_speech TEXT NOT NULL,
    language TEXT NOT NULL,

    PRIMARY KEY (_id)
);

CREATE INDEX IF NOT EXISTS lemma_word_idx ON lemmas(lemma);
