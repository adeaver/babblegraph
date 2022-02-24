CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE TABLE IF NOT EXISTS billing_external_id_mapping(
    _id uuid DEFAULT uuid_generate_v4 (),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT timezone('utc', now()),
    last_modified_at TIMESTAMP WITH TIME ZONE DEFAULT timezone('utc', now()),
    id_type TEXT NOT NULL,
    external_id TEXT NOT NULL,

    PRIMARY KEY(_id)
);

CREATE TABLE IF NOT EXISTS billing_information(
    _id uuid DEFAULT uuid_generate_v4 (),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT timezone('utc', now()),
    last_modified_at TIMESTAMP WITH TIME ZONE DEFAULT timezone('utc', now()),
    user_id uuid REFERENCES users(_id),
    external_id_mapping_id uuid NOT NULL REFERENCES billing_external_id_mapping(_id),

    PRIMARY KEY(_id)
);

CREATE UNIQUE INDEX IF NOT EXISTS billing_information_unique_user_idx ON billing_information(user_id) WHERE user_id IS NOT NULL;

CREATE TABLE IF NOT EXISTS billing_premium_newsletter_subscription(
    _id uuid DEFAULT uuid_generate_v4 (),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT timezone('utc', now()),
    last_modified_at TIMESTAMP WITH TIME ZONE DEFAULT timezone('utc', now()),
    billing_information_id UUID NOT NULL REFERENCES billing_information(_id),
    external_id_mapping_id uuid NOT NULL REFERENCES billing_external_id_mapping(_id),
    is_terminated BOOLEAN NOT NULL DEFAULT false,

    PRIMARY KEY(_id)
);

CREATE UNIQUE INDEX IF NOT EXISTS billing_premium_newsletter_subscription_active ON billing_premium_newsletter_subscription(billing_information_id) WHERE is_terminated = FALSE;

CREATE TABLE IF NOT EXISTS billing_premium_newsletter_subscription_debounce_record(
    created_at TIMESTAMP WITH TIME ZONE DEFAULT timezone('utc', now()),
    last_modified_at TIMESTAMP WITH TIME ZONE DEFAULT timezone('utc', now()),
    billing_information_id UUID NOT NULL REFERENCES billing_information(_id),

    PRIMARY KEY(billing_information_id) -- get free unique index here
);

CREATE TABLE IF NOT EXISTS billing_premium_newsletter_sync_request(
    created_at TIMESTAMP WITH TIME ZONE DEFAULT timezone('utc', now()),
    last_modified_at TIMESTAMP WITH TIME ZONE DEFAULT timezone('utc', now()),
    premium_newsletter_subscription_id UUID NOT NULL,
    update_type TEXT NOT NULL,
    attempt_number NUMERIC NOT NULL DEFAULT 0,
    hold_until TIMESTAMP WITH TIME ZONE DEFAULT timezone('utc', now()),

    PRIMARY KEY(premium_newsletter_subscription_id)
);

CREATE TABLE IF NOT EXISTS billing_stripe_event(
    _id uuid DEFAULT uuid_generate_v4 (),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT timezone('utc', now()),
    last_modified_at TIMESTAMP WITH TIME ZONE DEFAULT timezone('utc', now()),
    type TEXT NOT NULL,
    processed BOOLEAN NOT NULL DEFAULT FALSE,
    data JSONB NOT NULL,

    PRIMARY KEY(_id)
);
