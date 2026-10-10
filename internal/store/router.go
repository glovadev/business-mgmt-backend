package store

import (
	"net/http"
)

func RegisterRoutes(mux *http.ServeMux, handler *Handler, authMiddleware func(http.Handler) http.Handler) {
	// Assuming these are all protected by authentication middleware
	
	mux.Handle("POST /api/v1/stores", authMiddleware(http.HandlerFunc(handler.CreateStore)))
	mux.Handle("GET /api/v1/stores", authMiddleware(http.HandlerFunc(handler.ListStores)))
	mux.Handle("GET /api/v1/stores/{store_id}", authMiddleware(http.HandlerFunc(handler.GetStore)))
	mux.Handle("PATCH /api/v1/stores/{store_id}", authMiddleware(http.HandlerFunc(handler.UpdateStore)))
	mux.Handle("POST /api/v1/stores/{store_id}/activate", authMiddleware(http.HandlerFunc(handler.ActivateStore)))
	mux.Handle("POST /api/v1/stores/{store_id}/deactivate", authMiddleware(http.HandlerFunc(handler.DeactivateStore)))
}
