package spcard

import (
	"fmt"
	"log"

	"rongshu.tech/oiling/proxy_gateway/config"
)

const (
	SmsRecvProcStatusDone    = 0
	SmsRecvProcStatusPending = 1
	SmsRecvProcStatusDoing   = 2
)

type SmsRecv struct {
	ID            int64  `gorm:"Column:Id"`
	PortNum       int64  `gorm:"Column:PortNum"`
	PhoNum        string `gorm:"Column:PhoNum"`
	IMSI          string `gorm:"Column:IMSI"`
	ICCID         string `gorm:"Column:ICCID"`
	SmsDate       string `gorm:"Column:smsDate"`
	SmsNumber     string `gorm:"Column:smsNumber"`
	SmsContent    string `gorm:"Column:smsContent"`
	SmsProcStatus int64  `gorm:"Column:smsProcStatus"`
}

func (sr SmsRecv) TableName() string {
	return "sms_recv"
}

func (sr SmsRecv) DealProcSms() []SmsRecv {
	srs := []SmsRecv{}
	err := config.NewMysql().Where("smsNumber = ? and  smsProcStatus = ?",
		operatorNumber, SmsRecvProcStatusPending).Find(&srs).Error
	if err != nil {
		log.Println(err.Error())
	}

	ids := []int64{}
	for _, v := range srs {
		ids = append(ids, v.ID)
	}

	err = config.NewMysql().Model(SmsRecv{}).Where("Id IN (?)", ids).
		Update("smsProcStatus", SmsRecvProcStatusDoing).Error
	if err != nil {
		log.Println(err.Error())
	}

	return srs

}

func (sr *SmsRecv) SetProcSmsStatus() {
	log.Println("set proc sms status ", sr)

	err := config.NewMysql().Model(SmsRecv{}).Where("id = ?", sr.ID).
		Update("smsProcStatus", SmsRecvProcStatusDone).Error
	if err != nil {
		log.Println(err.Error())

	}

}

type SmsSend struct {
	ID         int64  `gorm:"Column:Id"`
	PortNum    int64  `gorm:"Column:PortNum"`
	SmsNumber  string `gorm:"Column:smsNumber"`
	SmsSubject string `gorm:"Column:smsSubject"`
	SmsContent string `gorm:"Column:smsContent"`
	SmsType    int64  `gorm:"Column:smsType"`
	PhoNum     string `gorm:"Column:PhoNum"`
	SmsState   int64  `gorm:"Column:smsState"`
}

func (ss SmsSend) TableName() string {
	return "sms_send"
}

func (ss *SmsSend) NewSendMsg() (err error) {
	check := config.NewMysql().NewRecord(ss)
	if !check {
		fmt.Errorf("cardNo is already exist")
		return
	}

	err = config.NewMysql().Create(ss).Error
	if err != nil {
		log.Println(err.Error())
	}
	return
}

type CardReflect struct {
	OrderNo    string `gorm:"Column:OrderNo;primary_key"`
	CardNo     string `gorm:"Column:CardNo"`
	Pwd        string `gorm:"Column:Pwd"`
	CardSuffix string `gorm:"Column:CardSuffix"`
}

func (cr CardReflect) TableName() string {
	return "card_reflect"
}

//NewReflect 如果卡号相同，订单号相同，则不做任何处理，直接返回
//如果卡号相同，订单号不同，则返回错误，已存在
//如果卡号不同，订单号不同，则插入新的映射，返回
func (cr *CardReflect) NewReflect() (err error) {
	var dbcr []CardReflect
	err = config.NewMysql().Model(&CardReflect{}).Where("CardSuffix = ? or Pwd = ?",
		cr.CardSuffix, cr.Pwd).Find(&dbcr).Error
	if err != nil || len(dbcr) > 1 || (len(dbcr) == 1 && dbcr[0].OrderNo != cr.OrderNo) {
		log.Println(cr.CardNo, len(dbcr))
		return fmt.Errorf("cardNo is already exist")
	} else if len(dbcr) == 1 && dbcr[0].OrderNo == cr.OrderNo {
		return
	}

	err = config.NewMysql().Create(cr).Error
	if err != nil {
		log.Println(err.Error())
		return
	}

	return
}

func (cr *CardReflect) FindByPwd() {
	err := config.NewMysql().Where("Pwd = ?", cr.Pwd).First(cr).Error
	if err != nil {
		log.Println(err.Error())
		return
	}
}

func (cr *CardReflect) FindByCardSuffix() {
	err := config.NewMysql().Where("CardSuffix = ?", cr.CardSuffix).First(cr).Error
	if err != nil {
		log.Println(err.Error())
		return
	}
}

func (cr *CardReflect) FindByCardNo() {
	err := config.NewMysql().Where("CardNo = ?", cr.CardNo).First(cr).Error
	if err != nil {
		log.Println(err.Error())
		return
	}
}

func (cr *CardReflect) Pop() (err error) {
	err = config.NewMysql().Where("OrderNo = ?", cr.OrderNo).Delete(&CardReflect{}).Error
	if err != nil {
		log.Println(err.Error())
		return
	}
	return
}
