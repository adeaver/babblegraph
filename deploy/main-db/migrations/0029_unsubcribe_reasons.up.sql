CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS unsubscribe_reasons(
    _id uuid DEFAULT uuid_generate_v4 (),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT timezone('utc', now()),
    last_modified_at TIMESTAMP WITH TIME ZONE DEFAULT timezone('utc', now()),
    user_id uuid NOT NULL REFERENCES users(_id),
	language_code TEXT NOT NULL,
    reason VARCHAR(500) NOT NULL,

    PRIMARY KEY (_id)
);

CREATE INDEX IF NOT EXISTS unsubscribe_reason_user_id ON unsubscribe_reasons(user_id);
