CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS admin_user(
    _id uuid DEFAULT uuid_generate_v4 (),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT timezone('utc', now()),
    last_modified_at TIMESTAMP WITH TIME ZONE DEFAULT timezone('utc', now()),
    email_address TEXT NOT NULL,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,

    PRIMARY KEY (_id)
);

CREATE INDEX IF NOT EXISTS admin_user_email_address_idx ON admin_user(email_address);

CREATE TABLE IF NOT EXISTS admin_user_password(
    _id uuid DEFAULT uuid_generate_v4 (),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT timezone('utc', now()),
    last_modified_at TIMESTAMP WITH TIME ZONE DEFAULT timezone('utc', now()),
    admin_user_id uuid NOT NULL REFERENCES admin_user(_id),
    password_hash TEXT NOT NULL,
    salt TEXT NOT NULL,
    is_active BOOLEAN NOT NULL DEFAULT FALSE,

    PRIMARY KEY (_id)
);

CREATE UNIQUE INDEX IF NOT EXISTS admin_user_password_admin_user_id ON admin_user_password(admin_user_id);

CREATE TABLE IF NOT EXISTS admin_access_permission(
    _id uuid DEFAULT uuid_generate_v4 (),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT timezone('utc', now()),
    last_modified_at TIMESTAMP WITH TIME ZONE DEFAULT timezone('utc', now()),
    admin_user_id uuid NOT NULL REFERENCES admin_user(_id),
    permission TEXT NOT NULL,
    is_active BOOLEAN NOT NULL,

    PRIMARY KEY (_id)
);

CREATE UNIQUE INDEX IF NOT EXISTS admin_access_permission_admin_user_permission ON admin_access_permission(admin_user_id, permission);

CREATE TABLE IF NOT EXISTS admin_2fa_codes(
    created_at TIMESTAMP WITH TIME ZONE DEFAULT timezone('utc', now()),
    expires_at TIMESTAMP WITH TIME ZONE,
    admin_user_id uuid NOT NULL REFERENCES admin_user(_id),
    code TEXT NOT NULL,

    PRIMARY KEY (code)
);

CREATE INDEX IF NOT EXISTS admin_2fa_codes_user_id ON admin_2fa_codes(admin_user_id);
CREATE UNIQUE INDEX IF NOT EXISTS admin_2fa_codes_user_id_unique ON admin_2fa_codes(admin_user_id) WHERE expires_at IS NULL;

CREATE TABLE IF NOT EXISTS admin_access_token(
    created_at TIMESTAMP WITH TIME ZONE DEFAULT timezone('utc', now()),
    expires_at TIMESTAMP WITH TIME ZONE,
    admin_user_id uuid NOT NULL REFERENCES admin_user(_id),
    token TEXT NOT NULL,

    PRIMARY KEY (token)
);

CREATE UNIQUE INDEX IF NOT EXISTS admin_access_token_user_id ON admin_access_token(admin_user_id);
CREATE UNIQUE INDEX IF NOT EXISTS admin_access_token_token ON admin_access_token(token);
