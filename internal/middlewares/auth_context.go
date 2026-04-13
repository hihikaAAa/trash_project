package middlewares

import (
	"strings"

	"github.com/gin-gonic/gin"
)

// Подстройте названия cookie под ваш auth-сервис.
const (
	cookieUserID = "user_id"
	cookieRole   = "role"
)

func InjectActorFromCookies() gin.HandlerFunc {
	return func(c *gin.Context) {
		if v, err := c.Cookie(cookieUserID); err == nil && strings.TrimSpace(v) != "" {
			c.Set("user_id", strings.TrimSpace(v))
		}

		if v, err := c.Cookie(cookieRole); err == nil && strings.TrimSpace(v) != "" {
			c.Set("role", strings.ToLower(strings.TrimSpace(v)))
		}

		c.Next()
	}
}
