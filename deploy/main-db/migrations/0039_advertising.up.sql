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

CREATE TABLE IF NOT EXISTS advertising_sources(
    _id uuid DEFAULT uuid_generate_v4 (),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT timezone('utc', now()),
    last_modified_at TIMESTAMP WITH TIME ZONE DEFAULT timezone('utc', now()),
    is_active BOOLEAN DEFAULT false,
    name TEXT NOT NULL,
    url TEXT NOT NULL,
    type TEXT NOT NULL,

    PRIMARY KEY(_id)
);

CREATE UNIQUE INDEX IF NOT EXISTS advertising_sources_url_unique_idx ON advertising_sources(url);

CREATE TABLE IF NOT EXISTS advertising_campaigns(
    _id uuid DEFAULT uuid_generate_v4 (),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT timezone('utc', now()),
    last_modified_at TIMESTAMP WITH TIME ZONE DEFAULT timezone('utc', now()),
    is_active BOOLEAN DEFAULT false,
    should_apply_to_all_users BOOLEAN DEFAULT false,
    name TEXT NOT NULL,
    url TEXT NOT NULL,
    vendor_id uuid NOT NULL REFERENCES advertising_vendors(_id),
    advertising_sources uuid NOT NULL REFERENCES advertising_sources(_id),

    PRIMARY KEY (_id),
);

CREATE INDEX IF NOT EXISTS advertising_campaigns_by_vendor_idx ON advertising_campaigns(vendor_id);
CREATE INDEX IF NOT EXISTS advertising_campaigns_by_source_idx ON advertising_campaigns(source_id);
