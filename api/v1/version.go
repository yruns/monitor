package v1

import (
	"github.com/gin-gonic/gin"
	"monitor/pkg/response"
	"monitor/service"
)

var (
	versionService = service.VersionService{}
)

func GetVersionList(c *gin.Context) {
	//var versionService service.VersionService

	res := versionService.GetVersionList()
	response.Result(res, c)
}

func SetModelVersion(c *gin.Context) {
	//var versionService service.VersionService
	version := c.Query("version")

	res := versionService.SetModelVersion(version)
	response.Result(res, c)
}

func SetVersionList(c *gin.Context) {
	//var versionService service.VersionService

	res := versionService.GetModelVersion()
	response.Result(res, c)
}
