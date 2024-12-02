package config

type OTP struct {
	OTPExchange   OTPExchange `yaml:"otp_exchange"`
	OTPQueue      QueueConfig `yaml:"otp_queue"`
	OTPConsumer   OTPConsumer `yaml:"otp_consumer"`
	OTPRoutingKey string      `yaml:"otp_routing_key"`
}

type OTPQueue struct {
	Name string `yaml:"name"`
}

type OTPExchange struct {
	Name        string       `yaml:"name"`
	Kind        string       `yaml:"kind"`
	Durable     bool         `yaml:"durable"`
	AutoDeleted bool         `yaml:"auto_deleted"`
	Internal    bool         `yaml:"internal"`
	NoWait      bool         `yaml:"no_wait"`
	Args        ExchangeArgs `yaml:"args"`
}

type OTPConsumer struct {
	Queue        string       `yaml:"queue"`
	Consumer     string       `yaml:"consumer"`
	AutoAck      bool         `yaml:"autoAck"`
	Exclusive    bool         `yaml:"exclusive"`
	NoLocal      bool         `yaml:"noLocal"`
	NoWait       bool         `yaml:"noWait"`
	ConsumerArgs ConsumerArgs `yaml:"args"`
}
