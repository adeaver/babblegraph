CREATE TABLE IF NOT EXISTS corpora(
    _id TEXT NOT NULL,
    language TEXT NOT NULL REFERENCES languages(code),
    name TEXT NOT NULL,

    PRIMARY KEY (_id)
);
