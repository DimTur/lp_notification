package telegram

import (
	"log/slog"
	"net/http"
)

type TgClient struct {
	host     string
	basePath string
	client   *http.Client

	logger *slog.Logger
}

func NewTgClient(
	host string,
	token string,

	logger *slog.Logger,
) (*TgClient, error) {
	const op = "tg client"

	logger = logger.With(
		slog.String("op", op),
		slog.String("host", host),
	)

	tgClient := &TgClient{
		host:     host,
		basePath: newBasePath(token),
		client:   &http.Client{},

		logger: logger,
	}

	return tgClient, nil
}

func newBasePath(token string) string {
	return "bot" + token
}
