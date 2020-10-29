CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS links(
    _id uuid DEFAULT uuid_generate_v4 (),
    domain TEXT NOT NULL,
    url TEXT NOT NULL,
    has_fetched BOOLEAN DEFAULT false,
    position SERIAL NOT NULL,

    PRIMARY KEY (_id)
);

CREATE UNIQUE INDEX IF NOT EXISTS links_url_idx ON links(url);
CREATE INDEX IF NOT EXISTS links_domain ON links(domain) WHERE has_fetched = FALSE;
