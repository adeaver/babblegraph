CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS bgstripe_customer(
    _id uuid DEFAULT uuid_generate_v4 (),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT timezone('utc', now()),
    last_modified_at TIMESTAMP WITH TIME ZONE DEFAULT timezone('utc', now()),
    babblegraph_user_id uuid NOT NULL REFERENCES users(_id),
    stripe_customer_id TEXT NOT NULL,
    default_payment_method TEXT,

    PRIMARY KEY (_id)
);

CREATE UNIQUE INDEX IF NOT EXISTS bgstripe_customer_babblegraph_user_unique ON bgstripe_customer(babblegraph_user_id);

CREATE UNIQUE INDEX IF NOT EXISTS bgstripe_customer_stripe_id_unique ON bgstripe_customer(stripe_customer_id);

CREATE TABLE IF NOT EXISTS bgstripe_subscription(
    _id uuid DEFAULT uuid_generate_v4 (),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT timezone('utc', now()),
    last_modified_at TIMESTAMP WITH TIME ZONE DEFAULT timezone('utc', now()),
    babblegraph_user_id uuid NOT NULL REFERENCES users(_id),
    stripe_subscription_id TEXT NOT NULL,
    payment_state INTEGER NOT NULL,
    stripe_product_id TEXT NOT NULL,

    PRIMARY KEY (_id)
);

CREATE UNIQUE INDEX IF NOT EXISTS bgstripe_subscription_babblegraph_user ON bgstripe_subscription(babblegraph_user_id) WHERE payment_state != 5;

CREATE UNIQUE INDEX IF NOT EXISTS bgstripe_subscription_stripe_id_unique ON bgstripe_subscription(stripe_subscription_id);

CREATE TABLE IF NOT EXISTS bgstripe_payment_method(
    _id uuid DEFAULT uuid_generate_v4 (),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT timezone('utc', now()),
    last_modified_at TIMESTAMP WITH TIME ZONE DEFAULT timezone('utc', now()),
    babblegraph_user_id uuid NOT NULL REFERENCES users(_id),
	stripe_payment_method_id TEXT NOT NULL,
	card_type TEXT NOT NULL,
	last_four_digits TEXT NOT NULL,
	expiration_month TEXT NOT NULL,
	expiration_year TEXT NOT NULL,

    PRIMARY KEY (_id)
);

CREATE INDEX IF NOT EXISTS bgstripe_payment_method_user ON bgstripe_payment_method(babblegraph_user_id);

CREATE UNIQUE INDEX IF NOT EXISTS bgstripe_payment_method_stripe_id_unique ON bgstripe_payment_method(stripe_payment_method_id);
