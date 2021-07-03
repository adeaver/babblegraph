CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS user_newsletter_schedule_day_metadata(
    _id uuid DEFAULT uuid_generate_v4 (),
    user_id uuid NOT NULL REFERENCES users(_id),
    day_of_week_index_utc INTEGER NOT NULL CHECK (day_of_week_index_utc >= 0 AND day_of_week_index_utc <= 6),
    hour_of_day_index_utc INTEGER NOT NULL CHECK (hour_of_day_index_utc >= 0 AND hour_of_day_index_utc <= 23),
    quarter_hour_index_utc INTEGER NOT NULL CHECK (quarter_hour_index_utc >= 0 AND quarter_hour_index_utc <= 3),
    language_code TEXT NOT NULL,
    content_topics TEXT,
    number_of_articles INTEGER NOT NULL CHECK (number_of_articles >= 4 AND number_of_articles <= 12),
    is_active BOOLEAN NOT NULL,

    PRIMARY KEY (_id)
);

CREATE UNIQUE INDEX IF NOT EXISTS user_newsletter_schedule_day_metadata_user_unique_day ON user_newsletter_schedule_day_metadata(user_id, language_code, day_of_week_index_utc);
