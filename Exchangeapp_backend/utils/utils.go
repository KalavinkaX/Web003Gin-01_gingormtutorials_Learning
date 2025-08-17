package utils

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

//密码加密算法
func HashPassword(pwd string) (string, error) {
	hashedPwd,err := bcrypt.GenerateFromPassword([]byte(pwd),12)
	return string(hashedPwd),err
}

func GenerateJWT(username string) (string,error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256,jwt.MapClaims{
		"username" : username,
		"exp" : time.Now().Add(time.Hour * 72).Unix(),//据1970年，现在时间后3days的时间戳
	})
	//最后加密前两部分
	signedToken,err := token.SignedString([]byte("secret"))//这里最后加密JWT需要[]byte类型的
	return "Bearer " + signedToken,err
}
//检查密码加密后是否和数据库hash一致，返回bool
func CheckPassword(password string,hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash),[]byte(password))
	return err == nil
}

//解析出JWT的用户名
func ParseJWT(tokenString string) (string,error) {
	//这里处理掉请求头的Authorization前缀包含"Bearer"，使得tokenString去掉了Bearer
	if len(tokenString) > 7 && tokenString[:7] == "Bearer " {
		tokenString = tokenString[7:]
	}

	// jwt.Parse()最终返回 1.解析出的所有信息的 *jwt.Token 结构体和2.error
	// 检查JWT是否是HMAC加密方法家族的一员
	token,err := jwt.Parse(tokenString,func (token *jwt.Token)(interface{},error){
		//jwt.Parse()回调方法内 1.先检查签名验证: 确保令牌未被篡改。
		//2.算法检查: 在回调函数中，确保令牌的签名算法是服务器所期望的。
		//3.标准声明验证: 自动检查 exp (过期时间), nbf (生效时间), iat (签发时间) 等标准声明。
		if _,ok := token.Method.(*jwt.SigningMethodHMAC); !ok{
			return nil,errors.New("unexpected Signing Method!")
		}
		return []byte("secret") , nil //jwt.Parse的第二参数函数，返回值(密钥,error)返回到了jwt.Parse()函数内部，而不是直接返回到ParseJWT()!!!
	})
	fmt.Println("first return value = ",token)
	if err != nil{
		return "",err
	}

	//提取载荷(获取JWT中间部分自己储存的信息)
	if claims,ok := token.Claims.(jwt.MapClaims);ok && token.Valid{
		fmt.Println("Claims : ",claims)
		username,ok := claims["username"].(string)
		if !ok {
			return "",errors.New("username in claims is not a String!")
		}
		return username,nil
	}
	return "",err

}