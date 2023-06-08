package v1

import (
	"github.com/gin-gonic/gin"
	"monitor/database"
	"monitor/pkg/response"
	"monitor/service"
	"net/http"
	"strconv"
	"strings"
)

func VerifyEmail(c *gin.Context) {

	userId, _ := c.Get("userId")

	codeFromUser := c.Query("code")
	codeFromDB, _ := database.Redis.Get(c, "user:verify:"+strconv.Itoa(int(userId.(uint)))).Result()

	if strings.EqualFold(codeFromDB, codeFromUser) {
		response.Result(response.Ok(), c)
	} else {
		response.FailWithStatusCode(http.StatusOK, "验证失败", c)
	}
}

func SendEmailCode(c *gin.Context) {
	var emailService service.EmailService

	if err := c.ShouldBind(&emailService); err != nil {
		response.FailWithStatusCode(http.StatusBadRequest, "参数绑定失败", c)
		return
	}

	userId, _ := c.Get("userId")

	res := emailService.SendEmailCode(int64(userId.(uint)))
	response.Result(res, c)
}

func SendWarningEmail(c *gin.Context) {
	var warningEmailService service.WarningEmailService

	if err := c.ShouldBind(&warningEmailService); err != nil {
		response.FailWithStatusCode(http.StatusBadRequest, "参数绑定失败", c)
		return
	}

	userId, _ := c.Get("userId")

	res := warningEmailService.SendWarningEmail(int64(userId.(uint)))
	response.Result(res, c)
}
