package rabbitmq

import (
	"http_server/cmd/config"

	"github.com/streadway/amqp"
)

type RabbitMQConsumer struct{
	*RabbitMQBase
}

func NewRabbitMQConsumer(cfg config.RabbitMQ) (*RabbitMQConsumer, error){
	base, err := newRabbitMQBase(cfg)
	if err != nil {
		return nil, err
	}
	return &RabbitMQConsumer{RabbitMQBase: base}, nil
}

func (c *RabbitMQConsumer) Consume() (<-chan amqp.Delivery, error) {
    msgs, err := c.channel.Consume(
        c.queueName,
        "",
        false,
        false,
        false,
        false,
        nil,
    )
    if err != nil {
        return nil, err
    }
    return msgs, nil
}