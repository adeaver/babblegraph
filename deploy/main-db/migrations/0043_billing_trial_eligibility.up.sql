CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS billing_newsletter_subscription_trials(
    _id uuid DEFAULT uuid_generate_v4 (),
    email_address TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT timezone('utc', now()),

    PRIMARY KEY(_id)
);

CREATE INDEX IF NOT EXISTS billing_newsletter_trial_email_address ON billing_newsletter_subscription_trials(email_address);
