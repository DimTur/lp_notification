package telegram

import (
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"path"
	"strconv"

	"github.com/DimTur/lp_notification/lib/e"
)

const (
	getUpdatesMethod  = "getUpdates"
	sendMessageMethod = "sendMessage"
)

func (c *TgClient) Updates(offset int, limit int) (updates []Update, err error) {
	const op = "internal.clients.telegram.requests.Updates"

	log := c.logger.With(
		slog.String("op", op),
	)

	log.Info("getting updates")

	defer func() { err = e.WrapIfErr(log, op, "can't get updates", err) }()

	q := url.Values{}
	q.Add("offset", strconv.Itoa(offset))
	q.Add("limit", strconv.Itoa(limit))

	data, err := c.doRequest(getUpdatesMethod, q)
	if err != nil {
		return nil, err
	}

	var res UpdatesResponse

	if err := json.Unmarshal(data, &res); err != nil {
		return nil, err
	}

	log.Info("updates got successfully")

	return res.Result, nil
}

func (c *TgClient) SendMessage(chatID int, text string) error {
	const op = "internal.clients.telegram.requests.SendMessage"

	log := c.logger.With(
		slog.String("op", op),
		slog.Int("sending message to chat id:", chatID),
	)

	log.Info("sending message")

	q := url.Values{}
	q.Add("chat_id", strconv.Itoa(chatID))
	q.Add("text", text)

	_, err := c.doRequest(sendMessageMethod, q)
	if err != nil {
		return e.Wrap(log, op, "can't send message", err)
	}

	log.Info("message sended successfully")

	return nil
}

func (c *TgClient) doRequest(method string, query url.Values) (data []byte, err error) {
	const op = "internal.clients.telegram.requests.Updates"

	log := c.logger.With(
		slog.String("op", op),
	)

	defer func() { err = e.WrapIfErr(log, op, "can't do request", err) }()

	u := url.URL{
		Scheme: "https",
		Host:   c.host,
		Path:   path.Join(c.basePath, method),
	}

	req, err := http.NewRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, err
	}

	req.URL.RawQuery = query.Encode()

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}
