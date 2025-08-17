package controllers

import (
	"exchangeapp/global"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
)

func LikeArticle(ctx *gin.Context) {
	//获取URL的articleId并拼接为key
	articleId := ctx.Param("id")
	likeKey := "article:" + articleId + ":likes"

	//点赞数++
	//状态型命令 (Status-Returning Commands): 这类命令的主要目的是修改 Redis 中的数据，我们只关心这个操作是否成功。例如 SET, INCR, DEL, EXPIRE 等。
	//用.Err()的原因：Incr 命令本身在 Redis 中执行后，会返回增加后的新值（一个整数）。但在很多业务场景下，我们调用 Incr 只是为了让计数器加一，并不关心加一之后具体等于多少，我们只关心这个“加一”的动作有没有成功。
	if err := global.RedisDB.Incr(likeKey).Err(); err!=nil{
		ctx.JSON(http.StatusInternalServerError,gin.H{"error":err.Error()})
		return
	}

	ctx.JSON(http.StatusOK,gin.H{"message": "Successfully liked the article"})
}

func GetArticleLikes(ctx *gin.Context)  {
	articleId := ctx.Param("id")
	likeKey := "article:" + articleId + ":likes"
	
	//从Redis取点赞数
	//取值型命令: 这类命令的主要目的是从 Redis 获取一个数据。例如 GET, HGET, LLEN, SCARD 等。
	//.Result()的原因：因为您调用 Get 的唯一目的就是为了获取那个值！如果 go-redis 只提供一个 .Err() 方法，您就只能知道操作是否成功，却永远拿不到您想要的数据。
	likesValue,err := global.RedisDB.Get(likeKey).Result()
	if err == redis.Nil{//err == redis.Nil是指redis中key不存在
		likesValue = "0"
	}else if err != nil{
		//
		ctx.JSON(http.StatusInternalServerError,gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK,gin.H{"likes":likesValue})

}