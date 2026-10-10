package configuration

import (
	"net/http"
)

func RegisterRoutes(mux *http.ServeMux, handler *Handler, authMiddleware func(http.Handler) http.Handler) {
	// Only accessible by Super Admins conceptually
	mux.Handle("GET /api/v1/configuration/company", authMiddleware(http.HandlerFunc(handler.GetCompanySettings)))
	mux.Handle("PATCH /api/v1/configuration/company", authMiddleware(http.HandlerFunc(handler.UpdateCompanySettings)))
}
