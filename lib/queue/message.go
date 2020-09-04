package queue

import "github.com/jmoiron/sqlx"

type MessageID string

type Message struct {
	ID          MessageID
	MessageBody string
}

type dbMessage struct {
	ID            MessageID `db:"_id"`
	Topic         string    `db:"topic"`
	QueuePosition int       `db:"queue_position"`
	IsEnqueued    bool      `db:"is_enqueued"`
	Body          string    `db:"body"`
}

func (m *dbMessage) ToMessage() Message {
	return Message{
		ID:          m.ID,
		MessageBody: m.Body,
	}
}

func saveMessage(tx *sqlx.Tx, msg dbMessage) error {
	// TODO: figure out a less brittle way to use this
	_, err := tx.Exec(
		"INSERT INTO queue_messages (topic, body) VALUES ($1, $2)",
		msg.Topic, msg.Body,
	)
	return err
}

func getMostRecentMessageForQueue(tx *sqlx.Tx, topicName string) (*dbMessage, error) {
	var msg dbMessage
	if err := tx.Get(&msg, "SELECT * FROM queue_messages WHERE topic=$1 AND is_enqueued=TRUE ORDER BY queue_position ASC LIMIT 1", topicName); err != nil {
		return nil, err
	}
	return &msg, nil
}

func clearMessage(tx *sqlx.Tx, messageID MessageID) error {
	_, err := tx.Exec("UPDATE queue_messages SET is_enqueued=FALSE WHERE _id=$1", messageID)
	return err
}
