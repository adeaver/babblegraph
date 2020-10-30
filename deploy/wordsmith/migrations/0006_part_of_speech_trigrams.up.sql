CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE part_of_speech_trigram_counts(
    _id uuid DEFAULT uuid_generate_v4 (),
    language TEXT NOT NULL REFERENCES languages(code),
    corpus_id uuid NOT NULL REFERENCES corpora(_id),

    first_token_id uuid NOT NULL REFERENCES parts_of_speech(_id),
    second_token_id uuid NOT NULL REFERENCES parts_of_speech(_id),
    third_token_id uuid NOT NULL REFERENCES parts_of_speech(_id),
    occurrences NUMERIC NOT NULL,

    PRIMARY KEY (_id)
);