package main

import (
	"log"
	"net/http"

	conf "rongshu.tech/oiling/proxy_gateway/config"
	"rongshu.tech/oiling/proxy_gateway/controller"
	"rongshu.tech/oiling/proxy_gateway/util/charge/spcard"
	"rongshu.tech/oiling/proxy_gateway/util/config"
	"rongshu.tech/oiling/proxy_gateway/util/mq/rabbitmq"
)

func init() {
	log.SetFlags(log.Lshortfile | log.Ldate)
}

//欠一块，配置文件的读取
func main() {

	cfg := config.Conf{FileName: "gateway.ini"}
	cfg.Instance()

	cfg.AddConfig("rabbitmq", &conf.Rabbitmq)
	cfg.AddConfig("mysql", &conf.Mysql)
	cfg.AddConfig("recharge", &conf.Spcard)

	rabbitmq.InitMQ(conf.Rabbitmq.Host, conf.Rabbitmq.Port, conf.Rabbitmq.Vhost,
		conf.Rabbitmq.UserName, conf.Rabbitmq.Password)
	controller.Init()

	log.Println("start listen MQ ")
	log.Println("start server listen at port", conf.Spcard.ListenPort)

	spcard.Init()
	go controller.MsgProcessing()
	go controller.ResultStatistics()
	go controller.DealLetterQueue()

	http.HandleFunc("/health", controller.HealthCheck)
	log.Fatal(http.ListenAndServe(conf.Spcard.ListenPort, nil))

}
