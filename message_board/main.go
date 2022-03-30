package main

import (
	"message/config"
	"message/model"
	"message/routes"
	"fmt"
)

func main() {
	//連線資料庫
	if err := config.InitMySql() ;err != nil {
		panic(err)
	}
	//連結模型
	config.Sql.AutoMigrate(&model.Message{})
	//註冊路由
	r := routes.SetRouter()
	//啟動埠為8081的專案
	fmt.Println("開啟127.0.0.0.1:8081...")
	r.Run("127.0.0.1:8081")
}
