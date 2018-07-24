package controller

import (
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/streadway/amqp"
	"github.com/vgmdj/utils/logger"
	"github.com/vgmdj/tmp/config"
	"github.com/vgmdj/tmp/util/charge/proxy"
	"github.com/vgmdj/tmp/util/mq/rabbitmq"
)

var (
	//单位毫秒，设置为2.5分钟

	chargeExc = ""
	chargeQue = ""
	chargeKey = ""
)

var (
	ChargeMsgSuccess = "charge_success"
	ChargeMsgFailed  = "charge_failed"

	ackMsg sync.Map
)

type (
	RechargeQue struct {
		OrderNo       string
		OilCard       string
		ChargeCardPwd string
		Amount        string
	}
)

func Init() {
	chargeExc = fmt.Sprintf("%s_exchange", config.Rabbitmq.RechargeQue)
	chargeQue = fmt.Sprintf("%s_queue", config.Rabbitmq.RechargeQue)
	chargeKey = fmt.Sprintf("%s_processing", config.Rabbitmq.RechargeQue)

	resultExc = fmt.Sprintf("%s_exc", config.Rabbitmq.ResultQeu)
	resultKey = fmt.Sprintf("%s_key", config.Rabbitmq.ResultQeu)

	log.Println(chargeExc, resultExc)
}

//MsgProcessing 监听充值队列
//每当有新的充值指令进来后，开始执行
func MsgProcessing() {
	msgs, err := rabbitmq.ReceiveFromMQ(chargeExc, chargeKey, chargeQue, nil)
	if err != nil {
		log.Println(err.Error())
		SetHealth(err)
		return
	}

	for delivery := range msgs {
		log.Printf("Received a message: %s", delivery.Body)
		order, err := sendMsgToWorkSpace(delivery.Body)
		if err == nil {
			continue
		}

		log.Println(err.Error())

		result := MsgResult{
			Result:  ChargeMsgFailed,
			OilCard: order.OilCard,
			Pwd:     order.ChargeCardPwd,
			Money:   "0",
			Message: "",
			ErrInfo: err.Error(),
		}

		err = sendMsgToResultQue(result)
		if err != nil {
			logger.Error(err.Error())
			continue
		}

		delivery.Ack(false)
	}

}

//TODO 计数加失败处理
func cacheQue(delivery amqp.Delivery) {
	log.Println("发送失败，缓存30s后再次进入队列尝试： ", string(delivery.Body))
	time.Sleep(time.Second * 30)

	log.Println("重新进入队列： ", string(delivery.Body))
	delivery.Reject(true)
}

//parseDeliveryBody
func parseDeliveryBody(msg []byte) (q RechargeQue, err error) {
	err = json.Unmarshal(msg, &q)
	if err != nil {
		log.Println(err.Error())
		return
	}

	return
}

//sendMsgToWorkSpace
func sendMsgToWorkSpace(msg []byte) (order RechargeQue, err error) {
	order, err = parseDeliveryBody(msg)
	if err != nil {
		log.Println(err.Error())
		return
	}

	if err = proxy.Recharge(order.OrderNo, order.OilCard, order.ChargeCardPwd, order.Amount); err != nil {
		log.Println(err.Error())
		return
	}

	return
}
