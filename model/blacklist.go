package model

import "time"

type Blacklist struct {
	Id         uint      `json:"id" gorm:"primary_key;auto_increment"`
	Ip         string    `json:"ip"`
	Grade      uint      `json:"grade"`
	Status     string    `json:"status"`
	CreateTime time.Time `json:"create_time"`
	UpdateTime time.Time `json:"update_time"`
}
