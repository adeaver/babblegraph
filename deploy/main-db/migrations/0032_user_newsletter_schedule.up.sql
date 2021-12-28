CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

ALTER TABLE user_newsletter_schedule_day_metadata RENAME COLUMN day_of_week_index_utc TO day_of_week_index;
ALTER TABLE user_newsletter_schedule_day_metadata DROP COLUMN hour_of_day_index_utc;
ALTER TABLE user_newsletter_schedule_day_metadata DROP COLUMN quarter_hour_index_utc;

CREATE TABLE IF NOT EXISTS user_newsletter_schedule(
    _id uuid DEFAULT uuid_generate_v4 (),
    user_id uuid NOT NULL REFERENCES users(_id),
    language_code TEXT NOT NULL,
    iana_timezone TEXT NOT NULL,
    hour_of_day_index INTEGER NOT NULL CHECK (hour_of_day_index >= 0 AND hour_of_day_index <= 23),
    quarter_hour_index INTEGER NOT NULL CHECK (quarter_hour_index >= 0 AND quarter_hour_index <= 3),

    PRIMARY KEY (_id)
);

CREATE UNIQUE INDEX IF NOT EXISTS user_newsletter_schedule_language ON user_newsletter_schedule(user_id, language_code);
