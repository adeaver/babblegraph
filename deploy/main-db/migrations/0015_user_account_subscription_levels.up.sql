CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS user_account_subscription_levels(
    _id uuid DEFAULT uuid_generate_v4 (),
    user_id uuid NOT NULL REFERENCES users(_id),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT timezone('utc', now()),
    subscription_level TEXT NOT NULL,
    expires_at TIMESTAMP WITH TIME ZONE NOT NULL,
    is_active BOOLEAN NOT NULL DEFAULT true,

    PRIMARY KEY(_id)
);

CREATE UNIQUE INDEX IF NOT EXISTS user_account_subscription_levels_user_idx ON user_account_subscription_levels(user_id);
