CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE user_documents(
    _id uuid DEFAULT uuid_generate_v4 (),
    user_id uuid NOT NULL REFERENCES users(_id),
    document_id TEXT NOT NULL,
    sent_on TIMESTAMP WITH TIME ZONE DEFAULT timezone('utc', now()),

    PRIMARY KEY (_id)
);

CREATE UNIQUE INDEX IF NOT EXISTS user_document_unique_idx ON user_documents(user_id, document_id);
