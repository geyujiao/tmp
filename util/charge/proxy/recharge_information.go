package proxy

import (
	"encoding/json"
	"fmt"
	"github.com/vgmdj/utils/logger"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"
)

type RechargeInformaRequest struct {
	AccountId string `json:"accountid"` //账户id
	OrderId   string `json:"orderid"`   //（自定义订单号，同trade_order_id）
}
type RechargeInformaResponse struct {
	Code        int    //状态码 0
	Status      string //订单状态
	Msg         string //订单状态说明
	Amount      string //实际金额
	Orderid     string //订单号
	OutOrderNo  string //官方流水号
	From_status string
}
type ErrorResult struct {
	Code int    //状态码 0
	Msg  string //订单状态说明
}

func (err *ErrorResult) Error() string {
	return err.Msg
}

const (
	ServerExp       = iota //服务器异常
)
const (
	ResponseSuccess = "2" //返回结果：充值成功
	ResponseFail    = "3" //返回结果：充值失败
)

func RechargeInformation(accountId, orderId string, notify int) (response *RechargeInformaResponse, err error) {
	var (
		app = NewClient()
	)
	app.SetUrl("http://47.93.136.39:8170/api/charge_status")
	body := make(map[string]string)
	body["accountid"] = accountId
	body["orderid"] = orderId
	body["notify"] = strconv.Itoa(notify)

	response, err = GetRechargeInformation(app.Url, body)
	return response, err
}

func GetRechargeInformation(baseUrl string, body map[string]string) (result *RechargeInformaResponse, err error) {
	result = new(RechargeInformaResponse)
	errResult := new(ErrorResult)
	urlQuery := "?"
	for k, v := range body {
		urlQuery += fmt.Sprintf("%s=%s", k, v)
		urlQuery += "&"
	}
	url := baseUrl + urlQuery
	resp, err := http.Get(url)
	defer resp.Body.Close()
	if err != nil {
		logger.Error(err.Error())
		errResult.Code = ServerExp
		errResult.Msg = "服务器异常"
		return nil, errResult
	}
	if body["notify"] == "1" {
		//回调模式
		logger.Info(body["notify"])
		return nil, nil
	}
	result = ParseData(resp.Body)
	return result, err
}

func ParseData(data io.ReadCloser) (result *RechargeInformaResponse) {
	result = new(RechargeInformaResponse)
	bytes, err := ioutil.ReadAll(data)
	if err != nil {
		logger.Error(err.Error(), data)
		return nil
	}
	err = json.Unmarshal(bytes, result)
	if err != nil {
		logger.Error(err.Error(), string(bytes))
		return nil
	}
	return result
}
