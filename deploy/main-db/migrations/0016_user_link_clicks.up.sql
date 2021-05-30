CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS user_link_clicks(
    _id uuid DEFAULT uuid_generate_v4 (),
    user_id uuid NOT NULL REFERENCES users(_id),
    domain TEXT NOT NULL,
    url_identifier TEXT NOT NULL,
    email_record_id TEXT NOT NULL REFERENCES email_records(_id),
    access_month TEXT NOT NULL,
    first_accessed_at TIMESTAMP WITH TIME ZONE DEFAULT timezone('utc', now()),

    PRIMARY KEY (_id)
);

CREATE UNIQUE INDEX IF NOT EXISTS user_link_clicks_unique_url_month ON user_link_clicks(user_id, url_identifier, access_month);
CREATE INDEX IF NOT EXISTS user_link_clicks_user_access_month ON user_link_clicks(user_id, access_month);
