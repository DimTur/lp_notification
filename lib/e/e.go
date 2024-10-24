package e

import (
	"fmt"
	"log/slog"
)

func Wrap(log *slog.Logger, op string, msg string, err error) error {
	log.Error(msg, slog.String("err", err.Error()))
	return fmt.Errorf("%s: %s: %w", op, msg, err)
}

func WrapIfErr(log *slog.Logger, op string, msg string, err error) error {
	if err == nil {
		return nil
	}

	return Wrap(log, op, msg, err)
}
