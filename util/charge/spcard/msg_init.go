package spcard

import (
	"log"
	"rongshu.tech/oiling/proxy_gateway/config"
)

const (

	//chinaMobileChargeNumber 中国移动充值
	chinaMobileChargeNumber = "106575650318"
	//chinaUnicomChargeNumber 中国联通充值
	chinaUnicomChargeNumber = "1065502795105888"
	//chinaTelecomChargeNumber 中国电信充值
	chinaTelecomChargeNumber = "10659057195105888"

	testChargeNumber = "18921776016"
)

var (
	NoGsmModemConnections = chargeResult{9901, "当前没有短信猫连接存活，请稍后再试"}
	NoChargeMsg           = chargeResult{9902, "充值失败，等待超时，无响应"}
	ChargeFailed          = chargeResult{9903, "充值交易失败"}
	Success               = chargeResult{0, "OK"}

	operatorNumber = testChargeNumber
)

type (
	chargeMsg struct {
		MsgInfo         string
		DestPhoneNumber string
	}

	chargeResult struct {
		Code int64
		Msg  string
	}
)

func Init() {
	log.Println("start init")

	log.SetFlags(log.Ldate | log.Lshortfile)

	if config.Spcard.ChargeNumber != "" {
		operatorNumber = config.Spcard.ChargeNumber
	}

	log.Println("set charge number :", operatorNumber)

	go Listen()

}

//send
func send(msg chargeMsg) (err error) {
	log.Println("prepare to send a message", msg)

	ss := SmsSend{
		SmsNumber:  msg.DestPhoneNumber,
		SmsContent: msg.MsgInfo,
		SmsType:    0,
	}

	return ss.NewSendMsg()
}

//receive read msg that status is no proc from mysql
//TODO copy the failed msg to err queue
func receive(recv SmsRecv) {
	r := msgParse(recv.SmsContent)
	if r.OilCard == "" {
		log.Println("unexpect err , cant's parse msg info")
	}

	spcardResult <- r

	recv.SmsProcStatus = SmsRecvProcStatusDone

	recv.SetProcSmsStatus()

	if r.OrderNo == "" {
		log.Println("error parse ", recv)
		return
	}
}

func (cr chargeResult) Error() string {
	return cr.Msg
}
