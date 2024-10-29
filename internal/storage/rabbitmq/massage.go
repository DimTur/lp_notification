package rabbitmq_store

import (
	"time"
)

type OTP struct {
	UserID    string    `json:"user_id"`
	Code      string    `json:"code"`
	ExpiresAt time.Time `json:"expires_at"`
	Used      bool      `json:"used"`
}

type MsgOTP struct {
	Otp    OTP `json:"otp"`
	ChatID int `json:"chat_id"`
}

type UserTg struct {
	TgLink string `json:"tg_link"`
	ChatID string `json:"chat_id"`
}
