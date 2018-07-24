package spcard

import (
	"log"
	"strings"
)

const (
	ChargeMsgSuccess = "charge_success"
	ChargeMsgFailed  = "charge_failed"
	ChargeMsgLost    = "charge_lost"
	RetryKeyWord     = "]充值失败"
	MsgTailFlag      = "【中国石化】"
	defaultMoney     = "0"

	MsgErrInvalidPwd = "invalid pwd"
	MsgErrPwdNoMoney = "there is no money in this pwd"
)

type (
	MsgResult struct {
		OrderNo string
		OilCard string
		Pwd     string
		Result  string
		Money   string
		Message string
		ErrInfo string
	}
)

func msgParse(msgInfo string) (result MsgResult) {
	if !strings.Contains(msgInfo, MsgTailFlag) {
		log.Println("msg info is : ", msgInfo)
		return MsgParseFailed
	}

	return setResult(msgInfo)
}

/*
一共有三种类型的短信
1、【中国石化】您尾号为 776323的加油卡于 12月29日 14时56分充值成功，金额100元，订单号：2417122914560481
2、【中国石化】加油卡[1000113300007776323]充值失败
3、【中国石化】充值卡卡密[31827502865321518755]验证失败，请重新发送充值短信
4、【中国石化】短信充值指令为：DS05-加油卡卡号#充值卡卡密,请重新发送充值短信,请注意加油卡卡号为19位数字,充值卡卡密为20位数字!
5、【中国石化】加油卡[1000113300007776323]验证失败，只能给主卡并且卡状态正常的加油卡充值，请重新发送充值短信！
*/

func setResult(msgInfo string) (result MsgResult) {
	msgInfo = strings.TrimSpace(msgInfo)
	msgInfo = strings.Replace(msgInfo, " ", "", -1)

	result.Money = "0"

	rs := strings.Index(msgInfo, "充值成功")
	st := strings.Index(msgInfo, "尾号为")
	end := strings.Index(msgInfo, "的加油卡")

	if rs != -1 && len(msgInfo) > 80 && st != -1 && end != -1 {
		return setResultSuccess(msgInfo[st+9:end], msgInfo)

	}

	rf := strings.Index(msgInfo, "充值失败")
	ocId := strings.Index(msgInfo, "加油卡[")

	if rf != -1 && ocId != -1 {
		return setResultFailed(msgInfo[ocId+10:rf-1], msgInfo)

	}

	rf = strings.Index(msgInfo, "验证失败")
	pwd := strings.Index(msgInfo, "充值卡卡密[")
	if rf != -1 && pwd != -1 {
		cr := new(CardReflect)
		cr.Pwd = msgInfo[pwd+16 : rf-1]
		cr.FindByPwd()
		if cr.CardNo != "" {
			return setResultFailed(cr.CardNo, msgInfo)
		}

		log.Println("验证失败且格式解析错误", msgInfo)

	}

	if rf != -1 && ocId != -1 {
		return setResultFailed(msgInfo[ocId+10:rf-1], msgInfo)
	}

	log.Println("异常短信", msgInfo)

	return
}

func setResultSuccess(num string, msgInfo string) (r MsgResult) {
	sm := strings.Index(msgInfo, "金额")
	em := strings.Index(msgInfo, "元")
	if sm > 0 && em > 0 {
		r.Money = msgInfo[sm+6 : em]
	}

	card := new(CardReflect)
	card.CardSuffix = num
	card.FindByCardSuffix()
	if card.CardNo == "" {
		card.CardNo = "*" + num
		log.Println("未找到完整卡号", num)
	}

	r.OrderNo = card.OrderNo
	r.OilCard = card.CardNo
	r.Pwd = card.Pwd
	r.Result = ChargeMsgSuccess
	r.Message = msgInfo

	return
}

func setResultFailed(cardId string, msgInfo string) (r MsgResult) {
	card := new(CardReflect)
	card.CardNo = cardId
	card.FindByCardNo()

	r.OrderNo = card.OrderNo
	r.OilCard = cardId
	r.Pwd = card.Pwd
	r.Result = ChargeMsgFailed
	r.ErrInfo = msgInfo

	return
}
