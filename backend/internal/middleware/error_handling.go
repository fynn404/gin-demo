package middleware

import (
	"github.com/fynn404/gin-demo/backend/internal/common/config"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

// ErrorHandlingMiddleware 捕获 panic 并返回标准化的错误响应
func ErrorHandlingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if r := recover(); r != nil {
				// 记录详细的错误日志
				log.Printf("Panic recovered: %v", r)

				// 检查环境变量以确定是否在开发环境
				if config.GlobalConfig.Server.Mode == "release" {
					// 在生产环境中返回通用错误信息
					c.JSON(http.StatusInternalServerError, gin.H{
						"code":    http.StatusInternalServerError,
						"message": "internal server error",
					})
				} else {
					// 在开发环境中返回详细错误信息
					c.JSON(http.StatusInternalServerError, gin.H{
						"code":    http.StatusInternalServerError,
						"message": "internal server error",
						"error":   r,
					})
				}

				// 中止请求
				c.Abort()
			}
		}()

		// 继续处理请求
		c.Next()
	}
}
