package controller

import (
	"encoding/json"
	"log"

	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/vgmdj/tmp/util/charge/proxy"
	"github.com/vgmdj/tmp/util/mq/rabbitmq"
	"github.com/vgmdj/utils/encrypt"
	"github.com/vgmdj/utils/logger"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
)

var (
	resultExc = "result_exc"
	resultKey = "result_key"
)

type MsgResult struct {
	OrderNo    string
	OilCard    string
	Pwd        string
	Result     string
	Money      string
	OutOrderNo string
	Message    string
	ErrInfo    string
}

//ResultStatistics 充值结果
func ResultStatistics(c *gin.Context) {

	result, err := parseData(c.Request.Body)
	if err != nil {
		logger.Error(err.Error())
		return
	}

	err = sendMsgToResultQue(result)
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

func parseData(body io.Reader) (result MsgResult, err error) {
	tmp, err := ioutil.ReadAll(body)
	if err != nil {
		logger.Error(err.Error())
		return
	}

	values, err := url.ParseQuery(string(tmp))
	if err != nil {
		logger.Error(err.Error())
		return
	}

	//Code        string  //状态码 0
	//Status      string //订单状态
	//Msg         string //订单状态说明
	//Amount      string //实际金额
	//Orderid     string //订单号
	//OutOrderNo  string //官方流水号
	//From_status string

	tmpRes := make(map[string]string)
	sign := ""
	keys := []string{"Orderid", "amount", "code", "from_status", "msg", "outOrderNo", "status"}
	for _, v := range keys {
		var (
			value []string
			ok    bool
		)
		if value, ok = values[v]; !ok {
			logger.Error("no such key", v, values)
			continue
		}

		tmpRes[v] = value[0]
		sign = fmt.Sprintf("%s%s%s", sign, v, value[0])
	}

	sign = encrypt.Md5(sign + proxy.NewClient().AppMd5Secret)
	if sign != tmpRes["sign"] {
		logger.Error("invalid sign ", values, sign)
	}

	result = MsgResult{
		OrderNo:    tmpRes["Orderid"],
		Money:      tmpRes["amount"],
		Message:    tmpRes["msg"],
		ErrInfo:    tmpRes["msg"],
		OutOrderNo: tmpRes["outOrderNo"],
	}

	if tmpRes["status"] == proxy.ResponseSuccess {
		result.Result = ChargeMsgSuccess

	} else if tmpRes["status"] == proxy.ResponseFail {
		result.Result = ChargeMsgFailed

	} else {
		return result, fmt.Errorf("invalid status")
	}

	return
}
