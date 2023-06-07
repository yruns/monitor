package routes

import (
	"github.com/gin-gonic/gin"
	api "monitor/api/v1"
	"monitor/middleware"
	"net/http"
)

func NewRouter() *gin.Engine {
	r := gin.Default()

	// 全局中间件
	r.Use(middleware.Cors())
	// 设置静态资源路径
	r.StaticFS("/static", http.Dir("./static"))
	// 设置路由组
	v1 := r.Group("/api/v1")
	{
		v1.GET("/ping", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"message": "pong",
			})
		})

		// 用户操作
		v1.POST("/user/register", api.UserRegister)
		v1.POST("/user/login", api.UserLogin)

		authed := v1.Group("/")                 // 需要登录操作
		authed.Use(middleware.AuthMiddleware()) // jwt鉴权
		authed.PUT("/user", api.UserUpdate)
		authed.POST("/user/upload", api.UploadAvatar)

		authed.POST("/email/code", api.SendEmailCode)
		authed.GET("/email/verify", api.VerifyEmail)
		authed.POST("/email/warning", api.SendWarningEmail)

		authed.GET("/records/list", api.QueryRecords)
		authed.GET("/records/tables", api.GetTableData)
		authed.GET("/records/reports", api.DetailedReport)
		authed.GET("/records/ip", api.IPStatistics)

		authed.GET("/version/list", api.GetVersionList)
	}
	return r
}
