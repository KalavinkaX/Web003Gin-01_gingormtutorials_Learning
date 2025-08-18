package controllers

import (
	"encoding/json"
	"errors"
	"exchangeapp/config"
	"exchangeapp/global"
	"exchangeapp/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"gorm.io/gorm"
)

var cacheKey = "articles"

// 插入文章
func CreateArticle(ctx *gin.Context) {
	var article models.Article
	//基础模型绑定
	if err := ctx.ShouldBindJSON(&article); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	//数据库自动迁移
	if err := global.Db.AutoMigrate(&article); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	//插入数据
	if err := global.Db.Create(&article).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	//新增数据后删除缓存
	if err := global.RedisDB.Del(cacheKey).Err(); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}
	//最后成功插入Article数据
	ctx.JSON(http.StatusOK, article)
}

// 获取所有文章
func GetArticles(ctx *gin.Context) {
	cacheData, err := global.RedisDB.Get(cacheKey).Result()
	//1.redis中取不到数据，需要从DB中取
	if err == redis.Nil { //key在Redis中不存在
		var articles []models.Article
		//链式操作: GORM 的几乎所有操作（Find, Create, Where, First, Update 等）都是“链式”的，这意味着每个方法都会返回一个 *gorm.DB 类型的结构体对象。
		if err := global.Db.Find(&articles).Error; err != nil { //最后的.Error是获取*gorm.DB结构体的Error属性的内容
			if errors.Is(err, gorm.ErrRecordNotFound) {
				ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			} else {
				ctx.JSON(http.StatusInternalServerError, gin.H{"err": err.Error()})
			}
			return
		}
		//序列化对象
		articleJson, err := json.Marshal(cacheData)
		if err != nil { //序列化Redis数据为对象时出错
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err})
			return
		}
		//缓存DB查询的数据到Redis
		if err := global.RedisDB.Set(cacheKey, articleJson, 10*time.Minute).Err(); err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusOK, articles)

		//2.从Redis取缓存时出错
	} else if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return

		//3.成功从Redis中取出缓存数据
	} else {
		var articles []models.Article
		//反序列化Redis数据
		if err := json.Unmarshal([]byte(cacheData), &articles); err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusOK, articles)
	}
}

func GetArticleById(ctx *gin.Context) {
	//！通过Param取出url里的参数
	id := ctx.Param("id")
	var article models.Article
	//！健壮性：错误类型准确检查
	if err := global.Db.Where("id = ?", id).First(&article).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			//找不到id为?的数据，返回准确err类型
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	ctx.JSON(http.StatusOK, article)

}
