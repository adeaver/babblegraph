CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS email_records(
    _id TEXT NOT NULL,
    sent_at TIMESTAMP WITH TIME ZONE DEFAULT timezone('utc', now()),
    ses_message_id TEXT NOT NULL,
    user_id uuid NOT NULL REFERENCES users(_id),
    first_opened_at TIMESTAMP WITH TIME ZONE,
    type TEXT NOT NULL,

    PRIMARY KEY (_id)
);

CREATE INDEX IF NOT EXISTS email_records_user_id ON email_records(user_id);

ALTER TABLE user_documents ADD COLUMN IF NOT EXISTS email_id TEXT REFERENCES email_records(_id);
CREATE INDEX IF NOT EXISTS user_documents_email ON user_documents(email_id);
