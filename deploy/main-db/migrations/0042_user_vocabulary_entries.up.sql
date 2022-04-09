CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS user_vocabulary_entries(
    _id uuid DEFAULT uuid_generate_v4 (),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT timezone('utc', now()),
    last_modified_at TIMESTAMP WITH TIME ZONE DEFAULT timezone('utc', now()),
    is_active BOOLEAN DEFAULT false,
    is_visible BOOLEAN DEFAULT false,
    user_id uuid REFERENCES users(_id),
    language_code TEXT NOT NULL,
    vocabulary_id TEXT,
    vocabulary_type TEXT NOT NULL,
    vocabulary_display TEXT NOT NULL,
    study_note VARCHAR(500),
    unique_hash TEXT NOT NULL,

    PRIMARY KEY(_id)
);

CREATE UNIQUE INDEX IF NOT EXISTS user_vocabulary_entries_unique_idx ON user_vocabulary_entries(user_id, language_code, unique_hash);
