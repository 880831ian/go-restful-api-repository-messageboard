package repository

import (
	"message/config"
	"message/model"
)

//查詢全部留言
func GetAllMessage() (message []*model.Message,err error) {
	if err := config.Sql.Find(&message).Error; err != nil {
		return nil,err
	}
	return
}

//查詢 {id} 留言
func GetMessage(message *model.Message, id string) (err error) {
	if err := config.Sql.Where("id=?", id).First(&message).Error; err != nil {
		return err
	}
	return nil
}

//新增留言
func CreateMessage(message *model.Message) (err error) {
	if err = config.Sql.Create(&message).Error; err != nil {
		return err
	}
	return nil
}

//更新 {id} 留言
func UpdateMessage(message *model.Message, content, id string) (err error) {
	if err = config.Sql.Model(&message).Where("id=?", id).Update("content" ,content).Error; err != nil {
		return err
	}
	return nil
}

//刪除 {id} 留言
func DeleteMessage(message *model.Message, id string) (err error) {
	if err = config.Sql.Where("id=?", id).Delete(&message).Error; err != nil {
		return err
	}
	return nil
}