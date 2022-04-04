package router

import (
	"message/controller"
	"github.com/gin-gonic/gin"
)

func SetRouter() *gin.Engine {
	//顯示 debug 模式
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()

	v1 := r.Group("api/v1")
	{
		//新增留言
		v1.POST("/message", controller.Create)
		//查詢全部留言
		v1.GET("/message", controller.GetAll)
		//查詢 {id} 留言
		v1.GET("/message/:id", controller.Get)
		//修改 {id} 留言
		v1.PATCH("/message/:id", controller.Update)
		//刪除 {id} 留言
		v1.DELETE("/message/:id", controller.Delete)
	}
	return r
}
