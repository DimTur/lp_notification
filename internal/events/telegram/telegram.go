package telegram

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/DimTur/lp_notification/internal/clients/telegram"
	"github.com/DimTur/lp_notification/internal/events"
	"github.com/DimTur/lp_notification/internal/storage"
	"github.com/DimTur/lp_notification/lib/e"
)

type Processor struct {
	ctx     context.Context
	tg      *telegram.TgClient
	offset  int
	storage storage.Storage

	logger *slog.Logger
}

type Meta struct {
	ChatID   int
	Username string
}

var (
	ErrUnknownEventType = errors.New("unknown event type")
	ErrUnknownMetaType  = errors.New("unknown meta type")
)

func New(
	ctx context.Context,
	client *telegram.TgClient,
	storage storage.Storage,
	logger *slog.Logger,
) *Processor {
	return &Processor{
		ctx:     ctx,
		tg:      client,
		storage: storage,
		logger:  logger,
	}
}

func (p *Processor) Fetch(limit int) ([]events.Event, error) {
	const op = "internal.events.telegram.telegram.Fetch"

	log := p.logger.With(
		slog.String("op", op),
	)

	// log.Info("fetching events")

	updates, err := p.tg.Updates(p.offset, limit)
	if err != nil {
		return nil, e.Wrap(log, op, "can't get events", err)
	}

	if len(updates) == 0 {
		return nil, nil
	}

	res := make([]events.Event, 0, len(updates))

	for _, u := range updates {
		res = append(res, event(u))
	}

	p.offset = updates[len(updates)-1].ID + 1

	log.Info("events fetched successfully")

	return res, nil
}

func (p Processor) Process(event events.Event) error {
	const op = "internal.events.telegram.telegram.Process"

	log := p.logger.With(
		slog.String("op", op),
	)

	log.Info("processing event")

	switch event.Type {
	case events.Message:
		return p.processMessage(event)
	default:
		return e.Wrap(log, op, "can't process message", ErrUnknownEventType)
	}
}

func (p *Processor) processMessage(event events.Event) error {
	const op = "internal.events.telegram.telegram.processMessage"

	log := p.logger.With(
		slog.String("op", op),
	)

	meta, err := meta(event)
	if err != nil {
		return e.Wrap(log, op, "can't process message", err)
	}

	if err := p.doCmd(event.Text, meta.ChatID, meta.Username); err != nil {
		return e.Wrap(log, op, "can't process message", err)
	}

	return nil
}

func meta(event events.Event) (Meta, error) {
	const op = "internal.events.telegram.telegram.meta"

	res, ok := event.Meta.(Meta)
	if !ok {
		return Meta{}, fmt.Errorf("%s: %w", op, ErrUnknownMetaType)
	}

	return res, nil
}

func event(upd telegram.Update) events.Event {
	updType := fetchType(upd)

	res := events.Event{
		Type: updType,
		Text: fetchText(upd),
	}

	if updType == events.Message {
		res.Meta = Meta{
			ChatID:   upd.Message.Chat.ID,
			Username: upd.Message.From.Username,
		}
	}

	return res
}

func fetchType(upd telegram.Update) events.Type {
	if upd.Message == nil {
		return events.Unknown
	}

	return events.Message
}

func fetchText(upd telegram.Update) string {
	if upd.Message == nil {
		return ""
	}

	return upd.Message.Text
}
