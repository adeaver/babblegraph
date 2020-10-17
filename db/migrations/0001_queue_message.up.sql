CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS queue_messages(
    _id uuid DEFAULT uuid_generate_v4 (),
    topic TEXT NOT NULL,
    is_enqueued BOOLEAN DEFAULT true,
    queue_position SERIAL NOT NULL,
    body TEXT NOT NULL,
    PRIMARY KEY (_id)
);

CREATE UNIQUE INDEX IF NOT EXISTS queue_messages_pos_idx ON queue_messages(topic, queue_position);
