CREATE TABLE IF NOT EXISTS links2(
    url_identifier TEXT NOT NULL,
    domain TEXT NOT NULL,
    url TEXT NOT NULL,
    last_fetch_version INT,
    fetched_on TIMESTAMP WITH TIME ZONE DEFAULT timezone('utc', now()),
    seq_num SERIAL NOT NULL,

    PRIMARY KEY (url_identifier)
);

--  CREATE INDEX IF NOT EXISTS links2_domain ON links2(domain) WHERE last_fetch_version = NULL;
