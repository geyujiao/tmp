package controller

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/vgmdj/utils/logger"
	"github.com/vgmdj/tmp/util/mq/rabbitmq"
)

var (
	resultExc = "result_exc"
	resultKey = "result_key"
)

type MsgResult struct {
	OrderNo string
	OilCard string
	Pwd     string
	Result  string
	Money   string
	Message string
	ErrInfo string
}

//ResultStatistics 充值结果
func ResultStatistics() {
	result := MsgResult{
		OrderNo: "",
		OilCard: "",
		Pwd:     "",
		Money:   "",
		Result:  ChargeMsgSuccess,
		Message: "",
		ErrInfo: "",
	}

	err := sendMsgToResultQue(result)
	if err != nil {
		logger.Error(err.Error())
		return

	}

	result = MsgResult{
		OrderNo: "",
		OilCard: "",
		Pwd:     "",
		Money:   "",
		Result:  ChargeMsgFailed,
		Message: "",
		ErrInfo: "",
	}

	err = sendMsgToResultQue(result)
	if err != nil {
		logger.Error(err.Error())
		return
	}

}

//sendMsgToResultQue
func sendMsgToResultQue(result MsgResult) (err error) {
	body, err := json.Marshal(result)
	if err != nil {
		log.Println(err.Error(), result)
		return
	}

	if result.OrderNo == "" && result.OilCard == "" {
		log.Println(result)
		return fmt.Errorf("invalid result")
	}

	err = rabbitmq.SendToMQ(resultExc, resultKey, body)
	if err != nil {
		log.Println(err.Error())
		return
	}

	return
}
