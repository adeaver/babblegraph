CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS advertising_vendors(
    _id uuid DEFAULT uuid_generate_v4 (),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT timezone('utc', now()),
    last_modified_at TIMESTAMP WITH TIME ZONE DEFAULT timezone('utc', now()),
    is_active BOOLEAN DEFAULT false,
    name TEXT NOT NULL,
    website_url TEXT NOT NULL,

    PRIMARY KEY(_id)
);

CREATE UNIQUE INDEX IF NOT EXISTS advertising_vendors_website_url_unique_idx ON advertising_vendors(website_url);
