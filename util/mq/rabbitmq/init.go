package rabbitmq

import (
	"fmt"
	"github.com/streadway/amqp"
	"log"
	"time"
)

var (
	rabbitmqHost     = ""
	rabbitmqPort     = ""
	rabbitmqVhost    = ""
	rabbitmqUserName = ""
	rabbitmqPassword = ""

	conn *amqp.Connection
	ch   *amqp.Channel

	exist bool = false
)

const (
	exchange = "amqp_test_exchange"
	key      = "amqp_test_key"
	queue    = "amqp_test_queue"
)

//InitMQ
func InitMQ(host, port, vhost, userName, password string) {
	rabbitmqHost = host
	rabbitmqPort = port
	rabbitmqVhost = vhost
	rabbitmqUserName = userName
	rabbitmqPassword = password

	if vhost == "" {
		rabbitmqVhost = "/"
	} else if vhost[0] != '/' {
		rabbitmqVhost = "/" + vhost
	}

	connect()

	go checkQueue()
}

func connect() {
	var err error

	dialUrl := fmt.Sprintf("amqp://%s:%s@%s:%s%s", rabbitmqUserName,
		rabbitmqPassword, rabbitmqHost, rabbitmqPort, rabbitmqVhost)
	conn, err = amqp.Dial(dialUrl)
	if err != nil {
		log.Printf("%s: %s\n", "Failed to connect to RabbitMQ", err)
		return
	}

	ch, err = conn.Channel()
	if err != nil {
		log.Printf("%s: %s\n", "Failed to open a channel", err)
		return
	}
}

func checkQueue() {
	for {
		go receiveTest()

		err := SendToMQ(exchange, key, []byte("OK"))
		if err != nil {
			log.Println(err.Error())
			connect()
			time.Sleep(time.Second * 30)
			continue
		}

		time.Sleep(time.Minute * 3)
	}
}

func receiveTest() {
	if exist {
		return
	}

	exist = true
	defer func() { exist = false }()

	test, err := ReceiveFromMQ(exchange, key, queue, nil)
	if err != nil {
		log.Println(err.Error())
		return
	}

	for v := range test {
		if string(v.Body) != "OK" {
			log.Println(string(v.Body))
			return
		}

		v.Ack(false)
	}

}

//CloseQueue
func CloseQueue(queue string) {
	num, err := ch.QueueDelete(queue, false, false, false)
	if err != nil {
		log.Println(err.Error())
		return
	}

	log.Println("delete queue ", queue, num)

}
