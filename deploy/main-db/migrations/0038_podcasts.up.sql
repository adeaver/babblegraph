CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS content_podcast_metadata(
    _id uuid DEFAULT uuid_generate_v4 (),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT timezone('utc', now()),
    last_modified_at TIMESTAMP WITH TIME ZONE NOT NULL,
    content_id uuid NOT NULL REFERENCES content_source(_id),
    image_url TEXT,

    PRIMARY KEY(_id)
);
