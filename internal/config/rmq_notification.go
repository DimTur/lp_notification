package config

type Notification struct {
	NotificationQueue    NotificationQueue    `yaml:"notification_queue"`
	NotificationExchange NotificationExchange `yaml:"notification_exchange"`
	NotificationConsumer NotificationConsumer `yaml:"notification_consumer"`
}

type NotificationQueue struct {
	Name string `yaml:"name"`
}

type NotificationExchange struct {
	Name        string       `yaml:"name"`
	Kind        string       `yaml:"kind"`
	Durable     bool         `yaml:"durable"`
	AutoDeleted bool         `yaml:"auto_deleted"`
	Internal    bool         `yaml:"internal"`
	NoWait      bool         `yaml:"no_wait"`
	Args        ExchangeArgs `yaml:"args"`
}

type NotificationConsumer struct {
	Queue        string       `yaml:"queue"`
	Consumer     string       `yaml:"consumer"`
	AutoAck      bool         `yaml:"autoAck"`
	Exclusive    bool         `yaml:"exclusive"`
	NoLocal      bool         `yaml:"noLocal"`
	NoWait       bool         `yaml:"noWait"`
	ConsumerArgs ConsumerArgs `yaml:"args"`
}
