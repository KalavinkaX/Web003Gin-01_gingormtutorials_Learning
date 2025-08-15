package controllers

import (
	"exchangeapp/global"
	"exchangeapp/models"
	"exchangeapp/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Register(ctx *gin.Context) {//函数参数必须接收一个*gin.Context 类型的参数
	var user models.User
	if err := ctx.ShouldBind(&user); err != nil {
		ctx.JSON(http.StatusBadRequest,gin.H{
			"error" : err.Error(),
		})
		return
	}
	//密码加密
	hashedPwd,err := utils.HashPassword(user.Password)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError,gin.H{"error" : err.Error(),})
		return
	}
	user.Password = hashedPwd
	
	//生成JWT令牌
	signedToken,err := utils.GenerateJWT(user.Username)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError,gin.H{"error" : err.Error(),})
		return
	}

	//数据库自动创建表(传参 对象实体)
	if err := global.Db.AutoMigrate(&user);err != nil{
		ctx.JSON(http.StatusInternalServerError,gin.H{"error" : err.Error()})
		return
	}
	//创建表记录
	if err := global.Db.Create(&user).Error;err !=nil{
		ctx.JSON(http.StatusInternalServerError,gin.H{"error":err.Error()})
		return
	}

	//前面无error，最终成功注册
	ctx.JSON(http.StatusOK,gin.H{"token":signedToken})
}

func Login(ctx *gin.Context)  {
	var input struct{//新建专用登录操作，只包含username和password的结构体实体变量，不复用models的user结构体
		Username string `json:"username"`//后面用于识别前端传来的小写username到结构体实体
		Password string `json:"password"`
	}
	//绑定只需要username和password的User结构体
	if err := ctx.ShouldBind(&input); err != nil {
		ctx.JSON(http.StatusInternalServerError,err.Error())
		return
	}
	//查询数据库是否有当前username的数据
	var user models.User
	if err := global.Db.Where("username=?",input.Username).First(&user).Error;err !=nil {
		ctx.JSON(http.StatusUnauthorized,gin.H{"error":"can't find username"})
		return
	}
	//真正目标：检查密码是否正确
	if !utils.CheckPassword(input.Password,user.Password){
		ctx.JSON(http.StatusUnauthorized,gin.H{"error":"wrong password!"})
		return
	}

	//生成JWT令牌
	signedToken,err := utils.GenerateJWT(user.Username)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError,gin.H{"error" : err.Error(),})
		return
	}
	//前面无error，最终成功注册
	ctx.JSON(http.StatusOK,gin.H{"token":signedToken})
}