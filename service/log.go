package service

import (
	"monitor/database"
	"monitor/model"
	"monitor/pkg/response"
)

type LogService struct {
}

func (s *LogService) GetLogList(pageNum, pageSize int) *response.Response {
	var logs []model.Log

	offset := (pageNum - 1) * pageSize
	var total int64
	database.Mysql.Table("log").Count(&total)
	err := database.Mysql.Table("log").Offset(offset).Limit(pageSize).Find(&logs).Error
	if err != nil {
		return response.FailWithMessage("获取日志信息失败")
	}
	return response.OkWithData(model.PageResult{
		Data:  logs,
		Total: total,
	})
}
