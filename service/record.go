package service

import (
	"context"
	"encoding/json"
	"fmt"
	"monitor/database"
	"monitor/model"
	"monitor/model/dto"
	"monitor/pkg/response"
	"monitor/pkg/utils"
	"net"
	"strings"
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
	cache, e := database.Redis.Get(context.Background(), "data:reports").Result()
	if e == nil && cache != "" {
		var items []dto.Table

		e := json.Unmarshal([]byte(cache), &items)
		if e != nil {
			return response.FailWithMessage("json 反序列化失败" + e.Error())
		}
		return response.OkWithData(items)
	}

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

	// 写入redis
	bytes, err := json.Marshal(items)
	if err != nil {
		return response.FailWithMessage("json序列化失败")
	}

	database.Redis.SetEx(context.Background(), "data:reports", string(bytes), time.Hour*24*30)

	return response.OkWithData(items)
}

func (s *TableService) GetTableData() *response.Response {
	// 查询是否在缓存中
	cache, err := database.Redis.Get(context.Background(), "data:record").Result()
	if cache != "" && err == nil {
		var result dto.TotalResult
		err := json.Unmarshal([]byte(cache), &result)
		if err != nil {
			return response.FailWithMessage("json Unmarshal失败")
		}

		return response.OkWithData(result)
	}

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
		if record.Label == "Benign" {
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

	result := dto.TotalResult{
		OverviewResult:   overview,
		StatisticsResult: statistics,
		AnalysisResult:   analysis,
	}

	bytes, err := json.Marshal(result)
	if err != nil {
		return nil
	}

	database.Redis.SetEx(context.Background(), "data:record", string(bytes), time.Second*60*60*24*30)

	return response.OkWithData(result)
}

func (s *RecordService) QueryRecords(label string) *response.Response {
	pageNum := s.PageNum
	pageSize := s.PageSize
	offset := (pageNum - 1) * pageSize

	var records []model.Record
	var total int64
	db := database.Mysql.Table("record")
	if label != "" {
		db = db.Where("label = ?", label)
	}
	e1 := db.Count(&total).Error
	e2 := db.Limit(pageSize).Offset(offset).Find(&records).Error
	if e1 != nil || e2 != nil {
		return response.FailWithMessage("历史记录查询失败")
	}

	return response.OkWithData(model.PageResult{
		Data:  records,
		Total: total,
	})
}

func (s *RecordService) IPStatistics() *response.Response {

	// 从redis中读取
	cache, e := database.Redis.Get(context.Background(), "data:ipStatistics").Result()
	if e == nil && cache != "" {
		var items [10]dto.IpStatistics
		err := json.Unmarshal([]byte(cache), &items)
		if err != nil {
			return response.FailWithMessage("json反序列化失败，" + err.Error())
		}
		return response.OkWithData(items)
	}

	aWeekAgo := time.Now().AddDate(0, 0, -7)

	// 根据ip分组查询
	var ipCount []struct {
		SrcHost string
		//Records []model.Record `gorm:"foreignKey:src_host"`
		Count int64
	}
	err := database.Mysql.Table("record").
		Select("src_host, COUNT(*) as count").
		Where("date > ? AND label <> 'Benign'", aWeekAgo).
		Group("src_host").
		Order("count DESC").
		Limit(10).
		Find(&ipCount).Error
	if err != nil {
		return response.FailWithMessage("数据查询失败," + err.Error())
	}

	var items [10]dto.IpStatistics
	// 根据ip进行统计
	for i, count := range ipCount {
		count.SrcHost = strings.TrimSpace(count.SrcHost)
		ip := net.ParseIP(count.SrcHost)
		record, err := database.GeoIP.City(ip)
		if err != nil {
			return response.FailWithMessage("ip地址查询失败, " + err.Error())
		}

		var records []model.Record
		database.Mysql.Table("record").Where("src_host = ?", count.SrcHost).Find(&records)

		var attackName []string
		for _, r := range records {
			if !utils.Contains(attackName, r.Label) {
				attackName = append(attackName, r.Label)
			}
		}

		items[i] = dto.IpStatistics{
			IP:         count.SrcHost,
			Count:      int64(len(records)),
			Address:    record.Country.Names["zh-CN"] + record.City.Names["zh-CN"],
			AttackName: attackName,
			StartTime:  records[0].Date,
			EndTime:    records[len(records)-1].Date,
		}
	}

	// 写入到redis
	bytes, err := json.Marshal(items)
	if err == nil {
		err = database.Redis.SetEx(context.Background(), "data:ipStatistics", bytes, time.Hour*24*90).Err()
		if err != nil {
			fmt.Println(err.Error())
		}
	}

	return response.OkWithData(items)
}
