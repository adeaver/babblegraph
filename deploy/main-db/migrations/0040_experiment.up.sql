CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS experiments(
    _id uuid DEFAULT uuid_generate_v4 (),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT timezone('utc', now()),
    last_modified_at TIMESTAMP WITH TIME ZONE DEFAULT timezone('utc', now()),
    is_active BOOLEAN DEFAULT false,
    name TEXT NOT NULL,
    current_step INTEGER NOT NULL CHECK (current_step <= 100 AND current_step >= 0),
    previous_step INTEGER NOT NULL CHECK (previous_step <= 100 AND previous_step >= 0),

    PRIMARY KEY(_id)
);

CREATE UNIQUE INDEX IF NOT EXISTS experiments_unique_name_idx ON experiments(name);

CREATE TABLE IF NOT EXISTS experiments_user_variations(
    _id uuid DEFAULT uuid_generate_v4 (),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT timezone('utc', now()),
    last_modified_at TIMESTAMP WITH TIME ZONE DEFAULT timezone('utc', now()),
    user_id uuid NOT NULL REFERENCES users(_id),
    accessed_at_step INTEGER NOT NULL,
    experiment_id uuid NOT NULL REFERENCES experiments(_id),
    in_variation BOOLEAN NOT NULL,

    PRIMARY KEY(_id)
);

CREATE UNIQUE INDEX IF NOT EXISTS experiments_user_variations_user_for_experiment_idx ON experiments_user_variations(experiment_id, user_id);

ALTER TABLE experiments_user_variations ADD COLUMN IF NOT EXISTS in_experiment BOOLEAN NOT NULL DEFAULT TRUE;
