CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS user_surveys(
    _id TEXT NOT NULL,
    sent_at TIMESTAMP WITH TIME ZONE DEFAULT timezone('utc', now()),
    user_id uuid NOT NULL REFERENCES users(_id),
    first_opened_at TIMESTAMP WITH TIME ZONE,
    survey_type TEXT NOT NULL,

    PRIMARY KEY (_id)
);

CREATE INDEX IF NOT EXISTS user_surveys_user_id ON user_surveys(user_id);

CREATE TABLE IF NOT EXISTS user_survey_responses(
    _id uuid DEFAULT uuid_generate_v4 (),
    question_id TEXT NOT NULL,
    survey_id TEXT NOT NULL REFERENCES user_surveys(_id),
    submitted_at TIMESTAMP WITH TIME ZONE DEFAULT timezone('utc', now()),
    answer TEXT NOT NULL,
);
