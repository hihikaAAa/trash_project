// Package logger
package logger

import (
	"log/slog"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

const requestIDHeader = "X-Request-Id"

type Config struct {
	Level slog.Level
}

func Init(cfg Config) {
	h := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: cfg.Level,
	})
	slog.SetDefault(slog.New(h))
}

func InitFromEnv() {
	lvl := strings.ToLower(strings.TrimSpace(os.Getenv("LOG_LEVEL")))
	if lvl == "" {
		lvl = "info"
	}
	var level slog.Level
	switch lvl {
	case "debug":
		level = slog.LevelDebug
	case "info":
		level = slog.LevelInfo
	case "warn", "warning":
		level = slog.LevelWarn
	case "error":
		level = slog.LevelError
	default:
		level = slog.LevelInfo
	}
	Init(Config{Level: level})
}

func GinMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		c.Next()

		latency := time.Since(start)

		status := c.Writer.Status()
		level := levelByStatus(status)

		attrs := []any{
			"ts", start.UTC().Format(time.RFC3339Nano),
			"status", status,
			"method", c.Request.Method,
			"path", c.FullPath(),
			"raw_path", c.Request.URL.Path,
			"query", c.Request.URL.RawQuery,
			"client_ip", c.ClientIP(),
			"user_agent", c.Request.UserAgent(),
			"latency_ms", latency.Milliseconds(),
			"bytes_out", c.Writer.Size(),
		}

		if rid := requestID(c); rid != "" {
			attrs = append(attrs, "request_id", rid)
		}
		if errMsg := c.Errors.ByType(gin.ErrorTypePrivate).String(); errMsg != "" {
			attrs = append(attrs, "gin_errors", errMsg)
		}

		msg := "http_request"
		switch level {
		case slog.LevelError:
			slog.Error(msg, attrs...)
		case slog.LevelWarn:
			slog.Warn(msg, attrs...)
		case slog.LevelDebug:
			slog.Debug(msg, attrs...)
		default:
			slog.Info(msg, attrs...)
		}
	}
}

func ErrorCtx(c *gin.Context, err error, msg string, extra ...any) {
	if err == nil {
		return
	}
	attrs := []any{
		"err", err.Error(),
		"method", safeMethod(c),
		"path", safePath(c),
		"client_ip", safeIP(c),
	}
	if rid := requestID(c); rid != "" {
		attrs = append(attrs, "request_id", rid)
	}
	if len(extra) > 0 {
		attrs = append(attrs, extra...)
	}
	slog.Error(msg, attrs...)
}

func InfoCtx(c *gin.Context, msg string, extra ...any) {
	attrs := []any{
		"method", safeMethod(c),
		"path", safePath(c),
		"client_ip", safeIP(c),
	}
	if rid := requestID(c); rid != "" {
		attrs = append(attrs, "request_id", rid)
	}
	if len(extra) > 0 {
		attrs = append(attrs, extra...)
	}
	slog.Info(msg, attrs...)
}

func WithRequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		rid := c.GetHeader(requestIDHeader)
		if strings.TrimSpace(rid) == "" {
			rid = strconv.FormatInt(time.Now().UnixNano(), 10)
		}
		c.Header(requestIDHeader, rid)
		c.Set("request_id", rid)
		c.Next()
	}
}

func requestID(c *gin.Context) string {
	if c == nil {
		return ""
	}
	if v, ok := c.Get("request_id"); ok {
		if s, ok := v.(string); ok {
			return s
		}
	}
	rid := strings.TrimSpace(c.GetHeader(requestIDHeader))
	return rid
}

func levelByStatus(code int) slog.Level {
	switch {
	case code >= http.StatusInternalServerError:
		return slog.LevelError
	case code >= http.StatusBadRequest:
		return slog.LevelWarn
	default:
		return slog.LevelInfo
	}
}

func safeMethod(c *gin.Context) string {
	if c == nil || c.Request == nil {
		return ""
	}
	return c.Request.Method
}

func safePath(c *gin.Context) string {
	if c == nil || c.Request == nil {
		return ""
	}
	if p := c.FullPath(); p != "" {
		return p
	}
	return c.Request.URL.Path
}

func safeIP(c *gin.Context) string {
	if c == nil {
		return ""
	}
	return c.ClientIP()
}
