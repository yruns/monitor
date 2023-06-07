package v1

import (
	"github.com/gin-gonic/gin"
	"monitor/pkg/response"
	"monitor/service"
	"net/http"
	"strconv"
)

func GetBlackList(c *gin.Context) {

	pageNum, e1 := strconv.Atoi(c.Query("pageNum"))
	pageSize, e2 := strconv.Atoi(c.Query("pageSize"))

	if e1 != nil || e2 != nil {
		response.FailWithStatusCode(http.StatusBadRequest, "参数错误", c)
		return
	}

	var blacklistService service.BlacklistService
	res := blacklistService.GetBlackList(pageNum, pageSize)
	response.Result(res, c)
}

func AddBlackList(c *gin.Context) {
	var blacklistService service.BlacklistService
	if err := c.ShouldBind(&blacklistService); err != nil {
		response.FailWithStatusCode(http.StatusBadRequest, "参数错误", c)
		return
	}

	res := blacklistService.AddBlackList()
	response.Result(res, c)
}
