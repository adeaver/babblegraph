CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE word_bigram_counts(
    _id uuid DEFAULT uuid_generate_v4 (),
    language TEXT NOT NULL REFERENCES languages(code),
    corpus_id uuid NOT NULL REFERENCES corpora(_id),

    first_token_id uuid NOT NULL REFERENCES words(_id),
    second_token_id uuid NOT NULL REFERENCES words(_id),
    occurrences NUMERIC NOT NULL,

    PRIMARY KEY (_id)
);
