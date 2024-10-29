package sender

import (
	"context"
	"encoding/json"
	"log/slog"

	tgclient "github.com/DimTur/lp_notification/internal/clients/telegram"
	"github.com/DimTur/lp_notification/internal/storage"
	rabbitmq_store "github.com/DimTur/lp_notification/internal/storage/rabbitmq"
	amqp "github.com/rabbitmq/amqp091-go"
)

type ConsumeOTP struct {
	storage  storage.Storage
	tgClient *tgclient.TgClient
	logger   *slog.Logger
}

func NewConsumeOTP(
	storage storage.Storage,
	tgClient *tgclient.TgClient,
	logger *slog.Logger,
) *ConsumeOTP {
	return &ConsumeOTP{
		storage:  storage,
		tgClient: tgClient,
		logger:   logger,
	}
}

func (c *ConsumeOTP) Start(ctx context.Context, queueName string) error {
	const op = "ConsumeOTP.Start"

	log := c.logger.With(slog.String("op", op))
	log.Info("Starting to consume OTP messages")

	return c.storage.Consume(ctx, queueName, c.handleMessage)
}

func (c *ConsumeOTP) handleMessage(msg interface{}) error {
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
