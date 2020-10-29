CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE words(
    _id uuid DEFAULT uuid_generate_v4 (),
    language TEXT NOT NULL REFERENCES languages(code),
    corpus_id uuid NOT NULL REFERENCES corpora(_id),

    part_of_speech_id uuid NOT NULL REFERENCES parts_of_speech(_id),
    lemma_id uuid NOT NULL REFERENCES lemmas(_id) ON DELETE CASCADE,
    word_text TEXT NOT NULL,

    PRIMARY KEY (_id)
);
