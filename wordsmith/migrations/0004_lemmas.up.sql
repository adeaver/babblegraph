CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE lemmas(
    _id uuid DEFAULT uuid_generate_v4 (),
    lemma TEXT NOT NULL,
    part_of_speech uuid NOT NULL REFERENCES parts_of_speech(_id),
    language TEXT NOT NULL REFERENCES languages(code),

    PRIMARY KEY (_id)
);
