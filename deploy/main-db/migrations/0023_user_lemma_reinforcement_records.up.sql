CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS user_lemma_reinforcement_records(
    _id uuid DEFAULT uuid_generate_v4 (),
    user_id uuid NOT NULL REFERENCES users(_id),
    lemma_id TEXT NOT NULL,
    language_code TEXT NOT NULL,
    last_sent_on TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT timezone('utc', now()),
	number_of_times_sent INTEGER NOT NULL,

    PRIMARY KEY (_id)
);

CREATE UNIQUE INDEX IF NOT EXISTS user_lemma_reinforcement_records_user_lemma ON user_lemma_reinforcement_records(user_id, lemma_id);
