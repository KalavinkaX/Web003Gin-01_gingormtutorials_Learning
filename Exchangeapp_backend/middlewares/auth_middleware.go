package middlewares

import (
	"exchangeapp/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

//发起请求前验证Token是否合规的中间件
func AuthMiddleWare() gin.HandlerFunc {//返回值是个匿名函数
	return func(ctx *gin.Context) {
		token := ctx.GetHeader("Authorization")
		if token == "" {
			ctx.JSON(http.StatusUnauthorized,gin.H{"error":"Missing Authorization!"})
			// 2b. 调用 ctx.Abort()，这是中间件的精髓！
			ctx.Abort()// 它会立即终止这个请求的处理链（但当前匿名函数能执行完）。后续的任何中间件或业务处理器都不会被执行。
			return
		}

		//token不为""，解析jwt看是否符合标准
		username,err := utils.ParseJWT(token)
		if err != nil {
			ctx.JSON(http.StatusUnauthorized,err.Error())
			ctx.Abort()//作用：终止接下来的中间件，当前函数还会执行完？
			return
		}
		ctx.Set("username",username)//类似springboot的threadlocal？
		ctx.Next()// 调用 ctx.Next() 将控制权交给下一个处理器
	}
}