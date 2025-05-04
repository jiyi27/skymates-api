package middleware

import (
	"fmt"
	"net/http"
	"strings"
)

// CORSConfig 定义 CORS 中间件的配置选项
type CORSConfig struct {
	AllowOrigins     []string // 允许的源域名列表
	AllowMethods     []string // 允许的 HTTP 方法
	AllowHeaders     []string // 允许的请求头
	ExposeHeaders    []string // 暴露给客户端的响应头
	AllowCredentials bool     // 是否允许携带认证信息
	MaxAge           int      // 预检请求的有效期(秒)
}

// DefaultCORSConfig 提供默认的CORS配置
func DefaultCORSConfig() *CORSConfig {
	return &CORSConfig{
		// 源是协议(http/https)+域名+端口, 不是用户的 IP
		AllowOrigins: []string{"http://localhost:3000", "http://127.0.0.1:3000", "http://192.168.2.46:3000", "http://172.20.10.14:3000", "https://skymates-shwezhus-projects.vercel.app/"},
		AllowMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		// 不要使用 Access-Control-Allow-Headers: *
		// 当 credentials: 'include' 或 withCredentials: true 时，通配符 * 可能不会被正确解析
		AllowHeaders: []string{
			"Content-Type",
			"Authorization",
			"X-Requested-With",
			"Accept",
			"Origin",
			"Access-Control-Request-Method",
			"Access-Control-Request-Headers",
		},
		ExposeHeaders:    []string{},
		AllowCredentials: true,
		MaxAge:           86400,
	}
}

// CORS 返回一个处理跨域资源共享(CORS)的中间件
func CORS(next http.Handler) http.Handler {
	// 如果没有提供配置，使用默认配置
	config := DefaultCORSConfig()

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//log.Print("CORS middleware - start")
		//log.Printf("Origin: %s", r.Header.Get("Origin"))

		origin := r.Header.Get("Origin")
		// 处理 Origin 头
		if origin != "" {
			// 如果启用凭证，必须指定确切的 Origin（不能使用通配符）
			if config.AllowCredentials {
				for _, allowed := range config.AllowOrigins {
					if allowed == origin {
						w.Header().Set("Access-Control-Allow-Origin", origin)
						break
					}
				}
			} else {
				// 没有凭证时的处理
				if len(config.AllowOrigins) == 1 && config.AllowOrigins[0] == "*" {
					w.Header().Set("Access-Control-Allow-Origin", "*")
				} else {
					for _, allowed := range config.AllowOrigins {
						if allowed == origin {
							w.Header().Set("Access-Control-Allow-Origin", origin)
							break
						}
					}
				}
			}
		}

		// 设置其他CORS头
		w.Header().Set("Access-Control-Allow-Methods", strings.Join(config.AllowMethods, ","))
		w.Header().Set("Access-Control-Allow-Headers", strings.Join(config.AllowHeaders, ","))

		if len(config.ExposeHeaders) > 0 {
			w.Header().Set("Access-Control-Expose-Headers", strings.Join(config.ExposeHeaders, ","))
		}

		if config.AllowCredentials {
			w.Header().Set("Access-Control-Allow-Credentials", "true")
		}

		// 处理预检请求
		if r.Method == http.MethodOptions {
			w.Header().Set("Access-Control-Max-Age", fmt.Sprint(config.MaxAge))
			w.WriteHeader(http.StatusNoContent)
			return
		}

		// 继续处理实际请求
		next.ServeHTTP(w, r)
	})
}
