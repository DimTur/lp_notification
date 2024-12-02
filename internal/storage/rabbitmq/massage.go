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

type NotificationMsg struct {
	UserID    string `json:"user_id"`
	Email     string `json:"email"`
	TgLink    string `json:"tg_link"`
	ChatID    string `json:"chat_id"`
	ChannelID int64  `json:"channel_id"`
	PlanID    int64  `json:"plan_id"`
	CreatedBy string `json:"created_by"`
}
