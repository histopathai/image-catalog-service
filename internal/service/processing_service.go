package service

import (
	"context"
	"encoding/json"

	"github.com/histopathai/image-catalog-service/internal/models"
	"github.com/streadway/amqp"
)

type MessageQueue interface {
	PublishProcessingJob(ctx context.Context, job *models.ProcessingJob) error
}
type RabbitMQPublisher struct {
	channel *amqp.Channel
	queue   string
}

func NewRabbitMQPublisher(conn *amqp.Connection, queue string) (*RabbitMQPublisher, error) {
	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}
	_, err = ch.QueueDeclare(queue, true, false, false, false, nil)
	if err != nil {
		ch.Close()
		return nil, err
	}
	return &RabbitMQPublisher{
		channel: ch,
		queue:   queue,
	}, nil
}

func (p *RabbitMQPublisher) PublishProcessingJob(ctx context.Context, job *models.ProcessingJob) error {
	body, err := json.Marshal(job)
	if err != nil {
		return err
	}

	return p.channel.Publish("", p.queue, false, false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		})
}
