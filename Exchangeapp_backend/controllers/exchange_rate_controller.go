package controllers

import (
	"exchangeapp/global"
	"exchangeapp/models"
	// "fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func CreateExchangeRate(ctx *gin.Context) {
	var exchangeRate models.ExchangeRate

	//绑定ExchangeRate汇率模型
	if err := ctx.ShouldBind(&exchangeRate); err != nil {
		ctx.JSON(http.StatusInternalServerError,gin.H{"error":err.Error()})
		return
	}

	//时间赋值
	exchangeRate.Date = time.Now()

	//创建数据库表(汇率实体)
	if err := global.Db.AutoMigrate(&exchangeRate);err!=nil{
		ctx.JSON(http.StatusInternalServerError,gin.H{"error":err.Error()})
		return
	}
	//创建汇率表数据
	if err := global.Db.Create(&exchangeRate).Error;err != nil{
		ctx.JSON(http.StatusInternalServerError,gin.H{"error":err.Error()})
		return
	}
	//test:ctx.Get()自从中间件AuthMiddleWare里set后，可以在 同一条请求链上 的所有函数里取出
	//fmt.Println(ctx.Get("username"))	//能打印 "Kalavinka" true

	//返回
	ctx.JSON(http.StatusOK,exchangeRate)
}

func GetExchangeRates(ctx *gin.Context){
	var exchangeRates []models.ExchangeRate
	if err := global.Db.Find(&exchangeRates).Error;err != nil {
		ctx.JSON(http.StatusInternalServerError,gin.H{"error":err.Error()})
		return
	}
	ctx.JSON(http.StatusOK,exchangeRates)
}