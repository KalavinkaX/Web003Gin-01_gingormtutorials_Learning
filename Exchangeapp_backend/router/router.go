package router

import (
	"exchangeapp/controllers"
	"exchangeapp/middlewares"
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
	
	api := r.Group("/api")
	//以下这个GET请求不会走中间件，直接能发起请求
	api.GET("/exchangeRates",controllers.GetExchangeRates)
	//发起以下Use{}内请求前，需要走中间件函数
	api.Use(middlewares.AuthMiddleWare())
	{
		api.POST("/exchangeRates",controllers.CreateExchangeRate)
	}

	test := r.Group("/test")
	{
		test.GET("/test1",func(ctx *gin.Context)  {//函数参数必须接收一个*gin.Context 类型的参数
			ctx.JSON(http.StatusOK,gin.H{"msg" : "Test1"})
		})
	}
	return r

}