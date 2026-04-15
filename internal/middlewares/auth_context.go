package middlewares

import (
	"encoding/base64"
	"encoding/json"
	"strings"

	"github.com/gin-gonic/gin"
)

const (
	cookieUserID = "user_id"
	cookieRole   = "role"
)

type tokenPayload struct {
	Sub        string `json:"sub"`
	Permission string `json:"permission"`
}

func InjectActorFromCookies() gin.HandlerFunc {
	return func(c *gin.Context) {
		if v, err := c.Cookie(cookieUserID); err == nil && strings.TrimSpace(v) != "" {
			c.Set("user_id", strings.TrimSpace(v))
		}

		if v, err := c.Cookie(cookieRole); err == nil && strings.TrimSpace(v) != "" {
			c.Set("role", strings.ToLower(strings.TrimSpace(v)))
		}

		// Fallback for integration with auth-service JWT Bearer tokens.
		// We intentionally decode the payload to extract actor context so the
		// orders service can work with the existing auth flow.
		if _, hasUser := c.Get("user_id"); !hasUser {
			if token := extractBearerToken(c.GetHeader("Authorization")); token != "" {
				if payload, ok := parseTokenPayload(token); ok {
					if strings.TrimSpace(payload.Sub) != "" {
						c.Set("user_id", strings.TrimSpace(payload.Sub))
					}
					if strings.TrimSpace(payload.Permission) != "" {
						c.Set("role", strings.ToLower(strings.TrimSpace(payload.Permission)))
					}
				}
			}
		}

		c.Next()
	}
}

func extractBearerToken(header string) string {
	parts := strings.Fields(strings.TrimSpace(header))
	if len(parts) == 2 && strings.EqualFold(parts[0], "Bearer") {
		return parts[1]
	}
	return ""
}

func parseTokenPayload(token string) (tokenPayload, bool) {
	parts := strings.Split(strings.TrimSpace(token), ".")
	if len(parts) < 2 {
		return tokenPayload{}, false
	}

	payloadPart := parts[1]
	payloadPart = strings.ReplaceAll(payloadPart, "-", "+")
	payloadPart = strings.ReplaceAll(payloadPart, "_", "/")
	if mod := len(payloadPart) % 4; mod != 0 {
		payloadPart += strings.Repeat("=", 4-mod)
	}

	raw, err := base64.StdEncoding.DecodeString(payloadPart)
	if err != nil {
		return tokenPayload{}, false
	}

	var payload tokenPayload
	if err = json.Unmarshal(raw, &payload); err != nil {
		return tokenPayload{}, false
	}

	return payload, true
}
