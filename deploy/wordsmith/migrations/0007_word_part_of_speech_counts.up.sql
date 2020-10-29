CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE word_part_of_speech_counts(
    _id uuid DEFAULT uuid_generate_v4 (),
    language TEXT NOT NULL REFERENCES languages(code),
    corpus_id uuid NOT NULL REFERENCES corpora(_id),

    word_id uuid NOT NULL REFERENCES words(_id),
    occurrences NUMERIC NOT NULL,

    PRIMARY KEY (_id)
);
