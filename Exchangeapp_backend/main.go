package main

import (
	"exchangeapp/config"
	"exchangeapp/router"

)

func main() {
	config.InitConfig()
	r := router.StaupRouter()
	port := config.AppConfig.App.Port
	if port == "" {
		port = ":8080"
	}
	// fmt.Println(config.AppConfig.App.Port)

	// type Info struct{
	// 	Message string
	// }
	// infoTest := Info{
	// 	Message : "pooom",
	// }
	// r := gin.Default()
	// r.GET("/ping", func(c *gin.Context) {
	// 	c.JSON(200, infoTest)
	// })
	r.Run(port) //可以设置要监听的端口
}
