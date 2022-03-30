package model

import 	"gorm.io/gorm"

func (Message) TableName() string {
	return "message"
}

type Message struct {
	Id        int    `gorm:"primary_key,type:INT;not null;AUTO_INCREMENT"`
	User_Id   int    `json:"User_Id"  binding:"required"`
	Content   string `json:"Content"  binding:"required"`
	Version   int    `gorm:"default:0"`
	// 包含 CreatedAt 和 UpdatedAt 和 DeletedAt 欄位
    gorm.Model
}