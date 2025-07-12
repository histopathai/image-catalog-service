package service

import (
	"context"
	"encoding/json"

	"github.com/histopathai/image-catalog-service/internal/models"
	"github.com/streadway/amqp"
)

type MessageQueue interface {
	PublishProcessingJob(ctx context.Context, job *models.ProcessingJob) error
	ConsumeProcessingResults(ctx context.Context) (<-chan amqp.Delivery, error)
}
type RabbitMQ struct {
	channel     *amqp.Channel
	jobQueue    string
	resultQueue string
}

func NewRabbitMQ(conn *amqp.Connection, jobQueue, resultQueue string) (*RabbitMQ, error) {
	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}
	_, err = ch.QueueDeclare(jobQueue, true, false, false, false, nil)
	if err != nil {
		return nil, err
	}

	_, err = ch.QueueDeclare(resultQueue, true, false, false, false, nil)
	if err != nil {
		return nil, err
	}

	return &RabbitMQ{
		channel:     ch,
		jobQueue:    jobQueue,
		resultQueue: resultQueue,
	}, nil
}

func (r *RabbitMQ) PublishProcessingJob(ctx context.Context, job *models.ProcessingJob) error {
	body, err := json.Marshal(job)
	if err != nil {
		return err
	}

	return r.channel.Publish("", r.jobQueue, false, false, amqp.Publishing{
		ContentType: "application/json",
		Body:        body,
	})
}

func (r *RabbitMQ) ConsumeProcessingResults(ctx context.Context) (<-chan amqp.Delivery, error) {
	msgs, err := r.channel.Consume(
		r.resultQueue, // listening to this queue
		"",            // consumer tag
		true,          // auto-ack
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
