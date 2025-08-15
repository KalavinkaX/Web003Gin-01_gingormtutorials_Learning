package router

import (
	"exchangeapp/controllers"
	"net/http"

	"github.com/gin-gonic/gin"
)

func StaupRouter() *gin.Engine {//返回Gin路由引擎指针，用于main等函数使用
	r := gin.Default()//创建一个Gin路由引擎

	//路由组：统一管理"/api/auth"开头的请求路径
	auth := r.Group("/api/auth")
	{	
		//具体POST方法的编写。参数为(路径,函数名)
		auth.POST("/login",controllers.Login)

		auth.POST("/register",controllers.Register)
	}
	test := r.Group("/test")
	{
		test.GET("/test1",func(ctx *gin.Context)  {//函数参数必须接收一个*gin.Context 类型的参数
			ctx.JSON(http.StatusOK,gin.H{"msg" : "Test1"})
		})
	}
	return r

}