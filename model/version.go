package model

import "time"

type Version struct {
	Id       uint      `json:"id"`
	Tag      string    `json:"tag"`
	DataPath string    `json:"data_path"`
	Date     time.Time `json:"time"`
}

//type Version struct {
//	Id uint
//}
