package config

import "github.com/ilyakaznacheev/cleanenv"

type Config struct {
	Server      Server      `yaml:"server"`
	TelegramBot TelegramBot `yaml:"telegram_bot"`
	RabbitMQ    RabbitMQ    `yaml:"rabbit_mq"`
}

type Server struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
}

type TelegramBot struct {
	TgBotToken string `yaml:"tg_bot_token"`
	TgBotHost  string `yaml:"tg_bot_host"`
	BatchSize  int    `yaml:"batch_size"`
}

type RabbitMQ struct {
	UserName         string         `yaml:"username"`
	Password         string         `yaml:"password"`
	Host             string         `yaml:"host"`
	Port             int            `yaml:"port"`
	OTPQueue         OTPQueue       `yaml:"otp_queue"`
	ChatIDExchange   ChatIDExchange `yaml:"chat_id_exchange"`
	ChatIDQueue      ChatIDQueue    `yaml:"chat_id_queue"`
	ChatIDRoutingKey string         `yaml:"chat_id_routing_key"`
}

type OTPQueue struct {
	Name string `yaml:"name"`
}

type ChatIDExchange struct {
	Name        string       `yaml:"name"`
	Kind        string       `yaml:"kind"`
	Durable     bool         `yaml:"durable"`
	AutoDeleted bool         `yaml:"auto_deleted"`
	Internal    bool         `yaml:"internal"`
	NoWait      bool         `yaml:"no_wait"`
	Args        ExchangeArgs `yaml:"args"`
}

type ExchangeArgs struct {
	AltExchange string `yaml:"alternate_exchange"`
}

type ChatIDQueue struct {
	Name        string          `yaml:"name"`
	Durable     bool            `yaml:"durable"`
	AutoDeleted bool            `yaml:"auto_deleted"`
	Exclusive   bool            `yaml:"exclusive"`
	NoWait      bool            `yaml:"no_wait"`
	Args        ChatIDQueueArgs `yaml:"args"`
}

type ChatIDQueueArgs struct {
	XMessageTtl int32 `yaml:"x_message_ttl"`
}

func Parse(s string) (*Config, error) {
	c := &Config{}
	if err := cleanenv.ReadConfig(s, c); err != nil {
		return nil, err
	}

	return c, nil
}

func (e ExchangeArgs) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"alternate-exchange": e.AltExchange,
	}
}

func (q ChatIDQueueArgs) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"x-message-ttl": q.XMessageTtl,
	}
}
