CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS utm_page_hits(
    _id uuid DEFAULT uuid_generate_v4 (),
    source TEXT,
    medium TEXT,
    campaign_id TEXT,
    url_path TEXT,
    tracking_id TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT timezone('utc', now()),

    PRIMARY KEY(_id)
);

CREATE TABLE IF NOT EXISTS utm_events(
    _id uuid DEFAULT uuid_generate_v4 (),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT timezone('utc', now()),
    tracking_id TEXT NOT NULL,
    event_type TEXT NOT NULL,

    PRIMARY KEY(_id)
);
