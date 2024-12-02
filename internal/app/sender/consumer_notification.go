package sender

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"strconv"

	tgclient "github.com/DimTur/lp_notification/internal/clients/telegram"
	rabbitmq_store "github.com/DimTur/lp_notification/internal/storage/rabbitmq"
	amqp "github.com/rabbitmq/amqp091-go"
)

type ConsumeNotification struct {
	msgQueue MessageQueue
	tgClient *tgclient.TgClient
	logger   *slog.Logger
}

func NewConsumeNotification(
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

func (c *ConsumeNotification) Start(
	ctx context.Context,
	queueName, consumer string,
	autoAck, exclusive, noLocal, noWait bool,
	args map[string]interface{},
) error {
	const op = "ConsumeNotification.Start"

	log := c.logger.With(slog.String("op", op))
	log.Info("Starting to consume notification messages")

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

func (c *ConsumeNotification) handleMessage(ctx context.Context, msg interface{}) error {
	const op = "ConsumeNotification.handleMessage"

	// Casting a message to a type amqp.Delivery
	del, ok := msg.(amqp.Delivery)
	if !ok {
		c.logger.Error("failed to cast message to amqp.Delivery")
		return nil // Return nil to avoid calling Nack/Ack
	}

	var message rabbitmq_store.NotificationMsg
	// Decoding JSON message
	if err := json.Unmarshal(del.Body, &message); err != nil {
		c.logger.Error("failed to unmarshal message", slog.Any("err", err))
		return err
	}

	// TODO: sent link to email. Do some html template
	m := fmt.Sprintf("Нou have a new plan and lessons to go through. Link: http://localhost:8000/channels/%d/plans/%d", message.ChannelID, message.PlanID)
	fmt.Println(m)

	// Sending a message in Telegram
	if message.ChatID != "" {
		chatIDInt, err := strconv.Atoi(message.ChatID)
		if err != nil {
			c.logger.Error("err to convert chat_id to int", slog.String("err", err.Error()))
			return fmt.Errorf("%s: %w", op, err)
		}

		if err := c.tgClient.SendMessage(chatIDInt, m); err != nil {
			c.logger.Error("Error sending message to Telegram", slog.Any("err", err))
			return err
		}

		c.logger.Info("Message sent to Telegram", slog.Int("chat_id", chatIDInt))
	}

	c.logger.Info("Message sent to Email", slog.String("email", message.Email))

	return nil
}
