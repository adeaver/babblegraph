CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS user_lemma_reinforcement_preferences(
    _id uuid DEFAULT uuid_generate_v4 (),
    user_id uuid NOT NULL REFERENCES users(_id),
    language_code TEXT NOT NULL,
    should_include_lemma_reinforcement BOOLEAN NOT NULL DEFAULT true,

    PRIMARY KEY (_id)
);

CREATE UNIQUE INDEX IF NOT EXISTS user_lemma_reinforcement_preferences_user_language ON user_lemma_reinforcement_preferences(user_id, language_code);
