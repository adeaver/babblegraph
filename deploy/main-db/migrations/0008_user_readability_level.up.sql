CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS user_readability_level(
    _id uuid DEFAULT uuid_generate_v4 (),
    user_id uuid NOT NULL REFERENCES users(_id) ON DELETE CASCADE,
    language_code TEXT NOT NULL,
    readability_level NUMERIC NOT NULL,
    version NUMERIC NOT NULL,

    PRIMARY KEY (_id)
);

CREATE INDEX IF NOT EXISTS user_readability_level_user_idx ON user_readability_level(user_id);
