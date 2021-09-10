CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS newsletter_send_requests(
    _id TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT timezone('utc', now()),
    last_modified_at TIMESTAMP WITH TIME ZONE DEFAULT timezone('utc', now()),
    user_id uuid NOT NULL REFERENCES users(_id),
	language_code TEXT NOT NULL,
    date_of_send TEXT NOT NULL,
	hour_to_send_index_utc INTEGER,
	quarter_hour_to_send_index_utc INTEGER,
	payload_status TEXT NOT NULL,

    PRIMARY KEY (_id)
);

CREATE TABLE IF NOT EXISTS newsletter_send_request_debounce_records(
    _id uuid DEFAULT uuid_generate_v4 (),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT timezone('utc', now()),
    last_modified_at TIMESTAMP WITH TIME ZONE DEFAULT timezone('utc', now()),
    newsletter_send_request_id TEXT NOT NULL REFERENCES newsletter_send_requests(_id),
	to_payload_status TEXT NOT NULL,

    PRIMARY KEY (_id)
);

CREATE UNIQUE INDEX IF NOT EXISTS unique_newsletter_send_request_debounce_records ON newsletter_send_request_debounce_records(newsletter_send_request_id, to_payload_status);
