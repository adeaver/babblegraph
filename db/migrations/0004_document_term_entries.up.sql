CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE document_term_entries(
    _id uuid DEFAULT uuid_generate_v4 (),
    document_id uuid NOT NULL REFERENCES documents(_id),
    term TEXT NOT NULL,
    count NUMERIC NOT NULL,

    PRIMARY KEY (_id)
);

CREATE INDEX IF NOT EXISTS document_term_entries_term_idx ON document_term_entries(term);
