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

CREATE TABLE IF NOT EXISTS user_vocabulary_spotlight_records(
    created_at TIMESTAMP WITH TIME ZONE DEFAULT timezone('utc', now()),
    last_modified_at TIMESTAMP WITH TIME ZONE DEFAULT timezone('utc', now()),
    _id uuid DEFAULT uuid_generate_v4 (),
    user_id uuid NOT NULL REFERENCES users(_id),
    vocabulary_entry_id uuid NOT NULL REFERENCES user_vocabulary_entries(_id),
    language_code TEXT NOT NULL,
    last_sent_on TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT timezone('utc', now()),
	number_of_times_sent INTEGER NOT NULL,

    PRIMARY KEY (_id)
);

CREATE UNIQUE INDEX IF NOT EXISTS user_vocabulary_spotlight_records_entry_unique_idx ON user_vocabulary_spotlight_records(user_id, language_code, vocabulary_entry_id);
