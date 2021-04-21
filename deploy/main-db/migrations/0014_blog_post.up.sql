CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS blog_posts(
    _id uuid DEFAULT uuid_generate_v4 (),
    title TEXT NOT NULL,
    description TEXT NOT NULL,
    tags TEXT,
    tracking_tag TEXT NOT NULL,
    url_path TEXT NOT NULL, -- This is the suffix for the url
    hero_image_url TEXT NOT NULL, -- This is the suffix for the url
    first_published_date TIMESTAMP WITH TIME ZONE DEFAULT timezone('utc', now()),
    updated_date TIMESTAMP WITH TIME ZONE,
    is_visible BOOLEAN NOT NULL DEFAULT true,

    PRIMARY KEY (_id)
);

CREATE INDEX IF NOT EXISTS blog_posts_url_path ON blog_posts(url_path);
