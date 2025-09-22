package repository

import (
	"http_server/domain"

	"github.com/streadway/amqp"
)

type TaskSender interface {
	Send(object domain.Task) error
}

type TaskConsumer interface {
	Consume() (<-chan amqp.Delivery, error)
}