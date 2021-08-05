CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS user_account_notification_requests(
    _id uuid DEFAULT uuid_generate_v4 (),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT timezone('utc', now()),
    last_modified_at TIMESTAMP WITH TIME ZONE DEFAULT timezone('utc', now()),
    hold_until TIMESTAMP WITH TIME ZONE DEFAULT timezone('utc', now()),
    fulfilled_at TIMESTAMP WITH TIME ZONE,
    user_id uuid NOT NULL REFERENCES users(_id),
    notification_type TEXT NOT NULL,

    PRIMARY KEY (_id)
);

CREATE TABLE IF NOT EXISTS user_account_notification_request_debounce_fulfillment_records(
    _id uuid DEFAULT uuid_generate_v4 (),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT timezone('utc', now()),
    last_modified_at TIMESTAMP WITH TIME ZONE DEFAULT timezone('utc', now()),
    notification_request_id uuid NOT NULL REFERENCES user_account_notification_request(_id),

    PRIMARY KEY (_id)
);

CREATE UNIQUE INDEX IF NOT EXISTS user_account_request_debounce_idx  ON user_account_notification_request_debounce_fulfillment_record(notification_request_id);
