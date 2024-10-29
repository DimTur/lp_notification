package telegram

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"strconv"
	"strings"

	"github.com/DimTur/lp_notification/internal/storage"
)

const (
	StartCmd = "/start"

	exchangeChatID   = "chat_id"
	queueChatID      = "chat_id"
	chatIDRoutingKey = "chat_id"
)

func (p *Processor) doCmd(text string, chatID int, username string) error {
	const op = "internal.events.telegram.commands.doCmd"

	text = strings.TrimSpace(text)

	log := p.logger.With(
		slog.String("op", op),
		slog.Int("chat_id", chatID),
		slog.String("username", username),
	)

	log.Info("sending new command")

	switch text {
	case StartCmd:
		return p.saveChatID(username, chatID)
	default:
		return p.tg.SendMessage(chatID, msgUnknownCommand)
	}
}

func (p *Processor) saveChatID(username string, chatID int) error {
	const op = "runTg"

	log := p.logger.With(
		slog.String("op", op),
		slog.String("username", username),
		slog.Int("chatID", chatID),
	)

	chatIDStr := strconv.Itoa(chatID)

	// TODO: delete
	fmt.Printf("username: %s, chat ID: %s", username, chatIDStr)

	msgUserTg := &storage.UserTg{
		TgLink: username,
		ChatID: chatIDStr,
	}
	msgBody, err := json.Marshal(msgUserTg)
	if err != nil {
		log.Error("err to marshal otp", slog.String("err", err.Error()))
		return fmt.Errorf("%s: %w", op, err)
	}

	if err := p.storage.Publish(p.ctx, exchangeChatID, chatIDRoutingKey, msgBody); err != nil {
		p.tg.SendMessage(chatID, msgUserNotFound)
		return err
	}

	p.tg.SendMessage(chatID, msgHello)
	p.tg.SendMessage(chatID, msgSaved)
	return nil
}
