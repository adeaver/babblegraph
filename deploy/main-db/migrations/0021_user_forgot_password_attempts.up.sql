CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS user_forgot_password_attempts(
    _id uuid DEFAULT uuid_generate_v4 (),
    user_id uuid NOT NULL REFERENCES users(_id),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT timezone('utc', now()),
    fulfilled_at TIMESTAMP WITH TIME ZONE,
    is_archived BOOLEAN NOT NULL DEFAULT false,
    has_been_used BOOLEAN NOT NULL DEFAULT false,

    PRIMARY KEY(_id)
);

CREATE UNIQUE INDEX IF NOT EXISTS user_forgot_password_attempts_unfulfilled_attempt ON user_forgot_password_attempts(user_id) WHERE fulfilled_at IS NULL;
