CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE documents(
    _id uuid DEFAULT uuid_generate_v4 (),
    url TEXT NOT NULL,
    language TEXT,
    metadata JSONB,

    PRIMARY KEY (_id)
);

CREATE UNIQUE INDEX IF NOT EXISTS documents_url ON documents(url);
