package controller

import (
	"message/model"
	"message/repository"
	"net/http"
	"unicode/utf8"

	"github.com/gin-gonic/gin"
)

func GetAll(c *gin.Context) {
	message, err := repository.GetAllMessage()

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": message})
}

func Get(c *gin.Context) {
	var message model.Message

	if err := repository.GetMessage(&message, c.Param("id")); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "找不到留言"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": message})
}

func Create(c *gin.Context) {
	var message model.Message

	if c.PostForm("Content") == "" || utf8.RuneCountInString(c.PostForm("Content")) >= 20 {
		c.JSON(http.StatusBadRequest, gin.H{"message": "沒有輸入內容或長度超過20個字元"})
		return
	}

	c.Bind(&message)
	repository.CreateMessage(&message)
	c.JSON(http.StatusCreated, gin.H{"message": message})
}

func Update(c *gin.Context) {
	var message model.Message

	if c.PostForm("Content") == "" || utf8.RuneCountInString(c.PostForm("Content")) >= 20 {
		c.JSON(http.StatusBadRequest, gin.H{"message": "沒有輸入內容或長度超過20個字元"})
		return
	}

	if err := repository.UpdateMessage(&message, c.PostForm("Content"), c.Param("id")); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "找不到留言"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": message})
}

func Delete(c *gin.Context) {
	var message model.Message

	if err := repository.DeleteMessage(&message, c.Param("id")); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "找不到留言"})
		return
	}
	c.JSON(http.StatusNoContent, gin.H{"message": "刪除留言成功"})
}
