package proxy

import (
	"crypto/md5"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/vgmdj/utils/httplib"
	"github.com/vgmdj/utils/logger"
	"strings"
)

type DataMessage struct {
	TradeOrderId string `json:"trade_order_id"` //自定义订单号
	CardNo       string `json:"card_no"`        //加油卡卡号
	CardPwd      string `json:"card_pwd"`       //充值卡卡密
	Amount       string `json:"amount"`         //充值卡金额
}

type Response struct {
	Code string
	Msg  string
}
type RechargeResult struct {
	MsgType int
	Msg     string
}

func (err *RechargeResult) Error() string {
	return err.Msg
}

func getData(dataM DataMessage, appSecret string) (data string, err error) {
	//to do
	//Base64_encode(AES(Json_encode(消息体)))
	//AES加密方式：CBC128位/PKCS7(JAVA PKCS5)/iv(1234567890123456)

	bytes, err := json.Marshal(dataM)
	if err != nil {
		return "", err
	}
	strMesg := string(bytes)

	//AES
	arrEncrypt, err := AesEncrypt([]byte(strMesg), []byte(appSecret))
	if err != nil {
		return "", err
	}

	data = base64.StdEncoding.EncodeToString(arrEncrypt)
	return data, err
}

func getSign(data, appmd5secret string) (sign string, err error) {
	//to do
	//Lower(MD5(data + appmd5secret))
	str := data + appmd5secret
	md5Bytes := md5.Sum([]byte(str))
	md5Str := fmt.Sprintf("%x", md5Bytes)
	sign = strings.ToLower(md5Str)

	return sign, nil
}

func Recharge(tradeOrderId, cardNo, cardPwd, amount string) error {
	var (
		resp        = Response{}
		app         = NewClient()
		dataMessage = DataMessage{}
		dataStr     string
	)
	result := new(RechargeResult)

	dataMessage.TradeOrderId = tradeOrderId
	dataMessage.CardNo = cardNo
	dataMessage.CardPwd = cardPwd
	dataMessage.Amount = amount
	dataStr, err := getData(dataMessage, app.AppSecret)
	if err != nil {
		logger.Error(err.Error())
		result.MsgType = SendFail
		result.Msg = "数据加密异常"
		return result
	}

	sign, err := getSign(dataStr, app.AppMd5Secret)
	if err != nil {
		logger.Error(err.Error())
		result.MsgType = SendFail
		result.Msg = "数据加密异常"
		return result
	}

	app.SetUrl("http://47.93.136.39:8170/api/charge_add")

	body := make(map[string]string)
	body["accountid"] = app.AppId
	body["data"] = dataStr
	body["sign"] = sign

	headers := make(map[string]string)
	headers[httplib.ResponseResultContentType] = "application/json"

	err = httplib.PostForm(app.Url, &resp, body, headers)
	if err != nil {
		logger.Error(err.Error())
		result.MsgType = SendFail
		result.Msg = "服务器异常"
		return result
	}

	if resp.Code == "0" {
		logger.Info("提交成功")
		return nil
	}

	result.MsgType = SendFail
	result.Msg = resp.Msg

	return result
}
