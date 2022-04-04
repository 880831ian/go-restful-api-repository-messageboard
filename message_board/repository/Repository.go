package repository

import (
	"message/model"
	"message/sql"
)

//查詢全部留言
func GetAllMessage() (message []*model.Message, err error) {
	err = sql.Connect.Find(&message).Error
	return
}

//查詢 {id} 留言
func GetMessage(message *model.Message, id string) (err error) {
	err = sql.Connect.Where("id=?", id).First(&message).Error
	return
}

//新增留言
func CreateMessage(message *model.Message) (err error) {
	err = sql.Connect.Create(&message).Error
	return
}

//更新 {id} 留言
func UpdateMessage(message *model.Message, content, id string) (err error) {
	err = sql.Connect.Where("id=?", id).First(&message).Update("content", content).Error
	return
}

//刪除 {id} 留言
func DeleteMessage(message *model.Message, id string) (err error) {
	err = sql.Connect.Where("id=?", id).First(&message).Delete(&message).Error
	return
}
