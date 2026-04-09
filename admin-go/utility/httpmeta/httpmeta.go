package httpmeta

import (
	"strconv"
	"strings"
	"time"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"

	"gbaseadmin/utility/snowflake"
)

const (
	RequestIDHeader = "X-Request-ID"
	requestIDKey    = "request_id"
)

func RequestID(r *ghttp.Request) string {
	if r == nil {
		return ""
	}
	if value := strings.TrimSpace(r.Response.Header().Get(RequestIDHeader)); value != "" {
		return value
	}
	if value := strings.TrimSpace(r.GetCtxVar(requestIDKey).String()); value != "" {
		return value
	}
	return ""
}

func RequestIDMiddleware(r *ghttp.Request) {
	requestID := strings.TrimSpace(r.GetHeader(RequestIDHeader))
	if requestID == "" {
		requestID = strconv.FormatInt(int64(snowflake.Generate()), 10)
	}
	r.Response.Header().Set(RequestIDHeader, requestID)
	r.SetCtxVar(requestIDKey, requestID)
	r.Middleware.Next()
}

func AccessLogMiddleware(r *ghttp.Request) {
	startedAt := time.Now()
	r.Middleware.Next()

	requestID := RequestID(r)
	userID := strings.TrimSpace(r.GetCtxVar("jwt_user_id").String())
	if userID == "" {
		userID = "anonymous"
	}
	g.Log().Infof(
		r.Context(),
		"[access] request_id=%s method=%s path=%s status=%d duration_ms=%d user_id=%s",
		requestID,
		r.Method,
		r.URL.Path,
		r.Response.Status,
		time.Since(startedAt).Milliseconds(),
		userID,
	)
}
