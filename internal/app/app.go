package app

// import (
// 	"log/slog"

// 	"github.com/DimTur/lp_notification/internal/app/telegram"
// 	tgevents "github.com/DimTur/lp_notification/internal/events/telegram"
// 	"github.com/DimTur/lp_notification/internal/storage"
// )

// type RabbitMq interface {
// 	storage.Storage
// }

// type TgBotApp struct {
// 	EventClient telegram.EventClient
// }

// func NewTgBotApp(
// 	host string,
// 	token string,
// 	rabbitMq RabbitMq,
// 	logger *slog.Logger,
// ) (*TgBotApp, error) {

// 	tgEventClient, err := telegram.NewEventClient(host, token, logger)
// 	if err != nil {
// 		return nil, err
// 	}

// 	eventProcessor := tgevents.New(tgEventClient.tgClient, rabbitMq, logger)
// }
