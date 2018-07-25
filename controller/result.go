package controller

import (
	"encoding/json"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/vgmdj/tmp/util/mq/rabbitmq"
	"github.com/vgmdj/utils/logger"
	"net/http"
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
func ResultStatistics(c *gin.Context) {

	logger.Info(c.Request.Body)
	rechargeData := proxy.ParseData(c.Request.Body)
	logger.Info(rechargeData)

	result := MsgResult{
		OrderNo: rechargeData.Orderid,
		Money:   rechargeData.Amount,
		Message: rechargeData.Msg,
		ErrInfo: rechargeData.Msg,
		//OilCard: "",
		//Pwd:     "",
	}
	if rechargeData.Status == proxy.ResponseSuccess {
		result.Result = ChargeMsgSuccess

	} else if rechargeData.Status == proxy.ResponseFail {
		result.Result = ChargeMsgFailed
	} else {
		logger.Error("invalid status")
		return

	}

	err := sendMsgToResultQue(result)
	if err != nil {
		logger.Error(err.Error())
		return
	}

	c.String(http.StatusOK, "ok")

}

//sendMsgToResultQue
func sendMsgToResultQue(result MsgResult) (err error) {
	body, err := json.Marshal(result)
	if err != nil {
		log.Println(err.Error(), result)
		return
	}

	//if result.OrderNo == "" && result.OilCard == "" {
	//	log.Println(result)
	//	return fmt.Errorf("invalid result")
	//}

	err = rabbitmq.SendToMQ(resultExc, resultKey, body)
	if err != nil {
		log.Println(err.Error())
		return
	}

	return
}
