CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS blog_post_view(
    _id uuid DEFAULT uuid_generate_v4 (),
    blog_post_id uuid NOT NULL REFERENCES blog_posts(_id),
    viewed_at TIMESTAMP WITH TIME ZONE DEFAULT timezone('utc', now()),

    PRIMARY KEY (_id)
);

CREATE INDEX IF NOT EXISTS blog_post_views_blog_post_index ON blog_post_view(blog_post_id);
