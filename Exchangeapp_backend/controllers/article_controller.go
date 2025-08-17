package controllers

import (
	"errors"
	"exchangeapp/global"
	"exchangeapp/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

//插入文章
func CreateArticle(ctx *gin.Context) {
	var article models.Article
	//基础模型绑定
	if err := ctx.ShouldBindJSON(&article); err != nil{
		ctx.JSON(http.StatusBadRequest,gin.H{"error" : err.Error()})
		return
	}
	//数据库自动迁移
	if err := global.Db.AutoMigrate(&article);err != nil{
		ctx.JSON(http.StatusInternalServerError,gin.H{"error":err.Error()})
		return
	}
	//插入数据
	if err := global.Db.Create(&article).Error; err != nil{
		ctx.JSON(http.StatusInternalServerError,gin.H{"error":err.Error()})
		return
	}
	//最后成功插入Article数据
	ctx.JSON(http.StatusOK,article)
}

//获取所有文章
func GetArticles(ctx *gin.Context){
	var articles []models.Article
	//链式操作: GORM 的几乎所有操作（Find, Create, Where, First, Update 等）都是“链式”的，这意味着每个方法都会返回一个 *gorm.DB 类型的结构体对象。
	if err := global.Db.Find(&articles).Error; err != nil{//最后的.Error是获取*gorm.DB结构体的Error属性的内容
		ctx.JSON(http.StatusInternalServerError,gin.H{"err":err.Error()})
		return
	}
	ctx.JSON(http.StatusOK,articles)
}

func GetArticleById(ctx *gin.Context)  {
	//！通过Param取出url里的参数
	id := ctx.Param("id")
	var article models.Article
	//！健壮性：错误类型准确检查
	if err := global.Db.Where("id = ?",id).First(&article).Error; err != nil{
		if errors.Is(err,gorm.ErrRecordNotFound){
			//找不到id为?的数据，返回准确err类型
			ctx.JSON(http.StatusNotFound,gin.H{"error":err.Error()})
		}else{
			ctx.JSON(http.StatusInternalServerError,gin.H{"error":err.Error()})
		}
		return
	}

	ctx.JSON(http.StatusOK,article)

}