CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS paywall_reports(
    _id uuid DEFAULT uuid_generate_v4 (),
    user_id uuid NOT NULL REFERENCES users(_id),
    domain TEXT NOT NULL,
    url_identifier TEXT NOT NULL,
    email_record_id TEXT NOT NULL REFERENCES email_records(_id),
    access_month TEXT NOT NULL,

    PRIMARY KEY (_id)
);

CREATE UNIQUE INDEX IF NOT EXISTS user_paywalls_unique_url_month ON paywall_reports(user_id, url_identifier, access_month);
