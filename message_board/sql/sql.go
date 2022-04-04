package sql

import (
	"fmt"
	"io/ioutil"

	"gopkg.in/yaml.v2"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var Connect *gorm.DB

type conf struct {
	Host     string `yaml:"host"`
	UserName string `yaml:"username"`
	Password string `yaml:"password"`
	DbName   string `yaml:"dbname"`
	Port     string `yaml:"port"`
}

func (c *conf) getConf() *conf {
	//讀取config/connect.yaml檔案
	yamlFile, err := ioutil.ReadFile("sql/connect.yaml")

	//若出現錯誤，列印錯誤訊息
	if err != nil {
		fmt.Println(err.Error())
	}

	//將讀取的字串轉換成結構體conf
	err = yaml.Unmarshal(yamlFile, c)
	if err != nil {
		fmt.Println(err.Error())
	}
	return c
}

//初始化連線資料庫
func InitMySql() (err error) {
	var c conf

	//獲取yaml配置引數
	conf := c.getConf()

	//將yaml配置引數拼接成連線資料庫的url
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		conf.UserName,
		conf.Password,
		conf.Host,
		conf.Port,
		conf.DbName,
	)

	//連線資料庫
	Connect, err = gorm.Open(mysql.New(mysql.Config{DSN: dsn}), &gorm.Config{})
	return
}
