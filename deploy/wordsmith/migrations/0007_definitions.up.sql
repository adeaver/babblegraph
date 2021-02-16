CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS definition_mappings(
    _id uuid DEFAULT uuid_generate_v4 (),
    language TEXT NOT NULL REFERENCES languages(code),
    corpus_id TEXT NOT NULL REFERENCES corpora(_id),
    lemma_id TEXT NOT NULL REFERENCES lemmas(_id) ON DELETE CASCADE,
    part_of_speech_id TEXT NOT NULL REFERENCES parts_of_speech(_id),
    english_definition TEXT NOT NULL,
    extra_info TEXT,


    PRIMARY KEY (_id)
);

CREATE INDEX IF NOT EXISTS definition_mappings_lemma_id ON definition_mappings(lemma_id);
