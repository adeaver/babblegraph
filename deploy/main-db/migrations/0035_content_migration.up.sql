CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

ALTER TABLE user_content_topic_mappings ADD COLUMN IF NOT EXISTS content_topic_id uuid REFERENCES content_topic(_id);

ALTER TABLE links2 ADD COLUMN IF NOT EXISTS source_id uuid REFERENCES content_source(_id);
ALTER TABLE content_topic_mappings ADD COLUMN IF NOT EXISTS topic_mapping_id TEXT;
