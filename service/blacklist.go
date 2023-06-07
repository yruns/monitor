package service

import (
	"monitor/database"
	"monitor/model"
	"monitor/pkg/response"
	"time"
)

type BlacklistService struct {
	Ip string `json:"ip" form:"ip"`
}

func (s *BlacklistService) GetBlackList(pageNum int, pageSize int) *response.Response {
	var blacklist []model.Blacklist

	offset := (pageNum - 1) * pageSize
	var total int64
	database.Mysql.Table("blacklist").Count(&total)
	err := database.Mysql.Table("blacklist").Offset(offset).Limit(pageSize).Find(&blacklist).Error
	if err != nil {
		return response.FailWithMessage("获取日志信息失败")
	}
	return response.OkWithData(model.PageResult{
		Data:  blacklist,
		Total: total,
	})
}

func (s *BlacklistService) AddBlackList() *response.Response {
	black := model.Blacklist{
		Ip:         s.Ip,
		Grade:      3,
		Status:     "禁止访问",
		CreateTime: time.Now(),
	}
	err := database.Mysql.Table("blacklist").Create(&black).Error
	if err != nil {
		return response.FailWithMessage("添加黑名单失败")
	}

	return response.OkWithData(black)
}
