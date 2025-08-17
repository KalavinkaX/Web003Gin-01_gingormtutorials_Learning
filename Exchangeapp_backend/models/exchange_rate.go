package models

import "time"

type ExchangeRate struct {
	ID           uint    `gorm:"primaryKey" json:"_id"`//设置为主键
	FromCurrency string  `json:"fromCurrency" binding:"required"`//设置请求参数必须携带
	ToCurrency   string  `json:"toCurrency" binding:"required"`//设置请求参数必须携带
	Rate         float64 `json:"rate" binding:"required"`//设置请求参数必须携带
	Date         time.Time`json:"date"`
}