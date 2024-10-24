package rabbitmq_store

import (
	"fmt"

	"github.com/DimTur/lp_notification/internal/storage"
)

type Storage struct {
}

func NewRabbit() (*Storage, error) {
	return &Storage{}, nil
}

func (s *Storage) SaveUserChatID(userTg *storage.UserTg) (err error) {
	// const op = "internal.storage.rabbitmq.rabbitmq.SaveUserChatID"

	// defer func() { err = e.WrapIfErr(op, "can't save user chatId", err) }()
	fmt.Println(userTg)

	return nil
}
