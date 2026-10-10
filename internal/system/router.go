package system

import (
	"net/http"
)

func RegisterRoutes(mux *http.ServeMux, handler *Handler, authMiddleware func(http.Handler) http.Handler) {
	// Unauthenticated health checks
	mux.HandleFunc("GET /health/live", handler.Liveness)
	mux.HandleFunc("GET /health/ready", handler.Readiness)
	
	// Authenticated or unauthenticated version endpoint depending on preference, we'll leave it unauthenticated for monitoring
	mux.HandleFunc("GET /api/v1/system/version", handler.Version)
}
