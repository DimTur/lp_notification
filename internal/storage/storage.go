package storage

import "context"

type Storage interface {
	Publish(ctx context.Context, exchange, routingKey string, body []byte) error
	Consume(
		ctx context.Context,
		queueName, consumer string,
		autoAck, exclusive, noLocal, noWait bool,
		args map[string]interface{},
		handle func(ctx context.Context, msg interface{}) error,
	) error
}
