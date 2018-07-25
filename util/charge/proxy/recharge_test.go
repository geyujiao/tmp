package proxy

import (
	"testing"
)

func TestRecharge(t *testing.T) {
	dataMessage := DataMessage{}
	dataMessage.TradeOrderId = "2018071949286878"
	dataMessage.CardNo = "1000113300015785355"
	dataMessage.CardPwd = "37959739494645461781"
	dataMessage.Amount = "1"

	//result := Recharge("2018071949286878", "1000113300015785355", "37959739494645461781", "1")

	result, err := RechargeInformation("1045", "2018071949286878", 0)

	t.Error(result, err)
}
