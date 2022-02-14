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
