CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS phrase_definitions(
    _id uuid DEFAULT uuid_generate_v4 (),
    language TEXT NOT NULL REFERENCES languages(code),
    corpus_id TEXT NOT NULL REFERENCES corpora(_id),
    phrase TEXT NOT NULL,
    definition TEXT NOT NULL,

    PRIMARY KEY (_id)
);

CREATE TABLE IF NOT EXISTS lemma_phrase_definition_mappings(
    _id uuid DEFAULT uuid_generate_v4 (),
    language TEXT NOT NULL REFERENCES languages(code),
    corpus_id TEXT NOT NULL REFERENCES corpora(_id),
    lemma_phrase TEXT NOT NULL,
    phrase_definition_id uuid REFERENCES phrase_definitions(_id),

    PRIMARY KEY (_id)
);

CREATE INDEX IF NOT EXISTS lemma_pharse_for_mappings_idx ON lemma_phrase_definition_mappings(lemma_phrase);
