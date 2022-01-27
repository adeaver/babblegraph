CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

ALTER TABLE user_content_topic_mappings ADD COLUMN IF NOT EXISTS content_topic_id uuid REFERENCES content_topic(_id);
