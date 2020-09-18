CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE words(
    _id uuid DEFAULT uuid_generate_v4 (),
    word_text TEXT NOT NULL,
    lemma_id uuid NOT NULL REFERENCES lemmas(_id) ON DELETE CASCADE,
    part_of_speech uuid NOT NULL REFERENCES parts_of_speech(_id),
    language TEXT NOT NULL REFERENCES languages(code),

    PRIMARY KEY (_id)
);
