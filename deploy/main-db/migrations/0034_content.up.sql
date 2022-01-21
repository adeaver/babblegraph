CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS content_topic(
    _id uuid DEFAULT uuid_generate_v4 (),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT timezone('utc', now()),
    last_modified_at TIMESTAMP WITH TIME ZONE DEFAULT timezone('utc', now()),
    label TEXT NOT NULL,
    is_active BOOLEAN NOT NULL,

    PRIMARY KEY (_id)
);

CREATE UNIQUE INDEX IF NOT EXISTS content_topic_label_unique ON content_topic(label);

CREATE TABLE IF NOT EXISTS content_topic_display_name(
    _id uuid DEFAULT uuid_generate_v4 (),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT timezone('utc', now()),
    last_modified_at TIMESTAMP WITH TIME ZONE DEFAULT timezone('utc', now()),
    topic_id uuid NOT NULL REFERENCES content_topic(_id),
    language_code TEXT NOT NULL,
    label TEXT NOT NULL,
    is_active BOOLEAN NOT NULL,

    PRIMARY KEY (_id)
);

CREATE UNIQUE INDEX IF NOT EXISTS content_topic_display_name_for_language ON content_topic_display_name(topic_id, language_code);

CREATE TABLE IF NOT EXISTS content_source(
    _id uuid DEFAULT uuid_generate_v4 (),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT timezone('utc', now()),
    last_modified_at TIMESTAMP WITH TIME ZONE DEFAULT timezone('utc', now()),
    title TEXT NOT NULL,
    url TEXT NOT NULL,
    type TEXT NOT NULL,
    country TEXT NOT NULL,
    ingest_strategy TEXT NOT NULL,
    language_code TEXT NOT NULL,
    is_active BOOLEAN NOT NULL,
    should_use_url_as_seed_url BOOLEAN NOT NULL DEFAULT false,
    monthly_access_limit INTEGER,

    PRIMARY KEY (_id)
);

CREATE UNIQUE INDEX IF NOT EXISTS content_source_url ON content_source(url);

CREATE TABLE IF NOT EXISTS content_source_topic_mapping(
    _id uuid DEFAULT uuid_generate_v4 (),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT timezone('utc', now()),
    last_modified_at TIMESTAMP WITH TIME ZONE DEFAULT timezone('utc', now()),
    root_id uuid REFERENCES content_source(_id),
    topic_id uuid REFERENCES content_topic(_id),
    is_active BOOLEAN NOT NULL,

    PRIMARY KEY (_id)
);

CREATE TABLE IF NOT EXISTS content_source_seed(
    _id uuid DEFAULT uuid_generate_v4 (),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT timezone('utc', now()),
    last_modified_at TIMESTAMP WITH TIME ZONE DEFAULT timezone('utc', now()),
    root_id uuid NOT NULL REFERENCES content_source(_id),
    url TEXT NOT NULL,
    is_active BOOLEAN NOT NULL,

    PRIMARY KEY (_id)
);

CREATE INDEX IF NOT EXISTS content_source_seed_root_idx ON content_source_seed(root_id);

CREATE TABLE IF NOT EXISTS content_source_filter(
    _id uuid DEFAULT uuid_generate_v4 (),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT timezone('utc', now()),
    last_modified_at TIMESTAMP WITH TIME ZONE DEFAULT timezone('utc', now()),
    root_id uuid NOT NULL REFERENCES content_source(_id),
    url TEXT NOT NULL,
    is_active BOOLEAN NOT NULL,
    use_ld_json_validation BOOLEAN,
    paywall_classes TEXT,
    paywall_ids TEXT,

    PRIMARY KEY (_id)
);

CREATE UNIQUE INDEX IF NOT EXISTS content_source_filter_root_id_unique ON content_source_filter(root_id);

CREATE TABLE IF NOT EXISTS content_source_seed_topic_mapping(
    _id uuid DEFAULT uuid_generate_v4 (),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT timezone('utc', now()),
    last_modified_at TIMESTAMP WITH TIME ZONE DEFAULT timezone('utc', now()),
    topic_id uuid NOT NULL REFERENCES content_topic(_id),
    source_seed_id uuid NOT NULL REFERENCES content_source_seed(_id),
    is_active BOOLEAN NOT NULL,

    PRIMARY KEY (_id)
);

CREATE UNIQUE INDEX IF NOT EXISTS content_source_seed_topic_mapping_unique ON content_source_seed_topic_mapping(source_seed_id, topic_id);
