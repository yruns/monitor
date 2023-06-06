package v1

import (
	"github.com/gin-gonic/gin"
	"monitor/model"
	"monitor/pkg/response"
	"monitor/service"
	"net/http"
	"strconv"
)

func QueryRecords(c *gin.Context) {

	pageNum, e1 := strconv.Atoi(c.Query("pageNum"))
	pageSize, e2 := strconv.Atoi(c.Query("pageSize"))
	label := c.Query("label")

	if e1 != nil || e2 != nil {
		response.FailWithStatusCode(http.StatusBadRequest, "参数错误", c)
		return
	}

	var recordService = service.RecordService{
		Page: model.Page{
			PageNum:  pageNum,
			PageSize: pageSize,
		},
	}

	res := recordService.QueryRecords(label)
	response.Result(res, c)
}

func GetTableData(c *gin.Context) {
	var tableService service.TableService
	res := tableService.GetTableData()

	response.Result(res, c)
}

func DetailedReport(c *gin.Context) {
	var reportService service.ReportService
	res := reportService.DetailedReport()

	response.Result(res, c)
}

func IPStatistics(c *gin.Context) {
	var ipService service.RecordService
	res := ipService.IPStatistics()

	response.Result(res, c)
}
