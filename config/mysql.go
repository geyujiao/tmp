package config

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"log"
	"sync"
)

var Mysql mysql
var DB *gorm.DB
var once sync.Once

type mysql struct {
	Host     string `config:"host"`
	Port     string `config:"port"`
	UserName string `config:"username"`
	Password string `config:"password"`
	DBName   string `config:"dbname"`
}

func NewMysql() *gorm.DB {
	once.Do(func() {
		var err error
		dbInfo := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True",
			Mysql.UserName, Mysql.Password, Mysql.Host, Mysql.Port, Mysql.DBName)
		DB, err = gorm.Open("mysql", dbInfo)
		if err != nil {
			log.Println("db info is : ", dbInfo)
			log.Printf(err.Error())
		}

	})

	return DB
}
