package service

import (
	"monitor/database"
	"monitor/model"
	"monitor/model/dto"
	"monitor/pkg/response"
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
