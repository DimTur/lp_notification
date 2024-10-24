package telegram

import (
	"context"
	"log/slog"

	tgclient "github.com/DimTur/lp_notification/internal/clients/telegram"
	eventconsumer "github.com/DimTur/lp_notification/internal/consumer/event-consumer"
	tgevents "github.com/DimTur/lp_notification/internal/events/telegram"
	"github.com/DimTur/lp_notification/internal/storage"
)

type Client struct {
}

type Storage interface {
	storage.Storage
}

func RunTg(
	ctx context.Context,
	host string,
	token string,
	batchSize int,
	storage Storage,
	logger *slog.Logger,
) error {
	const op = "runTg"

	log := logger.With(
		slog.String("op", op),
		slog.String("host", host),
	)

	tgClient, err := tgclient.NewTgClient(host, token, logger)
	if err != nil {
		return err
	}

	eventProcessor := tgevents.New(tgClient, storage, logger)

	consumer := eventconsumer.New(ctx, eventProcessor, eventProcessor, batchSize, logger)
	if err := consumer.Start(); err != nil {
		log.Info("tg bot is stopped", slog.String("err", err.Error()))
	}

	return nil
}
