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
    source_id uuid NOT NULL REFERENCES advertising_sources(_id),
    expires_at TIMESTAMP WITH TIME ZONE,

    PRIMARY KEY (_id)
);

CREATE INDEX IF NOT EXISTS advertising_campaigns_by_vendor_idx ON advertising_campaigns(vendor_id);
CREATE INDEX IF NOT EXISTS advertising_campaigns_by_source_idx ON advertising_campaigns(source_id);

CREATE TABLE IF NOT EXISTS advertising_campaign_topic_mappings(
    _id uuid DEFAULT uuid_generate_v4 (),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT timezone('utc', now()),
    last_modified_at TIMESTAMP WITH TIME ZONE DEFAULT timezone('utc', now()),
    campaign_id uuid NOT NULL REFERENCES advertising_campaigns(_id),
    topic_id uuid NOT NULL REFERENCES content_topic(_id),
    is_active BOOLEAN DEFAULT false,

    PRIMARY KEY(_id)
);

CREATE UNIQUE INDEX IF NOT EXISTS advertising_campaign_topic_mappings_unique_idx ON advertising_campaign_topic_mappings(campaign_id, topic_id);

CREATE TABLE IF NOT EXISTS advertising_advertisements(
    _id uuid DEFAULT uuid_generate_v4 (),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT timezone('utc', now()),
    last_modified_at TIMESTAMP WITH TIME ZONE DEFAULT timezone('utc', now()),
    campaign_id uuid NOT NULL REFERENCES advertising_campaigns(_id),
    language_code TEXT NOT NULL,
    title TEXT NOT NULL,
    description TEXT NOT NULL,
    image_url TEXT NOT NULL,
    is_active BOOLEAN DEFAULT false,

    PRIMARY KEY (_id)
);

CREATE INDEX IF NOT EXISTS advertising_advertisement_campaign_id_idx ON advertising_advertisement(campaign_id);

CREATE TABLE IF NOT EXISTS advertising_user_advertisements(
    _id uuid DEFAULT uuid_generate_v4 (),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT timezone('utc', now()),
    last_modified_at TIMESTAMP WITH TIME ZONE DEFAULT timezone('utc', now()),
    campaign_id uuid NOT NULL REFERENCES advertising_campaigns(_id),
    user_id uuid NOT NULL REFERENCES users(_id),
    advertisement_id uuid NOT NULL REFERENCES advertising_advertisements(_id),
    email_record_id TEXT NOT NULL REFERENCES email_records(_id),

    PRIMARY KEY (_id)
);

CREATE INDEX IF NOT EXISTS advertising_user_advertisements_user_idx ON advertising_user_advertisements(user_id);

CREATE TABLE IF NOT EXISTS advertising_advertisement_clicks(
    _id uuid DEFAULT uuid_generate_v4 (),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT timezone('utc', now()),
    last_modified_at TIMESTAMP WITH TIME ZONE DEFAULT timezone('utc', now()),
    user_advertisement_id uuid NOT NULL REFERENCES advertising_user_advertisements(_id),

    PRIMARY KEY (_id)
);
