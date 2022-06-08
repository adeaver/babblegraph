CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS user_reader_tutorial_receipt(
    created_at TIMESTAMP WITH TIME ZONE DEFAULT timezone('utc', now()),
    last_modified_at TIMESTAMP WITH TIME ZONE DEFAULT timezone('utc', now()),
    _id uuid DEFAULT uuid_generate_v4 (),
    user_id uuid NOT NULL REFERENCES users(_id),

    PRIMARY KEY (_id)
);

CREATE UNIQUE INDEX IF NOT EXISTS user_reader_tutorial_receipt_user_idx ON user_reader_tutorial_receipt(user_id);
