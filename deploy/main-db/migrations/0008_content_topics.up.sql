CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS content_topic_mappings(
    _id uuid DEFAULT uuid_generate_v4 (),
    url_identifier TEXT NOT NULL,
    content_topic TEXT NOT NULL
);

CREATE UNIQUE INDEX IF NOT EXISTS content_topic_mapping_uniqueness_idx ON content_topic_mappings(url_identifier, content_topic);
