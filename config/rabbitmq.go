package config

var Rabbitmq rabbitmq

type rabbitmq struct {
	Host     string `config:"host"`
	Port     string `config:"port"`
	Vhost    string `config:"vhost"`
	UserName string `config:"username"`
	Password string `config:"password"`

	RechargeQue string `config:"rechargeQueue"`
	ResultQeu   string `config:"resultQueue"`
}
