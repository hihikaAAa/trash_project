// Package middlewares
package middlewares

import (
	"fmt"
	"strings"
	"time"

	"github.com/hihikaAAa/trash_project/internal/metrics"

	"github.com/gin-gonic/gin"
)

var MetricsSkipPaths = []string{"/metrics", "/healthz", "/swagger/"}

func CORS() gin.HandlerFunc {
	return func(c *gin.Context) {

		c.Writer.Header().Set("Access-Control-Allow-Origin", c.GetHeader("Origin"))
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Host, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	}
}

func ShouldSkipMetrics(path string) bool {
	for _, p := range MetricsSkipPaths {
		if strings.HasSuffix(p, "/") {
			if strings.HasPrefix(path, p) {
				return true
			}
			continue
		}
		if path == p {
			return true
		}
	}
	return false
}

func APILatencyMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if ShouldSkipMetrics(c.Request.URL.Path) {
			c.Next()
			return
		}
		start := time.Now()
		c.Next()

		route := c.FullPath()
		if route == "" {
			route = "unknown"
		}
		statusClass := fmt.Sprintf("%dxx", c.Writer.Status()/100)
		metrics.LatencyAPI.WithLabelValues(c.Request.Method, route, statusClass).Observe(time.Since(start).Seconds())
	}
}

func NoMethodHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(405, gin.H{"errorCode": "not allowed", "errorMessage": "method not allowed"})
	}
}

func NoRouteHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(404, gin.H{"errorCode": "not found", "errorMessage": "method not found"})
	}
}
