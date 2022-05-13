CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS billing_promotion_codes(
    created_at TIMESTAMP WITH TIME ZONE DEFAULT timezone('utc', now()),
    last_modified_at TIMESTAMP WITH TIME ZONE DEFAULT timezone('utc', now()),
    _id uuid DEFAULT uuid_generate_v4 (),
    type TEXT NOT NULL,
    external_id_mapping_id uuid NOT NULL REFERENCES billing_external_id_mapping(_id),
    code TEXT NOT NULL,

    PRIMARY KEY(_id)
);

CREATE UNIQUE INDEX IF NOT EXISTS billing_promotion_codes_external_id ON billing_promotion_codes(external_id_mapping_id);
CREATE UNIQUE INDEX IF NOT EXISTS billing_promotion_codes_code ON billing_promotion_codes(code);

CREATE TABLE IF NOT EXISTS billing_user_promotion(
    created_at TIMESTAMP WITH TIME ZONE DEFAULT timezone('utc', now()),
    last_modified_at TIMESTAMP WITH TIME ZONE DEFAULT timezone('utc', now()),
    promotion_id uuid NOT NULL REFERENCES billing_promotion_codes(_id),
    billing_information_id uuid NOT NULL REFERENCES billing_information(_id),

    PRIMARY KEY (billing_information_id, promotion_id)
);
