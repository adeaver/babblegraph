CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS item_cache(
    created_at TIMESTAMP WITH TIME ZONE DEFAULT timezone('utc', now()),
    expires_at TIMESTAMP WITH TIME ZONE NOT NULL,
    key TEXT NOT NULL,
    item JSONB NOT NULL,

    PRIMARY KEY(key)
);
