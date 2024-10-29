package telegram

import (
	"context"
	"log/slog"

	tgclient "github.com/DimTur/lp_notification/internal/clients/telegram"
	eventconsumer "github.com/DimTur/lp_notification/internal/consumer/event-consumer"
	tgevents "github.com/DimTur/lp_notification/internal/events/telegram"
	"github.com/DimTur/lp_notification/internal/storage"
)

type Storage interface {
	storage.Storage
}

func RunTg(
	ctx context.Context,
	tgClient *tgclient.TgClient,
	batchSize int,
	storage Storage,
	logger *slog.Logger,
) error {
	const op = "runTg"

	log := logger.With(
		slog.String("op", op),
	)

	eventProcessor := tgevents.New(ctx, tgClient, storage, logger)

	consumer := eventconsumer.New(ctx, eventProcessor, eventProcessor, batchSize, logger)
	if err := consumer.Start(); err != nil {
		log.Info("tg bot is stopped", slog.String("err", err.Error()))
	}

	return nil
}
