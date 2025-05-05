package v1

import (
	"net/http"
	"skymates-api/internal/handler"
	"skymates-api/internal/service"
)

// RegisterTermRoutes 注册V1版本的所有 Term API 路由
func RegisterTermRoutes(mux *http.ServeMux, termService service.TermService) {
	termHandler := handler.NewTermHandler(termService)

	// 公开路由
	mux.HandleFunc("GET /api/v1/terms/search", termHandler.SearchTerms)
	mux.HandleFunc("GET /api/v1/terms/{id}", termHandler.GetTermByID)
	mux.HandleFunc("GET /api/v1/categories/{categoryID}/terms", termHandler.ListTermsByCategory)
	mux.HandleFunc("POST /api/v1/terms", termHandler.CreateTerm)
	mux.HandleFunc("PUT /api/v1/terms/{id}", termHandler.UpdateTerm)
}
