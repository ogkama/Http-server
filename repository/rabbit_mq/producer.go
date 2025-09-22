package rabbitmq

import (
	"encoding/json"
	"http_server/cmd/config"
	"http_server/domain"
	"log"

	"github.com/streadway/amqp"
)

type RabbitMQSender struct{
	*RabbitMQBase
}

func NewRabbitMQSender(cfg config.RabbitMQ) (*RabbitMQSender, error){
	base, err := newRabbitMQBase(cfg)
	if err != nil {
		return nil, err
	}
	return &RabbitMQSender{RabbitMQBase: base}, nil
}

func (r *RabbitMQSender) Send(object domain.Task) error {
    log.Println("Ino Data Len consumer: ", len(object.Data))
    body, err := json.Marshal(object)
    if err != nil {
        return err
    }

    err = r.channel.Publish(
        "",              // exchange
        r.queueName,     // routing key
        false,           // mandatory
        false,           // immediate
        amqp.Publishing{
            ContentType: "application/json",
            Body:        body,
        })
    if err != nil {
        return err
    }

    return nil
}