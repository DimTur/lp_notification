package sender

import (
	"context"
	"encoding/json"
	"log/slog"

	tgclient "github.com/DimTur/lp_notification/internal/clients/telegram"
	rabbitmq_store "github.com/DimTur/lp_notification/internal/storage/rabbitmq"
	amqp "github.com/rabbitmq/amqp091-go"
)

type MessageQueue interface {
	Consume(
		ctx context.Context,
		queueName, consumer string,
		autoAck, exclusive, noLocal, noWait bool,
		args map[string]interface{},
		handle func(ctx context.Context, msg interface{}) error,
	) error
}

type ConsumeOTP struct {
	msgQueue MessageQueue
	tgClient *tgclient.TgClient
	logger   *slog.Logger
}

func NewConsumeOTP(
	msgQueue MessageQueue,
	tgClient *tgclient.TgClient,
	logger *slog.Logger,
) *ConsumeOTP {
	return &ConsumeOTP{
		msgQueue: msgQueue,
		tgClient: tgClient,
		logger:   logger,
	}
}

func (c *ConsumeOTP) Start(
	ctx context.Context,
	queueName, consumer string,
	autoAck, exclusive, noLocal, noWait bool,
	args map[string]interface{},
) error {
	const op = "ConsumeOTP.Start"

	log := c.logger.With(slog.String("op", op))
	log.Info("Starting to consume OTP messages")

	return c.msgQueue.Consume(
		ctx,
		queueName,
		consumer,
		autoAck,
		exclusive,
		noLocal,
		noWait,
		args,
		c.handleMessage,
	)
}

func (c *ConsumeOTP) handleMessage(ctx context.Context, msg interface{}) error {
	// Casting a message to a type amqp.Delivery
	del, ok := msg.(amqp.Delivery)
	if !ok {
		c.logger.Error("failed to cast message to amqp.Delivery")
		return nil // Return nil to avoid calling Nack/Ack
	}

	var message rabbitmq_store.MsgOTP
	// Decoding JSON message
	if err := json.Unmarshal(del.Body, &message); err != nil {
		c.logger.Error("failed to unmarshal message to MsgOTP", slog.Any("err", err))
		return err
	}

	// Sending a message in Telegram
	if err := c.tgClient.SendMessage(message.ChatID, message.Otp.Code); err != nil {
		c.logger.Error("Error sending message to Telegram", slog.Any("err", err))
		return err
	}

	c.logger.Info("Message sent to Telegram", slog.Int("chat_id", message.ChatID))

	return nil
}
