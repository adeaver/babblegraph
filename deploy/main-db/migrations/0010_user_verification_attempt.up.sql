CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS user_verification_attempts(
    _id uuid DEFAULT uuid_generate_v4 (),
    user_id uuid NOT NULL REFERENCES users(_id),
    fulfilled_at_timestamp TIMESTAMP WITH TIME ZONE,

    PRIMARY KEY(_id)
);

CREATE UNIQUE INDEX IF NOT EXISTS user_verification_attempts_unfulfilled_attempt ON user_verification_attempts(user_id) WHERE fulfilled_at_timestamp IS NULL;
