package service

import (
	"fmt"
	"monitor/database"
	"monitor/model"
	"monitor/model/dto"
	"monitor/pkg/response"
	"time"
)

type RecordService struct {
	model.Page
}

type TableService struct {
}

func (s *TableService) GetTableData() *response.Response {
	// 从数据库中查询最近两周的数据
	oneWeeksAgo := time.Now().AddDate(0, 0, -7)
	twoWeeksAgo := time.Now().AddDate(0, 0, -14)
	var onwWeeksAgoRecords, twoWeeksAgoRecords []model.Record

	e1 := database.Mysql.Table("record").Where("date BETWEEN ? AND ?", twoWeeksAgo, oneWeeksAgo).
		Find(&twoWeeksAgoRecords).Error
	e2 := database.Mysql.Table("record").Where("date > ?", oneWeeksAgo).
		Find(&onwWeeksAgoRecords).Error
	if e1 != nil || e2 != nil {
		return response.FailWithMessage("数据查询失败")
	}

	// 近两周数据总数
	twoWeeksAgoTotal, oneWeeksAgoTotal := len(twoWeeksAgoRecords), len(onwWeeksAgoRecords)
	fmt.Println(twoWeeksAgoTotal, oneWeeksAgoTotal, int64(oneWeeksAgoTotal-twoWeeksAgoTotal))

	twoWeeksAgoNormal, oneWeeksAgoNormal := 0, 0
	for _, record := range twoWeeksAgoRecords {
		if record.Label == "normal" {
			twoWeeksAgoNormal += 1
		}
	}

	// 攻击分类计数
	classificationCount := make(map[string]int64)
	// 近七天分天计数
	dayCount := [7]int64{0, 0, 0, 0, 0, 0, 0}

	for _, record := range onwWeeksAgoRecords {
		if record.Label == "normal" {
			oneWeeksAgoNormal += 1
			continue
		}

		classificationCount[record.Label] += 1

		date := record.Date
		if date.After(time.Now().AddDate(0, 0, -1)) {
			dayCount[6] += 1
		} else if date.After(time.Now().AddDate(0, 0, -2)) {
			dayCount[5] += 5
		} else if date.After(time.Now().AddDate(0, 0, -3)) {
			dayCount[4] += 1
		} else if date.After(time.Now().AddDate(0, 0, -4)) {
			dayCount[3] += 1
		} else if date.After(time.Now().AddDate(0, 0, -5)) {
			dayCount[2] += 1
		} else if date.After(time.Now().AddDate(0, 0, -6)) {
			dayCount[1] += 1
		} else if date.After(time.Now().AddDate(0, 0, -7)) {
			dayCount[0] += 1
		}
	}

	// 攻击事件概览
	overview := dto.Overview{
		AttackIncrement: int64((oneWeeksAgoTotal - oneWeeksAgoNormal) - (twoWeeksAgoTotal - twoWeeksAgoNormal)),
		AttackTotal:     int64(oneWeeksAgoTotal - oneWeeksAgoNormal),

		NormalIncrement: int64(oneWeeksAgoNormal - twoWeeksAgoNormal),
		NormalTotal:     int64(oneWeeksAgoNormal),

		Variation: int64(oneWeeksAgoTotal - twoWeeksAgoTotal),
		Total:     int64(oneWeeksAgoTotal),
	}

	// 攻击事件类型统计
	var attackName []string
	var attackNum []int64

	for name, count := range classificationCount {
		attackName = append(attackName, name)
		attackNum = append(attackNum, count)
	}

	statistics := dto.Statistics{
		AttackName: attackName,
		AttackNum:  attackNum,
	}

	// 近七天攻击事件
	analysis := dto.Analysis{
		AttackNum: dayCount,
	}

	return response.OkWithData(map[string]interface{}{
		"overview":   overview,
		"statistics": statistics,
		"analysis":   analysis,
	})
}

func (s *RecordService) QueryRecords() *response.Response {
	pageNum := s.PageNum
	pageSize := s.PageSize
	offset := (pageNum - 1) * pageSize

	var records []model.Record
	var total int64
	e1 := database.Mysql.Table("record").Count(&total).Error
	e2 := database.Mysql.Table("record").Limit(pageSize).Offset(offset).Find(&records).Error
	if e1 != nil || e2 != nil {
		return response.FailWithMessage("历史记录查询失败")
	}

	return response.OkWithData(model.PageResult{
		Data:  records,
		Total: total,
	})
}
