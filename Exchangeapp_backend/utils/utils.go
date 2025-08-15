package utils

import (
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
	signedToken,err := token.SignedString("secret")
	return "Bearer " + signedToken,err
}
//检查密码加密后是否和数据库hash一致，返回bool
func CheckPassword(password string,hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash),[]byte(password))
	return err == nil
}