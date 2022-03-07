CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS content_podcast_metadata(
    _id uuid DEFAULT uuid_generate_v4 (),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT timezone('utc', now()),
    last_modified_at TIMESTAMP WITH TIME ZONE DEFAULT timezone('utc', now()),
    content_id uuid NOT NULL REFERENCES content_source(_id),
    image_url TEXT,

    PRIMARY KEY(_id)
);

CREATE TABLE IF NOT EXISTS user_podcast_preferences(
    _id uuid DEFAULT uuid_generate_v4 (),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT timezone('utc', now()),
    last_modified_at TIMESTAMP WITH TIME ZONE DEFAULT timezone('utc', now()),
    language_code TEXT NOT NULL,
    user_id uuid NOT NULL REFERENCES users(_id),
    podcasts_enabled BOOLEAN NOT NULL,
    include_explicit_podcasts BOOLEAN NOT NULL,
    minimum_duration_nanoseconds NUMERIC,
    maximum_duration_nanoseconds NUMERIC,

    PRIMARY KEY(_id)
);

CREATE UNIQUE INDEX IF NOT EXISTS user_podcast_preferences_user_language_idx ON user_podcast_preferences(language_code, user_id);

CREATE TABLE IF NOT EXISTS user_podcast_source_preferences(
    _id uuid DEFAULT uuid_generate_v4 (),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT timezone('utc', now()),
    last_modified_at TIMESTAMP WITH TIME ZONE DEFAULT timezone('utc', now()),
    language_code TEXT NOT NULL,
    user_id uuid NOT NULL REFERENCES users(_id),
    source_id uuid NOT NULL REFERENCES content_source(_id),
    is_active BOOLEAN NOT NULL,

    PRIMARY KEY(_id)
);

CREATE UNIQUE INDEX IF NOT EXISTS user_podcast_source_preferences_user_language_idx ON user_podcast_source_preferences(user_id, language_code, source_id);

CREATE TABLE IF NOT EXISTS user_podcasts(
    _id uuid DEFAULT uuid_generate_v4 (),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT timezone('utc', now()),
    last_modified_at TIMESTAMP WITH TIME ZONE DEFAULT timezone('utc', now()),
    email_record_id TEXT NOT NULL REFERENCES email_records(_id),
    user_id uuid NOT NULL REFERENCES users(_id),
    episode_id TEXT NOT NULL,
    source_id uuid NOT NULL REFERENCES content_source(_id),
    first_opened_at TIMESTAMP WITH TIME ZONE,

    PRIMARY KEY(_id)
);

CREATE UNIQUE INDEX IF NOT EXISTS user_podcasts_unique_idx ON user_podcasts(user_id, episode_id, email_record_id);
