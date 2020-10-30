CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS html_pages(
    _id uuid DEFAULT uuid_generate_v4 (),
    language TEXT,
    url TEXT NOT NULL,
    metadata jsonb NOT NULL,

    PRIMARY KEY (_id)
);

CREATE UNIQUE INDEX IF NOT EXISTS html_pages_url ON html_pages(url);