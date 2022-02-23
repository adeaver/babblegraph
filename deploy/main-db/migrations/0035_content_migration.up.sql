CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

ALTER TABLE user_content_topic_mappings ADD COLUMN IF NOT EXISTS content_topic_id uuid REFERENCES content_topic(_id);

ALTER TABLE links2 ADD COLUMN IF NOT EXISTS source_id uuid REFERENCES content_source(_id);
ALTER TABLE content_topic_mappings ADD COLUMN IF NOT EXISTS topic_mapping_id TEXT;
ALTER TABLE user_link_clicks ADD COLUMN IF NOT EXISTS source_id uuid REFERENCES content_source(_id);

CREATE TABLE IF NOT EXISTS user_newsletter_schedule_day_topic_mapping(
    _id uuid DEFAULT uuid_generate_v4 (),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT timezone('utc', now()),
    last_modified_at TIMESTAMP WITH TIME ZONE DEFAULT timezone('utc', now()),
    topic_id uuid NOT NULL REFERENCES content_topic(_id),
    day_id uuid NOT NULL REFERENCES user_newsletter_schedule_day_metadata(_id),
    is_active BOOLEAN NOT NULL,

    PRIMARY KEY (_id)
);

CREATE UNIQUE INDEX IF NOT EXISTS user_newsletter_schedule_day_mapping_unique_day_and_topic_idx ON user_newsletter_schedule_day_topic_mapping(day_id, topic_id);
