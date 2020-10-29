CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE languages(
    _id uuid DEFAULT uuid_generate_v4 (),
    code TEXT NOT NULL,

    PRIMARY KEY (_id)
);

CREATE UNIQUE INDEX IF NOT EXISTS languages_code_idx ON languages(code);
