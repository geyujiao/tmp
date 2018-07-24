package proxy

import "testing"

func TestRecharge(t *testing.T) {
	dataMessage := DataMessage{}
	dataMessage.TradeOrderId = "2018071949286878"
	dataMessage.CardNo = "1000113300015785355"
	dataMessage.CardPwd = "37959739494645461781"
	dataMessage.Amount = "1"

	// data ok
	//result, err := getData(dataMessage, "zhwwtoo786bbsome")
	//t.Error(result)
	//if err != nil{
	//	t.Error("1111")
	//	t.Error(err.Error())
	//}

	////sign  ok
	//result, err := getSign("UEK4bXlGpKIgIdr0N9UogQ==", "o61uswq6")
	//t.Error(result)
	//if err != nil{
	//	t.Error("1111")
	//	t.Error(err.Error())
	//}

	result := Recharge("2018071949286878", "1000113300015785355", "37959739494645461781", "1")

	t.Error(result)
}
