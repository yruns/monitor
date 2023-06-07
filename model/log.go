package model

import "time"

type Log struct {
	Id   uint      `json:"id" form:"id" gorm:"primary_key"`
	Info string    `json:"info" form:"info"`
	From string    `json:"from" form:"to"`
	Date time.Time `json:"date" form:"time"`
}
