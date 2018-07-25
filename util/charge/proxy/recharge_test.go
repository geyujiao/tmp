package proxy

import (
	"github.com/vgmdj/tmp/util/charge/proxy"
	"github.com/vgmdj/utils/encrypt"
	"testing"
)

func TestRecharge(t *testing.T) {
	//dataMessage := DataMessage{}
	//dataMessage.TradeOrderId = "2018071949286878"
	//dataMessage.CardNo = "1000113300015785355"
	//dataMessage.CardPwd = "37959739494645461781"
	//dataMessage.Amount = "1"
	//
	////result := Recharge("2018071949286878", "1000113300015785355", "37959739494645461781", "1")
	//
	//result, err := RechargeInformation("1045", "2018071949286878", 0)
	//
	//t.Error(result, err)

	//info   Orderid=2018072498681213&code=0&status=3&msg=%E8%AE%A2%E5%8D%95%E6%8F%90%E4%BA%A4%E5%A4%B1%E8%B4%A5%EF%BC%8C%E8%AF%B7%
	//E9%87%8D%E8%AF%95%E3%80%82&amount=0&from_status=02&outOrderNo=&sign=4036460fdb8d990c9497e04033349762

	str := `Orderid2018072498681213amount0code0from_status02msg订单提交失败，请重试。outOrderNostatus3`
	t.Error(encrypt.Md5(str + proxy.NewClient().AppMd5Secret))
	t.Error(proxy.NewClient().AppMd5Secret)
	//
	str2 := `Orderid=2018072498681213&code=0&status=3&msg=%E8%AE%A2%E5%8D%95%E6%8F%90%E4%BA%A4%E5%A4%B1%E8%B4%A5%EF%BC%8C%E8%AF%B7%
	E9%87%8D%E8%AF%95%E3%80%82&amount=0&from_status=02`
	t.Error(encrypt.Md5(str2 + proxy.NewClient().AppMd5Secret))

	//
	//	str3 := `Orderid=2018072498681213&code=0&status=3&msg=订单提交失败，请重试。&amount=0&from_status=02&sign=4036460fdb8d990c9497e04033349762`
	//	str3 = `{"Orderid":"2018072498681213","code":0,"status":"3","msg":"订单提交失败，请重试。","amount":"0","from_status":"02"}`
	//	t.Error(encrypt.Md5(str3 + proxy.NewClient().AppMd5Secret))
	//
	//	str4 := `Orderid=2018072498681213&code=0&status=3&msg=订单提交失败，请重试。&amount=0&from_status=02`
	//	str4 = `{"Orderid":"2018072498681213","code":0,"status":"3","msg":"订单提交失败，请重试。","amount":"0","from_status":"02","outOrderNo":""}`
	//	t.Error(encrypt.Md5(str4 + proxy.NewClient().AppMd5Secret))
	//
	//	str5 := `{"Orderid":"2018072498681213","code":0,"status":"3","msg":"%E8%AE%A2%E5%8D%95%E6%8F%90%E4%BA%A4%E5%A4%B1%E8%B4%A5%EF%BC%8C%E8%AF%B7%
	//E9%87%8D%E8%AF%95%E3%80%82","amount":"0","from_status":"02","outOrderNo":""}`
	//	t.Error(encrypt.Md5(str5 + proxy.NewClient().AppMd5Secret))

	//index := strings.Index(str, "&sign")
	//var sign string
	//if index == -1 {
	//	t.Error("index==-1")
	//}
	//
	//data := str[:index]
	//t.Error(data)
	//sign = encrypt.Md5(data + proxy.NewClient().AppMd5Secret)

}
