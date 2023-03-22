package config

type RabbitMq struct {
	Amqp       string `mapstructure:"amqp" json:"amqp" yaml:"amqp"`
	MaxTaskNum int    `mapstructure:"maxTaskNum" json:"maxTaskNum" form:"maxTaskNum"`
}

func (r *RabbitMq) GetRabbitMqAddr() string {
	return r.Amqp
}
