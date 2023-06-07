package v1

import (
	"github.com/gin-gonic/gin"
	"monitor/pkg/response"
	"monitor/service"
)

func GetVersionList(c *gin.Context) {
	var versionService service.VersionService

	res := versionService.GetVersionList()
	response.Result(res, c)
}
