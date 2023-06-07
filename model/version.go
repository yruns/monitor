package model

import "time"

type Version struct {
	Id       uint      `json:"id" gorm:"primary_key"`
	Tag      string    `json:"tag"`
	DataPath string    `json:"data_path"`
	Date     time.Time `json:"time"`
}

//type Version struct {
//	Id uint
//}
