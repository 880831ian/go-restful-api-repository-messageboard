package controller

import (
	"net/http"
	"message/model"
	"message/repository"
	"github.com/gin-gonic/gin"
)

func GetAll(c *gin.Context) {
	message,err := repository.GetAllMessage()

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": message})
}

func Get(c *gin.Context) {
	var message model.Message

	if err := repository.GetMessage(&message, c.Param("id")); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "找不到留言"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": message})
}

func Create(c *gin.Context) {
	var message model.Message
	c.Bind(&message)
	
	if err := repository.CreateMessage(&message); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "錯誤請求"})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": message})
}

func Update(c *gin.Context) {
	var message model.Message
	
	if err := repository.GetMessage(&message, c.Param("id")); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "找不到留言"})
		return
	}
	if err := repository.UpdateMessage(&message, c.PostForm("Content"), c.Param("id"));err!=nil{
		c.JSON(http.StatusBadRequest, gin.H{"error": "錯誤請求"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": message})
}

func Delete(c *gin.Context) {
	var message model.Message

	if err := repository.DeleteMessage(&message, c.Param("id")); err!=nil{
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusNoContent, gin.H{"message": "刪除留言成功"})	
}