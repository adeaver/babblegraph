CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS blog_post_metadata(
    _id uuid DEFAULT uuid_generate_v4 (),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT timezone('utc', now()),
    last_modified_at TIMESTAMP WITH TIME ZONE DEFAULT timezone('utc', now()),
    published_at TIMESTAMP WITH TIME ZONE,
    title TEXT NOT NULL,
    author_name TEXT NOT NULL,
    description TEXT NOT NULL,
    url_path TEXT NOT NULL,
    status TEXT NOT NULL,

    PRIMARY KEY (_id)
);

CREATE UNIQUE INDEX IF NOT EXISTS blog_post_metadata_url_path ON blog_post_metadata(url_path);

CREATE TABLE IF NOT EXISTS blog_post_image_metadata(
    _id uuid DEFAULT uuid_generate_v4 (),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT timezone('utc', now()),
    last_modified_at TIMESTAMP WITH TIME ZONE DEFAULT timezone('utc', now()),
    path TEXT NOT NULL,
    blog_id uuid NOT NULL REFERENCES blog_post_metadata(_id),
    file_name TEXT NOT NULL,
    alt_text TEXT NOT NULL,
    is_hero_image BOOLEAN DEFAULT FALSE,
    caption TEXT,

    PRIMARY KEY (_id)
);

CREATE UNIQUE INDEX IF NOT EXISTS blog_post_image_file_name ON blog_post_image_metadata(blog_id, file_name);
CREATE UNIQUE INDEX IF NOT EXISTS blog_post_image_hero_image ON blog_post_image_metadata(blog_id) WHERE is_hero_image = TRUE;

CREATE TABLE IF NOT EXISTS blog_post_view(
    _id uuid DEFAULT uuid_generate_v4 (),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT timezone('utc', now()),
    last_modified_at TIMESTAMP WITH TIME ZONE DEFAULT timezone('utc', now()),
    blog_id uuid NOT NULL REFERENCES blog_post_metadata(_id),
    tracking_id TEXT,

    PRIMARY KEY (_id)
);

CREATE INDEX IF NOT EXISTS blog_post_view_blog_id ON blog_post_view(blog_id);
