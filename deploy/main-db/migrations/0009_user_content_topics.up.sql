CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS user_content_topic_mappings(
    _id uuid DEFAULT uuid_generate_v4 (),
    user_id uuid NOT NULL REFERENCES users(_id),
    content_topic TEXT NOT NULL,
    is_active BOOLEAN NOT NULL DEFAULT true,

    PRIMARY KEY(_id)
);

CREATE UNIQUE INDEX IF NOT EXISTS user_content_topic_mappings_user_idx ON user_content_topic_mappings(user_id, content_topic);
