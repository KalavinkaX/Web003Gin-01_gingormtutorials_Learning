package models

import (
	"gorm.io/gorm"
)

type User struct{
	gorm.Model//嵌入gorm预定义的结构体，包含id主键和三个时间
	Username string	`gorm:"unique"`//设定username属性为unique
	Password string
}