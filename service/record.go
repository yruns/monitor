package service

import (
	"monitor/database"
	"monitor/model"
	"monitor/pkg/response"
)

type RecordService struct {
	model.Page
}

func (s *RecordService) QueryRecords() *response.Response {
	pageNum := s.PageNum
	pageSize := s.PageSize
	offset := (pageNum - 1) * pageSize

	var records []model.Record
	err := database.Mysql.Table("record").Limit(pageSize).Offset(offset).Find(&records).Error
	if err != nil {
		return response.FailWithMessage("历史记录查询失败")
	}

	return response.OkWithData(records)
}
