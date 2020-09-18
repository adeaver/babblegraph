CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE corpora(
    _id uuid DEFAULT uuid_generate_v4 (),
    language TEXT NOT NULL REFERENCES languages(code),
    name TEXT NOT NULL,

    PRIMARY KEY (_id)
);
