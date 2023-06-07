package v1

import (
	"github.com/gin-gonic/gin"
	"monitor/pkg/response"
	"monitor/service"
	"net/http"
	"strconv"
)

func GetLogList(c *gin.Context) {
	var logService service.LogService

	pageNum, e1 := strconv.Atoi(c.Query("pageNum"))
	pageSize, e2 := strconv.Atoi(c.Query("pageSize"))

	if e1 != nil || e2 != nil {
		response.FailWithStatusCode(http.StatusBadRequest, "参数错误", c)
		return
	}

	res := logService.GetLogList(pageNum, pageSize)
	response.Result(res, c)
}
