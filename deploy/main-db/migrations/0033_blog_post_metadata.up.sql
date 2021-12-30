CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS blog_post_metadata(
    _id uuid DEFAULT uuid_generate_v4 (),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT timezone('utc', now()),
    last_modified_at TIMESTAMP WITH TIME ZONE DEFAULT timezone('utc', now()),
    published_at TIMESTAMP WITH TIME ZONE,
    hero_image_path TEXT,
    title TEXT NOT NULL,
    author_name TEXT NOT NULL,
    description TEXT NOT NULL,
    url_path TEXT NOT NULL,
    status TEXT NOT NULL,

    PRIMARY KEY (_id)
);

CREATE UNIQUE INDEX IF NOT EXISTS blog_post_metadata_url_path ON blog_post_metadata(url_path);
