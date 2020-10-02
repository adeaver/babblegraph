CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE users(
    _id uuid DEFAULT uuid_generate_v4 (),
    email_address TEXT NOT NULL,
    status TEXT NOT NULL,

    PRIMARY KEY (_id)
);

CREATE UNIQUE INDEX IF NOT EXISTS users_email_address_idx ON users(email_address);
