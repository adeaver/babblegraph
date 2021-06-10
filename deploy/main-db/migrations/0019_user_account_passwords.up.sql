CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS user_account_passwords(
    _id uuid DEFAULT uuid_generate_v4 (),
    user_id uuid NOT NULL REFERENCES users(_id),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT timezone('utc', now()),
    password_hash TEXT NOT NULL,
    salt TEXT NOT NULL,

    PRIMARY KEY(_id)
);

CREATE UNIQUE INDEX IF NOT EXISTS user_account_passwords_user_idx ON user_account_passwords(user_id);
