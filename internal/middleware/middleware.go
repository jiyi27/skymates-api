package middleware

import (
	"context"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"log"
	"net/http"
	"skymates-api/internal/auth"
	"strings"
	"time"
)

// Middleware 定义中间件函数类型
type Middleware func(http.HandlerFunc) http.HandlerFunc

// CORSConfig 定义 CORS 中间件的配置选项
type CORSConfig struct {
	AllowOrigins     []string // 允许的源域名列表
	AllowMethods     []string // 允许的 HTTP 方法
	AllowHeaders     []string // 允许的请求头
	ExposeHeaders    []string // 暴露给客户端的响应头
	AllowCredentials bool     // 是否允许携带认证信息
	MaxAge           int      // 预检请求的有效期(秒)
}

// Use 按照从前到后的顺序执行提供的 middlewares, 最后执行 handler
func Use(handler http.HandlerFunc, middlewares ...Middleware) http.HandlerFunc {
	// 每次循环都把 handler 包装一层, 第一个中间件被包装在最外层, 所以最先执行
	// 想象中间件是盒子, handler 是礼物, 我们使用多层盒子包装礼物, 用来最后包装礼物的盒子最先打开(执行), 也是最后返回函数
	for i := len(middlewares) - 1; i >= 0; i-- {
		handler = middlewares[i](handler)
	}
	return handler
}

func CORS(config *CORSConfig) Middleware {

	if config == nil {
		config = &CORSConfig{
			// 源是协议(http/https)+域名+端口, 不是用户的 IP
			AllowOrigins: []string{"http://localhost:3000", "http://127.0.0.1:3000", "http://192.168.2.46:3000"},
			AllowMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
			// 不要使用 Access-Control-Allow-Headers: *
			// 看起来是允许所有 header，但在某些浏览器中，
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

	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			log.Print("CORS middleware - start")
			log.Printf("Origin: %s", r.Header.Get("Origin"))

			// If Access-Control-Allow-Credentials = true
			// then Access-Control-Allow-Origin must not use *. Even you set it to *, it will have error.
			// This is a security requirement: https://stackoverflow.com/a/19744754/16317008
			if config.AllowCredentials {
				if origin := r.Header.Get("Origin"); origin != "" {
					for _, allowed := range config.AllowOrigins {
						if allowed == origin {
							w.Header().Set("Access-Control-Allow-Origin", origin)
							break
						}
					}
				}
			} else {
				// no credentials allowed, if AllowOrigins is ["*"], then we can set it to "*"
				if len(config.AllowOrigins) == 1 && config.AllowOrigins[0] == "*" {
					w.Header().Set("Access-Control-Allow-Origin", "*")
				} else if origin := r.Header.Get("Origin"); origin != "" {
					for _, allowed := range config.AllowOrigins {
						if allowed == origin {
							w.Header().Set("Access-Control-Allow-Origin", origin)
							break
						}
					}
				}
			}

			w.Header().Set("Access-Control-Allow-Methods", strings.Join(config.AllowMethods, ","))
			w.Header().Set("Access-Control-Allow-Headers", strings.Join(config.AllowHeaders, ","))

			if len(config.ExposeHeaders) > 0 {
				w.Header().Set("Access-Control-Expose-Headers", strings.Join(config.ExposeHeaders, ","))
			}

			if config.AllowCredentials {
				w.Header().Set("Access-Control-Allow-Credentials", "true")
			}

			// 处理预检请求, 浏览器在发送某些 CORS 请求之前, 会先发送一个 OPTIONS 请求, 用来检查服务器是否允许实际的请求
			// 这种情况通常发生在: 使用非简单方法 (除了 GET、HEAD、POST 之外的方法) 或者自定义请求头时
			if r.Method == "OPTIONS" {
				// 1. 设置缓存时间,告诉浏览器可以缓存预检请求的结果多长时间, 意味着在这段时间内浏览器不需要再发送预检请求
				w.Header().Set("Access-Control-Max-Age", fmt.Sprint(config.MaxAge))
				// 2. 返回 204 状态码（无内容）
				w.WriteHeader(http.StatusNoContent)
				// 3. 结束请求
				return
			}

			next(w, r)

			//log.Print("CORS middleware - end")
		}
	}
}

func Logger(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		//log.Print("Logger middleware - start")

		start := time.Now()
		// 在请求结束后打印日, 这里使用 defer 的主要优势在于:
		// 异常处理: 如果 next(w, r) 执行过程中发生 panic，defer 依然会执行，这样我们能记录到这个请求的日志
		defer func() {
			log.Printf("[%s] %s %v", r.Method, r.URL.Path, time.Since(start))
		}()

		next(w, r)

		//log.Print("Logger middleware - end")
	}
}

type Claims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}

func Auth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// log.Print("Auth middleware - start")
		authHeaderValue := r.Header.Get("Authorization")
		if authHeaderValue == "" {
			http.Error(w, "Authorization header is required", http.StatusUnauthorized)
			return
		}

		bearerToken := strings.Split(authHeaderValue, " ")
		if len(bearerToken) != 2 || bearerToken[0] != "Bearer" {
			http.Error(w, "Invalid authorization format", http.StatusUnauthorized)
			return
		}

		claims, err := auth.ValidateJwtToken(bearerToken[1])
		if err != nil {
			switch {
			case strings.Contains(err.Error(), "signature"):
				http.Error(w, err.Error(), http.StatusUnauthorized)
			case strings.Contains(err.Error(), "expired"):
				http.Error(w, err.Error(), http.StatusUnauthorized)
			default:
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
			}
			return
		}

		ctx := context.WithValue(r.Context(), "username", claims.Username)
		next(w, r.WithContext(ctx))

		// log.Print("Auth middleware - end")
	}
}
