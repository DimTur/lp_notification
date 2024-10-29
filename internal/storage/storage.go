package storage

import "context"

type Storage interface {
	Publish(ctx context.Context, exchange, routingKey string, body []byte) error
	Consume(ctx context.Context, queueName string, handle func(msg interface{}) error) error
}

type UserTg struct {
	TgLink string
	ChatID string
}
