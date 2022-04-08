ALTER TABLE user_newsletter_schedule_day_metadata DROP COLUMN IF EXISTS content_topics;
ALTER TABLE user_newsletter_schedule_day_metadata DROP COLUMN IF EXISTS number_of_articles;
DROP TABLE IF EXISTS user_newsletter_schedule_day_topic_mapping;

ALTER TABLE user_newsletter_schedule ADD COLUMN IF NOT EXISTS  number_of_articles_per_email INTEGER NOT NULL DEFAULT 12;
