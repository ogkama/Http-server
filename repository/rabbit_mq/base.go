package rabbitmq

import (
	"fmt"
	"http_server/cmd/config"

	"github.com/streadway/amqp"
)


type RabbitMQBase struct {
    connection *amqp.Connection
    channel    *amqp.Channel
    queueName  string
}


func newRabbitMQBase(cfg config.RabbitMQ) (*RabbitMQBase, error) {
    url := fmt.Sprintf("amqp://guest:guest@%s:%d", cfg.Host, cfg.Port)
    conn, err := amqp.Dial(url)
    if err != nil {
        return nil, fmt.Errorf("connecting to rabbitMQ: %w", err)
    }

    ch, err := conn.Channel()
    if err != nil {
        return nil, err
    }

    _, err = ch.QueueDeclare(
        cfg.QueueName, // name
        true,      // durable
        false,     // delete when unused
        false,     // exclusive
        false,     // no-wait
        nil,       // arguments
    )
    if err != nil {
        return nil, err
    }

    return &RabbitMQBase{
        connection: conn,
        channel:    ch,
        queueName:  cfg.QueueName,
    }, nil
}

func (r *RabbitMQBase) Close() {
    r.channel.Close()
    r.connection.Close()
}