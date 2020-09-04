package queue

import (
	"github.com/jmoiron/sqlx"
)

type Queue interface {
	GetTopicName() string
	ProcessMessage(*sqlx.Tx, Message) error
}
