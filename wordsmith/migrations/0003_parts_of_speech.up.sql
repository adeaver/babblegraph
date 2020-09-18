CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE parts_of_speech(
    _id uuid DEFAULT uuid_generate_v4 (),
    language TEXT NOT NULL REFERENCES languages(code),
    code TEXT NOT NULL,

    PRIMARY KEY (_id)
)
