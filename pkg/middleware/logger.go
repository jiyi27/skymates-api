package middleware

import (
	"log"
	"net/http"
	"time"
)

// Logger 是一个记录HTTP请求的中间件
func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// 创建一个自定义的响应写入器来捕获状态码
		rw := &responseWriter{
			ResponseWriter: w,
			statusCode:     http.StatusOK, // 默认状态码
		}

		// 调用下一个处理器
		next.ServeHTTP(rw, r)

		// 计算请求处理时间
		duration := time.Since(start)

		// 记录请求信息
		log.Printf(
			"[%s] %s %s %d %s",
			r.Method,
			r.RequestURI,
			r.RemoteAddr,
			rw.statusCode,
			duration,
		)
	})
}

// 自定义ResponseWriter以捕获状态码
type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}
