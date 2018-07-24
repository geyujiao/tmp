package spcard

import (
	"rongshu.tech/oiling/oiling_gateway/config"
	"testing"
)

func InitDB() {
	config.Mysql.Host = "10.11.22.123"
	config.Mysql.Port = "12306"
	config.Mysql.UserName = "rongshu"
	config.Mysql.Password = "MinkTech2501"
	config.Mysql.DBName = "spcard"
}

func TestReflect(t *testing.T) {
	InitDB()

	cr := new(CardReflect)
	cr.OrderNo = "00ea36406f0811e8a32e00163e0fdf37"
	cr.CardNo = "1000213300007776323"
	cr.Pwd = "33133393015448497595"
	cr.CardSuffix = "776323"

	t.Log(cr.NewReflect())

}

func TestDelete(t *testing.T) {
	InitDB()

	cr := new(CardReflect)
	cr.OrderNo = "00ea36406f0811e8a32e00163e0fdf37"
	cr.CardNo = "1000213300007776323"
	cr.Pwd = "33133393015448497595"
	cr.CardSuffix = "776323"

	t.Log(cr.Pop())

}
