package storage

type Storage interface {
	SaveUserChatID(userTg *UserTg) error
}

type UserTg struct {
	TgLink string
	ChatID string
}
