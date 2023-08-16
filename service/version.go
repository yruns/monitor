package service

import (
	"context"
	"monitor/database"
	"monitor/model"
	"monitor/model/dto"
	"monitor/pkg/response"
	"strconv"
	"time"
)

type VersionService struct {
}

func (s *VersionService) GetVersionList() *response.Response {
	var versions []model.Version

	err := database.Mysql.Table("version").Find(&versions).Error
	if err != nil {
		return response.FailWithMessage("模型版本信息获取失败")
	}

	var versionDtos []dto.Version
	for _, version := range versions {
		var aucList []model.Auc
		err := database.Mysql.Table("auc").Where("version_id = ?", version.Id).Find(&aucList).Error
		if err != nil {
			return response.FailWithMessage("获取AUC指标失败")
		}

		versionDtos = append(versionDtos, dto.Version{
			Version: version,
			AucList: aucList,
		})
	}

	return response.OkWithData(versionDtos)
}

func (s *VersionService) SetModelVersion(version string) *response.Response {

	vInt, err := strconv.Atoi(version)
	if err != nil {
		// 获取第一个作为默认值
		var v model.Version
		err := database.Mysql.Table("version").First(&v).Error
		if err != nil {
			return response.FailWithMessage("获取版本信息失败")
		}
		vInt = int(v.Id)
	}

	database.Redis.SetEx(context.Background(), "version:id", strconv.Itoa(vInt), time.Hour*24*30*30)
	return response.OkWithMessage("模型版本设置成功")
}

func (s *VersionService) GetModelVersion() *response.Response {

	vId, err := database.Redis.Get(context.Background(), "version:id").Result()
	if err != nil {
		// 获取第一个作为默认值
		var v model.Version
		err := database.Mysql.Table("version").First(&v).Error
		if err != nil {
			return response.FailWithMessage("获取版本信息失败")
		}
		vId = strconv.Itoa(int(v.Id))
	}

	return response.OkWithData(vId)
}
