package service

import (
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

type ReportService struct {
}

var (
	levelAndMeasure = map[string]struct {
		HazardLevel string
		Measure     string
		Weight      int64
	}{
		"SSH-Bruteforce":        {"高", "立即禁止IP地址", 8},
		"DDoS attack-HOIC":      {"非常高", "增加带宽、使用防火墙", 10},
		"Benign":                {"无", "无需采取措施", 0},
		"Bot":                   {"高", "阻止流量并清除病毒", 6},
		"Dos attacks-GoldenEye": {"非常高", "使用防火墙、增加带宽", 9},
		"Infiltration":          {"非常高", "立即禁止IP地址，并进行深入的检查", 7},
	}
)

func computeGrade(severity int64) string {
	if severity > 200 {
		return "极端"
	} else if severity > 100 {
		return "非常严重"
	} else if severity > 20 {
		return "严重"
	} else if severity > 1 {
		return "一般"
	} else {
		return "无"
	}
}

func (s *ReportService) DetailedReport() *response.Response {
	oneWeeksAgo := time.Now().AddDate(0, 0, -7)
	var oneWeeksAgoRecords []model.Record

	err := database.Mysql.Table("record").Where("date > ?", oneWeeksAgo).
		Find(&oneWeeksAgoRecords).Error
	if err != nil {
		return response.FailWithMessage("数据查询失败")
	}

	// 攻击分类计数
	classificationCount := make(map[string]int64)

	for _, record := range oneWeeksAgoRecords {
		classificationCount[record.Label] += 1
	}

	totalCount := len(oneWeeksAgoRecords)
	var items []dto.Table
	for k, v := range classificationCount {
		frequency := v / 7
		severity := frequency * levelAndMeasure[k].Weight
		grade := computeGrade(severity)

		items = append(items, dto.Table{
			AttackType:  k,
			Proportion:  float64(v) / float64(totalCount),
			HazardLevel: levelAndMeasure[k].HazardLevel,
			Frequency:   frequency,
			Severity:    severity,
			Grade:       grade,
			Measure:     levelAndMeasure[k].Measure,
		})
	}

	return response.OkWithData(items)
}

func (s *TableService) GetTableData() *response.Response {
	// 从数据库中查询最近两周的数据
	oneWeeksAgo := time.Now().AddDate(0, 0, -7)
	twoWeeksAgo := time.Now().AddDate(0, 0, -14)
	var oneWeeksAgoRecords, twoWeeksAgoRecords []model.Record

	e1 := database.Mysql.Table("record").Where("date BETWEEN ? AND ?", twoWeeksAgo, oneWeeksAgo).
		Find(&twoWeeksAgoRecords).Error
	e2 := database.Mysql.Table("record").Where("date > ?", oneWeeksAgo).
		Find(&oneWeeksAgoRecords).Error
	if e1 != nil || e2 != nil {
		return response.FailWithMessage("数据查询失败")
	}

	// 近两周数据总数
	twoWeeksAgoTotal, oneWeeksAgoTotal := len(twoWeeksAgoRecords), len(oneWeeksAgoRecords)

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

	for _, record := range oneWeeksAgoRecords {
		if record.Label == "normal" {
			oneWeeksAgoNormal += 1
			continue
		}

		classificationCount[record.Label] += 1

		date := record.Date
		if date.After(time.Now().AddDate(0, 0, -1)) {
			dayCount[6] += 1
		} else if date.After(time.Now().AddDate(0, 0, -2)) {
			dayCount[5] += 1
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
