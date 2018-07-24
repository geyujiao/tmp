package spcard

import (
	"fmt"
	"log"
	"time"
)

var (
	spcardResult = make(chan MsgResult, 10000)

	MsgParseFailed = MsgResult{
		Result:  ChargeMsgFailed,
		ErrInfo: "未收到短信，或短信解析失败",
	}
)

//Listen 监听
//一直监听模块收到的短信，把相关的信息返回给充值结果
func Listen() {
	for {
		receiveMsgBySpcard()

		time.Sleep(time.Second * 2)
	}
}

//ChargeBySpcard 通过新酷卡连接猫池进行充值
func ChargeBySpcard(orderNo, oilCard, chargeCardPwd string) (err error) {
	if len(oilCard) < 6 {
		return fmt.Errorf("加油卡号不正确")
	}

	msgInfo := chargeMsg{
		MsgInfo:         fmt.Sprintf("DS05-%s#%s", oilCard, chargeCardPwd),
		DestPhoneNumber: operatorNumber,
	}

	if err = setReflect(oilCard, chargeCardPwd, orderNo); err != nil {
		log.Println(err.Error())
		return
	}

	if err = sendMsgBySpcard(msgInfo); err != nil {
		log.Println("error msg:", msgInfo)
		return
	}

	return
}

//ReceiveResultBySpcard
func ReceiveResultBySpcard() <-chan MsgResult {

	return (<-chan MsgResult)(spcardResult)
}

//setReflect
func setReflect(cardNo, pwd, orderNo string) (err error) {
	ref := CardReflect{
		CardNo:     cardNo,
		Pwd:        pwd,
		CardSuffix: cardNo[len(cardNo)-6:],
		OrderNo:    orderNo,
	}

	return ref.NewReflect()
}

//sendMsgBySpcard
func sendMsgBySpcard(info chargeMsg) (err error) {
	log.Println("send info", info)

	return send(info)
}

func receiveMsgBySpcard() {
	sr := new(SmsRecv)
	msgs := sr.DealProcSms()

	for _, v := range msgs {
		log.Println("start deal sms ", v)
		receive(v)
	}
}
