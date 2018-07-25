package main

import (
	"github.com/gin-gonic/gin"
	"log"

	conf "github.com/vgmdj/tmp/config"
	"github.com/vgmdj/tmp/controller"
	"github.com/vgmdj/tmp/util/config"
	"github.com/vgmdj/tmp/util/mq/rabbitmq"
)

func init() {
	log.SetFlags(log.Lshortfile | log.Ldate)
}

//欠一块，配置文件的读取
func main() {

	cfg := config.Conf{FileName: "gateway.ini"}
	cfg.Instance()

	cfg.AddConfig("rabbitmq", &conf.Rabbitmq)

	rabbitmq.InitMQ(conf.Rabbitmq.Host, conf.Rabbitmq.Port, conf.Rabbitmq.Vhost,
		conf.Rabbitmq.UserName, conf.Rabbitmq.Password)
	controller.Init()

	log.Println("start listen MQ ")
	log.Println("start server listen at port", conf.Spcard.ListenPort)

	go controller.MsgProcessing()

	r := gin.Default()
	r.POST("/charge/result", controller.ResultStatistics)
	r.Run(":9099")

}
