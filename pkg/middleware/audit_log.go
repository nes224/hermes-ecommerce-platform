package middleware

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func AuditLog() gin.HandlerFunc {
	return func(c *gin.Context) {
		startTime := time.Now()

		traceID := c.GetHeader("X-Trace-ID")
		if traceID == "" {
			traceID = uuid.New().String()
		}
		c.Set("TraceID", traceID)
		c.Header("X-Trace-ID", traceID)

		c.Next()

		latency := time.Since(startTime)
		statusCode := c.Writer.Status()
		clientIP := c.ClientIP()
		method := c.Request.Method
		path := c.Request.URL.Path

		log.Printf("[AUDIT] trace_id=%s ip=%s method=%s path=%s status=%d latency=%s", traceID, clientIP, method, path, statusCode, latency)
	}
}
