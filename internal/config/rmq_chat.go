package config

type Chat struct {
	ChatIDExchange   ChatIDExchange `yaml:"chat_id_exchange"`
	ChatIDQueue      QueueConfig    `yaml:"chat_id_queue"`
	ChatIDRoutingKey string         `yaml:"chat_id_routing_key"`
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
