package middleware

import (
	"github.com/gin-gonic/gin"
	"monitor/pkg/response"
	"monitor/pkg/utils"
	"net/http"
	"time"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {

		// 从请求头中获取Token
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" {
			response.FailWithStatusCode(http.StatusUnauthorized, "请先登录", c)
			c.Abort()
			return
		}

		claims, err := utils.ValidToken(tokenString)
		if err != nil {
			response.FailWithStatusCode(http.StatusUnauthorized, "token验证失败", c)
			c.Abort()
			return
		}

		// 验证是否过期
		if claims.ExpiresAt.Before(time.Now()) {
			response.FailWithStatusCode(http.StatusUnauthorized, "token已过期", c)
			c.Abort()
			return
		}

		// 将用户信息写入上下文
		c.Set("username", claims.Username)
		c.Set("userId", claims.Id)

		c.Next()
	}
}
