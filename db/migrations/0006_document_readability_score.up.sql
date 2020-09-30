CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE document_reability_score(
    _id uuid DEFAULT uuid_generate_v4 (),
    document_id uuid NOT NULL REFERENCES documents(_id),
    readbility_score NUMERIC NOT NULL,

    PRIMARY KEY (_id)
);

CREATE UNIQUE INDEX IF NOT EXISTS document_reability_score_unique_documents_idx ON document_reability_score(document_id);
