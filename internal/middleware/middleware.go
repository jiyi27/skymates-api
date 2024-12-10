package middleware

import (
	"log"
	"net/http"
	"time"
)

// Logger 中间件用于记录请求日志
func Logger(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		defer func() {
			log.Printf("[%s] %s %s %v", r.Method, r.RemoteAddr, r.URL.Path, time.Since(start))
		}()
		next(w, r)
	}
}

// ErrorHandler 统一错误处理中间件
func ErrorHandler(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("panic: %v", err)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			}
		}()
		next(w, r)
	}
}

// Auth 认证中间件
func Auth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")
		if token == "" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		// TODO: 实现token验证逻辑
		next(w, r)
	}
}

// Chain 中间件链式调用
func Chain(handlers ...func(http.HandlerFunc) http.HandlerFunc) func(http.HandlerFunc) http.HandlerFunc {
	return func(final http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			last := final
			for i := len(handlers) - 1; i >= 0; i-- {
				last = handlers[i](last)
			}
			last(w, r)
		}
	}
}
