package middleware

import (
	"encoding/json"
	"fastApi/core/global"
	"fastApi/core/logger"
	"fastApi/util"
	"net"
	"net/http"
	"net/http/httputil"
	"os"
	"runtime/debug"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// GinZap returns a gin.HandlerFunc using configs
func GinZap(confSkipPaths []string) gin.HandlerFunc {

	skipPaths := make(map[string]bool, len(confSkipPaths))
	for _, path := range confSkipPaths {
		skipPaths[path] = true
	}

	return func(c *gin.Context) {
		logger := global.Log
		start := time.Now()
		// some evil middlewares modify this values
		path := c.Request.URL.Path
		//		query := c.Request.URL.RawQuery
		params := util.GetParams(c)

		c.Next()

		if _, ok := skipPaths[path]; !ok {
			end := time.Now()
			runtime := end.Sub(start)

			if len(c.Errors) > 0 {
				// Append error field if this is an erroneous request.
				for _, e := range c.Errors.Errors() {
					logger.Error(e)
				}
			} else {
				headers, _ := json.Marshal(c.Request.Header)

				paramsJson, _ := json.Marshal(params)

				fields := []zapcore.Field{
					zap.Int("userId", c.GetInt("userId")),
					zap.Int("status", c.Writer.Status()),
					zap.String("host", c.Request.Host),
					zap.String("url", c.Request.URL.String()),
					zap.String("method", c.Request.Method),
					zap.String("ip", c.ClientIP()),
					zap.String("headers", string(headers)),
					zap.String("params", string(paramsJson)),
					zap.Duration("runtime", runtime),
				}

				logger.Info("", fields...)
			}
		}
	}
}

// RecoveryWithZap returns a gin.HandlerFunc (middleware)
// that recovers from any panics and logs requests using uber-go/zap.
// All errors are logged using zap.Error().
// stack means whether output the stack info.
// The stack info is easy to find where the error occurs but the stack info is too large.
func RecoveryWithZap(stack bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		logger := global.Log
		defer func() {
			if err := recover(); err != nil {
				// Check for a broken connection, as it is not really a
				// condition that warrants a panic stack trace.
				var brokenPipe bool
				if ne, ok := err.(*net.OpError); ok {
					if se, ok := ne.Err.(*os.SyscallError); ok {
						if strings.Contains(strings.ToLower(se.Error()), "broken pipe") || strings.Contains(strings.ToLower(se.Error()), "connection reset by peer") {
							brokenPipe = true
						}
					}
				}

				httpRequest, _ := httputil.DumpRequest(c.Request, false)
				if brokenPipe {
					logger.Error(c.Request.URL.Path,
						zap.Any("error", err),
						zap.String("request", string(httpRequest)),
					)
					// If the connection is dead, we can't write a status to it.
					c.Error(err.(error)) // nolint: errcheck
					c.Abort()
					return
				}

				if stack {
					logger.Error("[Recovery from panic]",
						zap.Time("time", time.Now()),
						zap.Any("error", err),
						zap.String("request", string(httpRequest)),
						zap.String("stack", string(debug.Stack())),
					)
				} else {
					logger.Error("[Recovery from panic]",
						zap.Time("time", time.Now()),
						zap.Any("error", err),
						zap.String("request", string(httpRequest)),
					)
				}
				c.AbortWithStatus(http.StatusInternalServerError)
			}
		}()
		c.Next()
	}
}

func AddTraceId() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// 每个请求生成的请求traceId具有全局唯一性
		traceId, spanId := logger.CalcTraceId(ctx.Request.Context())
		ctx.Set(logger.TraceId, traceId)
		ctx.Set(logger.SpanId, spanId)
		logger.With(
			ctx,
			zap.String("traceId", traceId),
		)
		ctx.Next()
	}
}
