package eventconsumer

import (
	"context"
	"log/slog"
	"time"

	"github.com/DimTur/lp_notification/internal/events"
)

type Consumer struct {
	ctx       context.Context
	fetcher   events.Fetcher
	processor events.Processor
	batchSize int

	logger *slog.Logger
}

func New(
	ctx context.Context,
	fetcher events.Fetcher,
	processor events.Processor,
	batchSize int,
	logger *slog.Logger,
) Consumer {
	return Consumer{
		ctx:       ctx,
		fetcher:   fetcher,
		processor: processor,
		batchSize: batchSize,
		logger:    logger,
	}
}

func (c *Consumer) Start() error {
	const op = "internal.consumer.event-consumer.event-consumer.Start"

	log := c.logger.With(
		slog.String("op", op),
	)

	log.Info("starting event-consumer")

	for {
		select {
		case <-c.ctx.Done():
			log.Info("consumer is shutting down")
			return nil
		default:
			gotEvents, err := c.fetcher.Fetch(c.batchSize)
			if err != nil {
				log.Error("consumer:", slog.String("err", err.Error()))
				time.Sleep(1 * time.Second)
				continue
			}

			if len(gotEvents) == 0 {
				time.Sleep(1 * time.Second)
				continue
			}

			if err := c.handleEvents(gotEvents); err != nil {
				log.Error("can't handle event:", slog.String("err", err.Error()))
			}
		}
	}
}

/*
Possible problems and solutions:
1. Loss of events: retry, return to storage, fallback, confirmation for fetcher
2. Processing the entire batch: stop after the first error, error counter
3. Parallel processing
*/
func (c *Consumer) handleEvents(events []events.Event) error {
	const op = "internal.consumer.event-consumer.event-consumer.handleEvents"

	log := c.logger.With(
		slog.String("op", op),
	)

	for _, event := range events {
		log.Info("got new", slog.String("event", event.Text))

		if err := c.retryProcess(event); err != nil {
			log.Error("failed to handle event after retries", slog.String("err", err.Error()))
			c.saveFailedEvent(event) // Saving a failed event
		}
	}

	return nil
}

func (c *Consumer) retryProcess(event events.Event) error {
	const maxRetries = 3
	var err error
	for attempt := 1; attempt <= maxRetries; attempt++ {
		err = c.processor.Process(event)
		if err == nil {
			return nil // Successful processing
		}
		time.Sleep(time.Duration(attempt) * time.Second) // Exponential backoff
	}
	return err
}

func (c *Consumer) saveFailedEvent(event events.Event) {
	// TODO: Logic for saving an event to a database or queue
	log := c.logger.With(slog.String("op", "internal.consumer.event-consumer.event-consumer.saveFailedEvent"))
	log.Info("saving failed event", slog.String("event", event.Text))
}
