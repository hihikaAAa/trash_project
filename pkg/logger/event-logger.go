//nolint:contextcheck
package logger

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
	"unicode/utf8"

	"github.com/gin-gonic/gin"
)

const eventLinkIDKeyGin = "event_link_id"

type eventLinkCtxKey string

const eventLinkIDKeyCtx eventLinkCtxKey = "event_link_id"

var (
	eventMu   sync.Mutex
	eventFile *os.File
	eventPath string
)

func InitEventLogFile(path string) error {
	path = strings.TrimSpace(path)
	if path == "" {
		return nil
	}

	eventMu.Lock()
	defer eventMu.Unlock()

	if eventFile != nil && eventPath == path {
		return nil
	}
	if eventFile != nil {
		_ = eventFile.Close()
		eventFile = nil
		eventPath = ""
	}

	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return fmt.Errorf("cannot create event directory: %w", err)
	}

	f, err := os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o644)
	if err != nil {
		return fmt.Errorf("cannot open event file: %w", err)
	}
	eventFile = f
	eventPath = path
	return nil
}

func CloseEventLogFile() error {
	eventMu.Lock()
	defer eventMu.Unlock()

	if eventFile == nil {
		return nil
	}
	err := eventFile.Close()
	eventFile = nil
	eventPath = ""
	return err
}

type bodyWriter struct {
	gin.ResponseWriter
	buf bytes.Buffer
}

func (w *bodyWriter) Write(b []byte) (int, error) {
	_, _ = w.buf.Write(b)
	return w.ResponseWriter.Write(b)
}

func (w *bodyWriter) WriteString(s string) (int, error) {
	_, _ = w.buf.WriteString(s)
	return w.ResponseWriter.WriteString(s)
}

func EventMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if !strings.HasPrefix(c.Request.URL.Path, "/trash/api/v1") && !strings.HasPrefix(c.Request.URL.Path, "/paper/api/v1/") {
			c.Next()
			return
		}
		reqBody, _ := io.ReadAll(c.Request.Body)
		c.Request.Body = io.NopCloser(bytes.NewBuffer(reqBody))
		linkID := extractLinkIDFromRequest(c, reqBody)
		if linkID != "" {
			SetEventLinkID(c, linkID)
		}
		writeRequestBlock(
			"REST",
			valueOrDash(linkID),
			c.Request.Method,
			fullURL(c.Request),
			prettyBody(c.ContentType(), reqBody),
		)
		bw := &bodyWriter{ResponseWriter: c.Writer}
		c.Writer = bw

		c.Next()

		if lid := GetEventLinkID(c); lid != "" {
			linkID = lid
		}
		if linkID == "" {
			linkID = extractLinkIDFromResponse(bw.buf.Bytes())
		}
		writeResponseBlock(
			"REST",
			valueOrDash(linkID),
			prettyBody(c.Writer.Header().Get("Content-Type"), bw.buf.Bytes()))
	}
}

func LogCRMRequest(ctx context.Context, method, reqURL string, body any) {
	writeRequestBlock(
		"CRM",
		valueOrDash(EventLinkID(ctx)),
		method,
		reqURL,
		prettyBodyFromAny(body),
	)
}

func LogCRMResponse(ctx context.Context, body []byte, contentType string) {
	writeResponseBlock(
		"CRM",
		valueOrDash(EventLinkID(ctx)),
		prettyBody(contentType, body),
	)
}

func WithEventLinkID(ctx context.Context, linkID string) context.Context {
	linkID = strings.TrimSpace(linkID)
	if linkID == "" {
		return ctx
	}
	return context.WithValue(ctx, eventLinkIDKeyCtx, linkID)
}

func EventLinkID(ctx context.Context) string {
	if ctx == nil {
		return ""
	}
	v := ctx.Value(eventLinkIDKeyCtx)
	if s, ok := v.(string); ok {
		return strings.TrimSpace(s)
	}
	return ""
}

func SetEventLinkID(c *gin.Context, linkID string) {
	linkID = strings.TrimSpace(linkID)
	if c == nil || linkID == "" {
		return
	}
	c.Set(eventLinkIDKeyGin, linkID)
	if c.Request != nil {
		c.Request = c.Request.WithContext(WithEventLinkID(c.Request.Context(), linkID))
	}
}

func GetEventLinkID(c *gin.Context) string {
	if c == nil {
		return ""
	}
	if v, ok := c.Get(eventLinkIDKeyGin); ok {
		if s, ok2 := v.(string); ok2 {
			return strings.TrimSpace(s)
		}
	}
	return ""
}

func fullURL(req *http.Request) string {
	if req == nil || req.URL == nil {
		return ""
	}
	scheme := "https"
	return fmt.Sprintf("%s://%s%s", scheme, req.Host, req.URL.RequestURI())
}

func extractLinkIDFromRequest(c *gin.Context, body []byte) string {
	if c != nil {
		if v := strings.TrimSpace(c.Query("linkId")); v != "" {
			return v
		}
	}
	if len(bytes.TrimSpace(body)) == 0 || !json.Valid(body) {
		return ""
	}
	var m map[string]any
	if err := json.Unmarshal(body, &m); err != nil {
		return ""
	}
	if v, ok := m["linkId"]; ok {
		if s, ok2 := v.(string); ok2 {
			return strings.TrimSpace(s)
		}
	}
	return ""
}

func extractLinkIDFromResponse(body []byte) string {
	if len(bytes.TrimSpace(body)) == 0 || !json.Valid(body) {
		return ""
	}
	var m map[string]any
	if err := json.Unmarshal(body, &m); err != nil {
		return ""
	}
	if raw, ok := m["linkId"]; ok {
		if s, ok2 := raw.(string); ok2 && strings.TrimSpace(s) != "" {
			return strings.TrimSpace(s)
		}
	}
	if raw, ok := m["URL"]; ok {
		if s, ok2 := raw.(string); ok2 && s != "" {
			s = strings.TrimRight(s, "/")
			parts := strings.Split(s, "/")
			if len(parts) > 0 {
				return strings.TrimSpace(parts[len(parts)-1])
			}
		}
	}
	return ""
}

func writeRequestBlock(kind, linkID, method, reqURL, pretty string) {
	line := fmt.Sprintf("%s %s %s\n\n%s %s\n%s\n", tsNow(), kind, linkID, method, reqURL, pretty)
	writeRaw(line)
}

func writeResponseBlock(kind, linkID, pretty string) {
	line := fmt.Sprintf("%s %s %s\n\n%s\n", tsNow(), kind, linkID, pretty)
	writeRaw(line)
}

func tsNow() string {
	return time.Now().Format("02.01.2006 15:04:05:000")
}

func writeRaw(line string) {
	eventMu.Lock()
	defer eventMu.Unlock()

	if eventFile == nil {
		return
	}
	_, _ = eventFile.WriteString(line)
}

func prettyBodyFromAny(v any) string {
	if v == nil {
		return "<empty>"
	}
	switch t := v.(type) {
	case []byte:
		return prettyBody("", t)
	case string:
		return prettyBody("", []byte(t))
	default:
		b, err := json.MarshalIndent(v, "", "  ")
		if err != nil {
			return fmt.Sprintf("%v", v)
		}
		return string(b)
	}
}

func prettyBody(contentType string, body []byte) string {
	body = bytes.TrimSpace(body)
	if len(body) == 0 {
		return "<empty>"
	}

	if strings.Contains(strings.ToLower(contentType), "json") || json.Valid(body) {
		var out bytes.Buffer
		if err := json.Indent(&out, body, "", "  "); err != nil {
			return string(body)
		}
		return out.String()
	}

	if !utf8.Valid(body) {
		return fmt.Sprintf("<binary %d bytes>", len(body))
	}
	return string(body)
}

func valueOrDash(v string) string {
	v = strings.TrimSpace(v)
	if v == "" {
		return "-"
	}
	return v
}
