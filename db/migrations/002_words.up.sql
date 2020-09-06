CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE words(
    _id uuid DEFAULT uuid_generate_v4 (),
    word TEXT NOT NULL,
    lemma_id uuid NOT NULL REFERENCES lemmas(_id) ON DELETE CASCADE,
    part_of_speech TEXT NOT NULL,
    language TEXT NOT NULL,

    PRIMARY KEY (_id)
);

CREATE UNIQUE INDEX IF NOT EXISTS word_lemma_unique_idx ON words(word, lemma_id, part_of_speech);
