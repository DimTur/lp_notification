package telegram

import (
	"fmt"
	"log/slog"
	"strconv"
	"strings"

	"github.com/DimTur/lp_notification/internal/storage"
)

const (
	StartCmd = "/start"
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
	chatIDStr := strconv.Itoa(chatID)

	fmt.Printf("username: %s, chat ID: %s", username, chatIDStr)
	userTg := &storage.UserTg{
		TgLink: username,
		ChatID: chatIDStr,
	}
	if err := p.storage.SaveUserChatID(userTg); err != nil {
		p.tg.SendMessage(chatID, msgUserNotFound)
		return err
	}

	p.tg.SendMessage(chatID, msgHello)
	p.tg.SendMessage(chatID, msgSaved)
	return nil
}
